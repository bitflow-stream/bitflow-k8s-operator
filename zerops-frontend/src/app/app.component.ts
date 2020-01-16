import {AfterContentInit, Component, HostListener} from '@angular/core';
import {
  AnalysisType,
  DataSource,
  DataSourceGraphElement,
  dataSourceGraphElementMap,
  DataSourceLabelKeyValuePair,
  dataSourceMap,
  DataSourceStepMatch,
  KeyValuePair,
  KubernetesGraph,
  Step,
  StepGraphElement,
  stepGraphElementMap,
  StepKeyValuePair,
  stepMap
} from './definitions/definitions';
import {uuidv4} from "./util/util";
import {dataSourcesRaw, stepsRaw} from "./data/data";
import {drawSvg} from "./util/d3Helper";

let matches: DataSourceStepMatch[] = [];

function getDepthOfDataSource(elementUuid: string): number {
  let element: DataSourceGraphElement = dataSourceGraphElementMap.get(elementUuid);
  if (!element) {
    return undefined;
  }
  if (!element.creatorStepGraphElement) {
    return 0;
  }
  return getDepthOfStep(element.creatorStepGraphElement) + 1;
}

function getDepthOfStep(elementUuid: string): number {
  let element: StepGraphElement = stepGraphElementMap.get(elementUuid);
  if (!element) {
    return undefined;
  }
  if (!element.sourceDataSourceGraphElements) {
    return undefined;
  }
  let depth: number = 0;
  element.sourceDataSourceGraphElements.forEach(sourceElement => {
    let sourceDepth = getDepthOfDataSource(sourceElement);
    if (sourceDepth !== undefined && sourceDepth > depth) {
      depth = sourceDepth;
    }
  });
  return depth + 1;
}

function getDepthByUuid(uuid: string): number {
  let dataSourceGraphElement: DataSourceGraphElement = dataSourceGraphElementMap.get(uuid);
  if (dataSourceGraphElement) {
    return getDepthOfDataSource(dataSourceGraphElement.uuid)
  }
  let stepGraphElement: StepGraphElement = stepGraphElementMap.get(uuid);
  if (stepGraphElement) {
    return getDepthOfStep(stepGraphElement.uuid);
  }
  return undefined;
}

function createKubernetesGraph(): KubernetesGraph {
  function getDataSourceStepMatches(): DataSourceStepMatch[] {
    function stepMatchesDataSource(step: Step, dataSource: DataSource): boolean {
      function dataSourceLabelMatchesStepKeyValuePair(dataSourceLabel: DataSourceLabelKeyValuePair, stepKeyValuePair: StepKeyValuePair): boolean {
        if (stepKeyValuePair.regex) {
          let keyRegex: RegExp = RegExp(stepKeyValuePair.key);
          if (!keyRegex.test(dataSourceLabel.key)) {
            return false;
          }
          let valueRegex: RegExp = RegExp(stepKeyValuePair.value);
          if (!valueRegex.test(dataSourceLabel.value)) {
            return false;
          }
        } else {
          if (dataSourceLabel.key !== stepKeyValuePair.key) {
            return false;
          }
          if (dataSourceLabel.value !== stepKeyValuePair.value) {
            return false;
          }
        }
        return true;
      }

      for (let stepKeyValuePair of step.keyValuePairs) {
        let foundMatchingLabel: boolean = false;
        for (let dataSourceLabel of dataSource.labels) {
          if (dataSourceLabelMatchesStepKeyValuePair(dataSourceLabel, stepKeyValuePair)) {
            foundMatchingLabel = true;
            break;
          }
        }
        if (!foundMatchingLabel) {
          return false;
        }
      }
      return true;
    }
    function stepIsAncestorOfDataSource(stepUuid: string, dataSourceUuid: string): boolean {
      let stepGraphElement: StepGraphElement = stepGraphElementMap.get(stepUuid);
      let dataSourceGraphElement: DataSourceGraphElement = dataSourceGraphElementMap.get(dataSourceUuid);
      if (dataSourceGraphElement.creatorStepGraphElement === stepUuid) {
        return true;
      }
      let sourceStepGraphElement: StepGraphElement = stepGraphElementMap.get(dataSourceGraphElement.creatorStepGraphElement);
      if (!sourceStepGraphElement) {
        return false;
      }
      let stepIsAncestor: boolean = false;
      sourceStepGraphElement.sourceDataSourceGraphElements.forEach(element => {
        if (stepIsAncestorOfDataSource(stepUuid, element)) {
          stepIsAncestor = true;
        }
      });
      return stepIsAncestor;
    }
    function stepAlreadyMatchedDataSource(stepUuid: string, dataSourceUuid: string): boolean {
      let alreadyMatched: boolean = false;
      matches.forEach(match => {
        if (match.step === stepUuid && match.dataSource === dataSourceUuid) {
          alreadyMatched = true;
        }
      });
      return alreadyMatched;
    }

    let dataSourceStepMatches: DataSourceStepMatch[] = [];

    for (let dataSource of getAllDataSources()) {
      for (let step of getAllSteps()) {
        if (stepMatchesDataSource(step, dataSource)) {
          if (!stepIsAncestorOfDataSource(step.uuid, dataSource.uuid)) {
            if (!stepAlreadyMatchedDataSource(step.uuid, dataSource.uuid)) {
              dataSourceStepMatches.push({dataSource: dataSource.uuid, step: step.uuid})
            }
          }
        }
      }
    }
    return dataSourceStepMatches;
  }
  function connectDataSourceStepMatchesInDataStructures(dataSourceStepMatches: DataSourceStepMatch[]) {
    dataSourceStepMatches.forEach(match => {
      let dataSourceGraphElement: DataSourceGraphElement = dataSourceGraphElementMap.get(match.dataSource);
      let stepGraphElement: StepGraphElement = stepGraphElementMap.get(match.step);

      dataSourceGraphElement.stepGraphElements.push(stepGraphElement.uuid);
      stepGraphElement.sourceDataSourceGraphElements.push(dataSourceGraphElement.uuid);
    });
  }
  function createOutputDataSources() {
    function handleOneToOneStep(stepGraphElement, currentStep: Step) {
      stepGraphElement.sourceDataSourceGraphElements.forEach(sourceDataSourceGraphElement => {
        currentStep.outputLabelsArray.forEach(outputLabels => {
          let currentDataSource: DataSource = dataSourceMap.get(sourceDataSourceGraphElement);
          let outputDataSource: DataSource = {
            uuid: uuidv4(),
            name: currentDataSource.name + '->' + currentStep.name,
            labels: [...currentDataSource.labels, ...outputLabels], // TODO prevent doubles, overwrite old ones
          };
          dataSourceMap.set(outputDataSource.uuid, outputDataSource);
          dataSourceGraphElementMap.set(outputDataSource.uuid, {
            uuid: outputDataSource.uuid,
            stepGraphElements: [],
            creatorStepGraphElement: currentStep.uuid
          });
          stepGraphElementMap.get(currentStep.uuid).outputDataSourceGraphElements.push(outputDataSource.uuid);
        });
      });
    }
    function handleAllToOneStep(stepGraphElement, currentStep: Step) {
      let sourceDataSources: DataSource[] = [];
      stepGraphElement.sourceDataSourceGraphElements.forEach(sourceDataSource => sourceDataSources.push(dataSourceMap.get(sourceDataSource)));
      let labelIntersection: KeyValuePair[] = [];
      if (sourceDataSources.length > 0) {
        sourceDataSources[0].labels.forEach(label => {
          let labelIsPartOfEverySourceDataSource: boolean = true;
          for (let sourceDataSource of sourceDataSources) {
            if (sourceDataSource.labels[label.key] === undefined || sourceDataSource.labels[label.key] !== label.value) {
              labelIsPartOfEverySourceDataSource = false;
              break;
            }
          }
          if (labelIsPartOfEverySourceDataSource) {
            labelIntersection.push(label);
          }
        });
      }
      currentStep.outputLabelsArray.forEach(outputLabels => {
        let uuid: string = uuidv4();
        let outputDataSource: DataSource = {
          uuid: uuid,
          name: '(all-to-one-' + uuid + ')->' + currentStep.name,
          labels: [...labelIntersection, ...outputLabels], // TODO prevent doubles, overwrite old ones
        };

        dataSourceMap.set(outputDataSource.uuid, outputDataSource);
        dataSourceGraphElementMap.set(outputDataSource.uuid, {
          uuid: outputDataSource.uuid,
          stepGraphElements: [],
          creatorStepGraphElement: currentStep.uuid
        });
        stepGraphElementMap.get(currentStep.uuid).outputDataSourceGraphElements.push(outputDataSource.uuid);
      });
    }

    getAllStepGraphElements().forEach(stepGraphElement => {
      let currentStep: Step = stepMap.get(stepGraphElement.uuid);
      if (currentStep.type === AnalysisType.ONE_TO_ONE) {
        handleOneToOneStep(stepGraphElement, currentStep);
      } else if (currentStep.type === AnalysisType.ALL_TO_ONE) {
        handleAllToOneStep(stepGraphElement, currentStep);
      }
    });
  }
  function buildKubernetesGraph() {
    return {
      dataSourceGraphElements: getAllDataSourceGraphElements(),
      stepGraphElements: getAllStepGraphElements()
    };
  }

  let dataSourceStepMatches: DataSourceStepMatch[] = getDataSourceStepMatches();
  console.log(dataSourceStepMatches);
  dataSourceStepMatches.forEach(match => matches.push(match));
  connectDataSourceStepMatchesInDataStructures(dataSourceStepMatches);
  createOutputDataSources();

  // dataSourceStepMatches = getDataSourceStepMatches();
  // dataSourceStepMatches.forEach(match => matches.push(match));
  // connectDataSourceStepMatchesInDataStructures(dataSourceStepMatches);
  // createOutputDataSources();

  return buildKubernetesGraph();
}

function getAllDataSources(): DataSource[] {
  return Array.from(dataSourceMap.keys()).map(dataSourceKey => dataSourceMap.get(dataSourceKey));
}

function getAllSteps(): Step[] {
  return Array.from(stepMap.keys()).map(stepKey => stepMap.get(stepKey));
}

function getAllDataSourceGraphElements(): DataSourceGraphElement[] {
  return Array.from(dataSourceGraphElementMap.keys()).map(dataSourceGraphElementKey => dataSourceGraphElementMap.get(dataSourceGraphElementKey));
}

function getAllStepGraphElements(): StepGraphElement[] {
  return Array.from(stepGraphElementMap.keys()).map(stepGraphElementKey => stepGraphElementMap.get(stepGraphElementKey));
}

function getRawDataSourcesAndSaveToMaps() {
  function getDataSourcesFromRawData() {
    return dataSourcesRaw.items.map(dataSourceRaw => ({
      name: dataSourceRaw.metadata.name,
      labels: Object.keys(dataSourceRaw.metadata.labels).map(dataSourceLabelKey => ({
        key: dataSourceLabelKey,
        value: dataSourceRaw.metadata.labels[dataSourceLabelKey]
      })),
      uuid: uuidv4(),
      depth: 0
    }));
  }

  function fillDataSourceMap(dataSources: DataSource[]): void {
    dataSources.forEach(dataSource => {
      dataSourceMap.set(dataSource.uuid, dataSource);
    });
  }

  function fillDataSourceGraphElementMap(dataSources: DataSource[]): void {
    dataSources.forEach(dataSource => {
      dataSourceGraphElementMap.set(dataSource.uuid, {
        uuid: dataSource.uuid,
        stepGraphElements: [],
        creatorStepGraphElement: null
      });
    });
  }

  let initialDataSources: DataSource[] = getDataSourcesFromRawData();
  fillDataSourceMap(initialDataSources);
  fillDataSourceGraphElementMap(initialDataSources);
}

function getRawStepsAndSaveToMaps() {
  function getStepsFromRawData() {
    function getStepsRaw(): any { // TODO to fix TSLint and allow multiple output labels
      return stepsRaw;
    }

    return getStepsRaw().items.map(stepRaw => ({
      name: stepRaw.metadata.name,
      keyValuePairs: stepRaw.spec.ingest.map(ingest => ({
        regex: ingest.check === 'regex',
        key: ingest.key,
        value: ingest.value
      })),
      uuid: uuidv4(),
      type: stepRaw.spec.type,
      outputLabelsArray: stepRaw.spec.outputs.map(output => {
        return Object.keys(output.labels).map((labelKey): KeyValuePair => ({
          key: labelKey,
          value: output.labels[labelKey]
        }));
      })
    }));
  }

  function fillStepMap(steps: Step[]): void {
    steps.forEach(step => stepMap.set(step.uuid, step));
  }

  function fillStepGraphElementMap(steps: Step[]): void {
    steps.forEach(step => stepGraphElementMap.set(step.uuid, {
      uuid: step.uuid,
      outputDataSourceGraphElements: [],
      sourceDataSourceGraphElements: []
    }));
  }

  let initialSteps: Step[] = getStepsFromRawData();
  fillStepMap(initialSteps);
  fillStepGraphElementMap(initialSteps);
}

function getNodeLayoutFromKubernetesGraph(kubernetesGraph: KubernetesGraph) {
  let nodeLayout: string[][] = [];

  let maxDataSourceDepth: number = kubernetesGraph.dataSourceGraphElements.map(element => getDepthOfDataSource(element.uuid)).reduce((p, c) => {
    if (!c) {
      return p;
    }
    return Math.max(p, c)
  });

  let maxStepDepth: number = kubernetesGraph.stepGraphElements.map(element => getDepthOfStep(element.uuid)).reduce((p, c) => {
    if (!c) {
      return p;
    }
    return Math.max(p, c)
  });

  let maxDepth = Math.max(maxDataSourceDepth, maxStepDepth);

  for (let i = 0; i <= maxDepth; i++) {
    nodeLayout[i] = [];
  }

  kubernetesGraph.dataSourceGraphElements.forEach(dataSourceGraphElement => {
    let depth: number = getDepthOfDataSource(dataSourceGraphElement.uuid);
    nodeLayout[depth].push(dataSourceGraphElement.uuid);
  });
  kubernetesGraph.stepGraphElements.forEach(stepGraphElement => {
    let depth: number = getDepthOfStep(stepGraphElement.uuid);
    nodeLayout[depth].push(stepGraphElement.uuid);
  });
  return nodeLayout;
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements AfterContentInit {
  @HostListener('click', ['$event.target']) onClick(target) {
    if (!target.closest('rect')) return;
    alert(target.id);
  }

  ngAfterContentInit() {
    getRawDataSourcesAndSaveToMaps();
    getRawStepsAndSaveToMaps();
    let kubernetesGraph: KubernetesGraph = createKubernetesGraph();
    let nodeLayout: string[][] = getNodeLayoutFromKubernetesGraph(kubernetesGraph);
    drawSvg.call(this, kubernetesGraph, nodeLayout);
  }
}
