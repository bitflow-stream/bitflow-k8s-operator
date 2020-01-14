import {AfterContentInit, Component, HostListener} from '@angular/core';
import {
  AnalysisType,
  DataSource,
  DataSourceGraphElement,
  DataSourceLabelKeyValuePair,
  dataSourceMap,
  DataSourceStepMatch,
  KeyValuePair,
  KubernetesGraph,
  Step,
  StepGraphElement,
  StepKeyValuePair,
  stepMap
} from './definitions/definitions';
import {uuidv4} from "./util/util";
import {dataSourcesRaw, stepsRaw} from "./data/data";
import {drawSvg} from "./util/d3Helper";

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

function stepMatchesDataSource(step: Step, dataSource: DataSource): boolean {
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

function getDataSourceStepMatches(dataSources: string[], steps: string[]): DataSourceStepMatch[] {
  let dataSourceStepMatches: DataSourceStepMatch[] = [];
  for (let dataSource of dataSources) {
    for (let step of steps) {
      if (stepMatchesDataSource(stepMap.get(step), dataSourceMap.get(dataSource))) {
        dataSourceStepMatches.push({dataSource: dataSource, step: step})
      }
    }
  }
  return dataSourceStepMatches;
}

function createKubernetesGraph(dataSources: DataSource[], steps: Step[]): KubernetesGraph {
  let dataSourceGraphElements: DataSourceGraphElement[] = [];
  dataSources.forEach(dataSource => dataSourceGraphElements.push({
    uuid: dataSource.uuid,
    stepGraphElements: [],
    creatorStepGraphElement: null
  }));

  let stepGraphElements: StepGraphElement[] = [];
  steps.forEach(step => stepGraphElements.push({
    uuid: step.uuid,
    outputDataSourceGraphElements: [],
    sourceDataSourceGraphElements: []
  }));

  let dataSourceStepMatches: DataSourceStepMatch[] = getDataSourceStepMatches(dataSources.map(dataSource => dataSource.uuid), steps.map(step => step.uuid));
  dataSourceStepMatches.forEach(match => {
    let dataSourceGraphElement: DataSourceGraphElement = dataSourceGraphElements.filter(dataSourceGraphElement => dataSourceGraphElement.uuid === match.dataSource)[0];
    let stepGraphElement: StepGraphElement = stepGraphElements.filter(stepGraphElement => stepGraphElement.uuid === match.step)[0];

    dataSourceGraphElement.stepGraphElements.push(stepGraphElement.uuid);
    stepGraphElement.sourceDataSourceGraphElements.push(dataSourceGraphElement.uuid);
  });

  stepGraphElements.forEach(stepGraphElement => {
    let currentStep: Step = stepMap.get(stepGraphElement.uuid);
    if (currentStep.type === AnalysisType.ONE_TO_ONE) {
      stepGraphElement.sourceDataSourceGraphElements.forEach(sourceDataSourceGraphElement => {
        currentStep.outputLabelsArray.forEach(outputLabels => {
          let currentDataSource: DataSource = dataSourceMap.get(sourceDataSourceGraphElement);
          let outputDataSource: DataSource = {
            uuid: uuidv4(),
            name: currentDataSource.name + '->' + currentStep.name,
            labels: [...currentDataSource.labels, ...outputLabels], // TODO prevent doubles, overwrite old ones
            depth: currentDataSource.depth + 1
          };
          dataSourceMap.set(outputDataSource.uuid, outputDataSource);
          dataSourceGraphElements.push({
            uuid: outputDataSource.uuid,
            stepGraphElements: [],
            creatorStepGraphElement: currentStep.uuid
          });
          stepGraphElements.filter(stepGraphElement => stepGraphElement.uuid === currentStep.uuid)[0].outputDataSourceGraphElements.push(outputDataSource.uuid);
        });
      });
    } else if (currentStep.type === AnalysisType.ALL_TO_ONE) {
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
          depth: Math.max(...sourceDataSources.map(dataSource => dataSource.depth)) + 1
        };

        dataSourceMap.set(outputDataSource.uuid, outputDataSource);
        dataSourceGraphElements.push({
          uuid: outputDataSource.uuid,
          stepGraphElements: [],
          creatorStepGraphElement: currentStep.uuid
        });
        stepGraphElements.filter(stepGraphElement => stepGraphElement.uuid === currentStep.uuid)[0].outputDataSourceGraphElements.push(outputDataSource.uuid);

      })
    }
  });


  return {
    dataSourceGraphElements: dataSourceGraphElements,
    stepGraphElements: stepGraphElements
  };
}

function fillDataSourceMap(dataSources: DataSource[]): void {
  dataSources.forEach(dataSource => {
    dataSourceMap.set(dataSource.uuid, dataSource);
  });
}

function fillStepMap(steps: Step[]): void {
  steps.forEach(step => stepMap.set(step.uuid, step));
}

function getAllDataSources(): DataSource[] {
  return Array.from(dataSourceMap.keys()).map(dataSourceKey => dataSourceMap.get(dataSourceKey));
}

function getAllSteps(): Step[] {
  return Array.from(stepMap.keys()).map(stepKey => stepMap.get(stepKey));
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements AfterContentInit {
  title = 'zerops-frontend';

  @HostListener('click', ['$event.target']) onClick(target) {
    if (!target.closest('rect')) return;
    alert(target.id);
  }

  ngAfterContentInit() {
    let initialDataSources: DataSource[] = dataSourcesRaw.items.map(dataSourceRaw => ({
      name: dataSourceRaw.metadata.name,
      labels: Object.keys(dataSourceRaw.metadata.labels).map(dataSourceLabelKey => ({
        key: dataSourceLabelKey,
        value: dataSourceRaw.metadata.labels[dataSourceLabelKey]
      })),
      uuid: uuidv4(),
      depth: 0
    }));
    fillDataSourceMap(initialDataSources);
    let initialSteps: Step[] = stepsRaw.items.map(stepRaw => ({
      name: stepRaw.metadata.name,
      keyValuePairs: stepRaw.spec.ingest.map(ingest => ({
        regex: ingest.check === 'regex',
        key: ingest.key,
        value: ingest.value
      })),
      uuid: uuidv4(),
      type: stepRaw.spec.type,
      outputLabelsArray: stepRaw.spec.outputs.map(output => {
        return Object.keys(output.labels).map(labelKey => ({
          key: labelKey,
          value: output.labels[labelKey]
        }));
      })
    }));
    fillStepMap(initialSteps);

    let kubernetesGraph: KubernetesGraph = createKubernetesGraph(getAllDataSources(), getAllSteps());

    drawSvg.call(this, kubernetesGraph);

  }


}
