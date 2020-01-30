import {AfterContentInit, Component, HostListener} from '@angular/core';
import {
  D3Edge,
  D3Node,
  DataSource,
  dataSourceMap,
  kubernetesGraph,
  KubernetesNode,
  Pod,
  podMap,
  Step,
  stepMap,
  VisualizationData
} from './definitions/definitions';
import {dataSourcesRuntime, podsRuntime, stepDataSourceMatches, stepsRuntime} from "./data/data";
import {drawSvg} from "./util/d3Helper";
import {svgHorizontalGap, svgNodeHeight, svgNodeWidth, svgVerticalGap} from "./config/config";
import {uuidv4} from "./util/util";

function getDepthOfDataSource(elementName: string): number | undefined {
  let element: DataSource | undefined = dataSourceMap.get(elementName);
  if (element == undefined) {
    return undefined;
  }
  if (element.creatorPodName == undefined) {
    return 0;
  }
  let depth: number | undefined = getDepthOfPod(element.creatorPodName);
  if (depth == undefined) {
    depth = 0;
  }
  return depth + 1;
}

function getDepthOfPod(podName: string): number | undefined {
  let element: Pod | undefined = podMap.get(podName);
  if (element == undefined) {
    return undefined;
  }
  if (element.creatorDataSourceNames == undefined || element.creatorDataSourceNames.length == undefined || element.creatorDataSourceNames.length === 0) {
    return undefined;
  }

  return 1 + element.creatorDataSourceNames.map(dataSourceName => {
    let depth: number | undefined = getDepthOfDataSource(dataSourceName);
    if (depth == undefined) {
      return 0;
    }
    return depth;
  }).reduce((p, c) => {
    if (c == undefined) {
      return p;
    }
    return Math.max(p, c)
  });
}

function getDepthOfStep(stepName: string): number | undefined {
  let element: Step | undefined = stepMap.get(stepName);
  if (element == undefined) {
    return undefined;
  }
  if (element.podNames == undefined || element.podNames.length == undefined || element.podNames.length === 0) {
    return 0;
  }
  let depth: number = 0;
  element.podNames.forEach(podName => {
    let podDepth = getDepthOfPod(podName);
    if (podDepth != undefined && podDepth > depth) {
      depth = podDepth;
    }
  });
  return depth;
}

function getAllDataSources(): DataSource[] {
  return Array.from(dataSourceMap.keys()).map(dataSourceKey => <DataSource>dataSourceMap.get(dataSourceKey)).filter(dataSource => dataSource != undefined);
}

function getAllSteps(): Step[] {
  return Array.from(stepMap.keys()).map(stepKey => <Step>stepMap.get(stepKey)).filter(step => step != undefined);
}

function getAllPods(): Pod[] {
  return Array.from(podMap.keys()).map(podKey => <Pod>podMap.get(podKey)).filter(pod => pod != undefined);
}

function getStepsFromRawDataAndSaveToMap() {
  function getStepsFromRawData(): Step[] {
    return stepsRuntime.map(stepRaw => {
      let name: string = stepRaw.metadata.name;
      return {
        name: name,
        podNames: []
      } as Step;
    });
  }

  getStepsFromRawData().forEach(step => stepMap.set(step.name, step));
}

function getDataSourcesFromRawDataAndSaveToMap() {
  function getDataSourcesFromRawData(): DataSource[] {
    return <DataSource[]>dataSourcesRuntime.map(dataSourceRaw => {
      let name: string = dataSourceRaw.metadata.name;
      let creatorPodName: string | undefined = dataSourceRaw.metadata.labels['zerops-pod'];
      let dataSource: DataSource = {
        name: name,
        creatorPodName: creatorPodName
      };
      return dataSource;
    }).filter(dataSource => dataSource != undefined);
  }

  getDataSourcesFromRawData().forEach(dataSource => dataSourceMap.set(dataSource.name, dataSource));
}

function getPodsFromRawDataAndSaveToMap() {
  function getPodsFromRawData(): Pod[] {
    return podsRuntime
      .filter(podRaw => {
        return podRaw.metadata.labels['zerops-analysis-step'] != undefined;
      })
      .map(podRaw => {
        let name: string = podRaw.metadata.name;
        let creatorStepName: string = <string>podRaw.metadata.labels['zerops-analysis-step'];

        let creatorDataSourceName: string | undefined = podRaw.metadata.labels['zerops-data-source-name'];
        let creatorDataSourceNames: string[];
        if (creatorDataSourceName) {
          creatorDataSourceNames = [creatorDataSourceName];
        } else {
          creatorDataSourceNames = stepDataSourceMatches[creatorStepName];
        }
        return {
          name: name,
          creatorStepName: creatorStepName,
          creatorDataSourceNames: creatorDataSourceNames
        } as Pod;
      });
  }

  getPodsFromRawData().forEach(pod => {
    let step: Step | undefined = stepMap.get(pod.creatorStepName);
    if (step == undefined) {
      return;
    }
    step.podNames.push(pod.name);
    podMap.set(pod.name, pod)
  });
}

// function generateNodeLayout() {
//   let nodeLayout: KubernetesNode[][] = [];
//
//   let maxDataSourceDepth: number = getAllDataSources().map(element => getDepthOfDataSource(element.name)).reduce((p, c) => {
//     if (!c) {
//       return p;
//     }
//     return Math.max(p, c)
//   });
//
//   let maxStepDepth: number = getAllSteps().map(element => getDepthOfStep(element.name)).reduce((p, c) => {
//     if (!c) {
//       return p;
//     }
//     return Math.max(p, c)
//   });
//
//   let maxPodDepth: number = getAllPods().map(element => getDepthOfPod(element.name)).reduce((p, c) => {
//     if (!c) {
//       return p;
//     }
//     return Math.max(p, c)
//   });
//
//   let maxDepth = Math.max(maxDataSourceDepth, maxStepDepth, maxPodDepth);
//
//   for (let i = 0; i <= maxDepth; i++) {
//     nodeLayout[i] = [];
//   }
//
//   // getAllDataSources().filter(dataSource => !dataSource.creatorPodName).forEach(dataSource => nodeLayout[0].push(dataSource));
//
//   getAllDataSources().forEach(dataSource => nodeLayout[getDepthOfDataSource(dataSource.name)].push(dataSource));
//   getAllSteps().forEach(step => nodeLayout[getDepthOfStep(step.name)].push(step));
//
//   return nodeLayout;
// }

function getVisualizationData(dataSources: DataSource[], steps: Step[], pods: Pod[]): VisualizationData | undefined {
  pods = pods;
  function getMaxColumnId(dataSources: DataSource[], steps: Step[]): number {
    let maxColumnId: number = 0;
    dataSources.forEach(dataSource => {
      let depth: number | undefined = getDepthOfDataSource(dataSource.name);
      if (depth != undefined && maxColumnId < depth) {
        maxColumnId = depth;
      }
    });
    steps.forEach(step => {
      let depth: number | undefined = getDepthOfStep(step.name);
      if (depth != undefined && maxColumnId < depth) {
        maxColumnId = depth;
      }
    });
    return maxColumnId;
  }
  let maxColumnId: number = getMaxColumnId(dataSources, steps);
  for (let i = 0; i <= maxColumnId; i++) {
    kubernetesGraph[i] = [];
  }

  const dataSourceGroupMap: Map<string, DataSource[]> = new Map();
  dataSources.forEach(dataSource => {
    let creatorPodName: string = dataSource.creatorPodName == undefined ? 'undefined' : dataSource.creatorPodName;


    let creatorStepName: string | undefined = ((): string | undefined => {
      if (creatorPodName == 'undefined') {
        return 'undefined';
      }
      let creatorPod: Pod | undefined = podMap.get(creatorPodName);
      if (creatorPod == undefined) {
        return undefined;
      }
      return creatorPod.creatorStepName == undefined ? 'undefined' : creatorPod.creatorStepName;
    })();

    if (creatorStepName == undefined) {
      return;
    }
    let dataSourceGroup: DataSource[] | undefined = dataSourceGroupMap.get(creatorStepName);
    if (dataSourceGroup == undefined || dataSourceGroup.length === 0) {
      dataSourceGroupMap.set(creatorStepName, [dataSource]);
    } else {
      dataSourceGroup.push(dataSource);
    }
  });

  Array.from(dataSourceGroupMap.keys()).forEach(key => {
    let dataSourceGroup: DataSource[] | undefined = dataSourceGroupMap.get(key);
    if (dataSourceGroup == undefined || dataSourceGroup.length === 0) {
      return;
    }
    let depth: number | undefined = getDepthOfDataSource(dataSourceGroup[0].name);
    if (depth == undefined) {
      return;
    }
    kubernetesGraph[depth].push({dataSources: dataSourceGroup});
  });

  steps.forEach(step => {
    let depth: number | undefined = getDepthOfStep(step.name);
    if (depth == undefined) {
      return;
    }
    kubernetesGraph[depth].push({step: step});
  });

  let nodes: D3Node[] = [];
  let edges: D3Edge[] = [];

  for (let columnId = 0; columnId <= maxColumnId; columnId++) {
    let currentHeight: number = 0;
    for (let rowId = 0; rowId < kubernetesGraph[columnId].length; rowId++) {
      let currentNode: KubernetesNode = kubernetesGraph[columnId][rowId];
      if (currentNode.dataSources != undefined) {
        let dataSources: DataSource[] = currentNode.dataSources;
        if (dataSources.length > 1) {
          nodes.push({
            id: 'data-source-stack-' + uuidv4(),
            text: 'data-source-stack-' + uuidv4() + ' (' + dataSources.length + ')',
            x: columnId * (svgNodeWidth + svgHorizontalGap),
            y: currentHeight,
            width: svgNodeWidth,
            height: svgNodeHeight,
            type: 'data-source-stack'
          });
        } else if (dataSources.length === 1) {
          let dataSource: DataSource = dataSources[0];
          nodes.push({
            id: dataSource.name,
            text: dataSource.name,
            x: columnId * (svgNodeWidth + svgHorizontalGap),
            y: currentHeight,
            width: svgNodeWidth,
            height: svgNodeHeight,
            type: 'data-source'
          });
        }
        currentHeight += svgNodeHeight + svgVerticalGap;
      } else if (currentNode.step != undefined) {
        let step: Step = currentNode.step;
        nodes.push({
          id: step.name,
          text: step.name,
          x: columnId * (svgNodeWidth + svgHorizontalGap),
          y: currentHeight,
          width: svgNodeWidth,
          height: svgNodeHeight,
          type: 'step'
        });
        currentHeight += svgNodeHeight + svgVerticalGap;
      }
    }
  }

  return {
    nodes: nodes,
    edges: edges
  };
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements AfterContentInit {
  @HostListener('click', ['$event.target']) onClick(target: any) {
    if (target.closest('rect') == undefined) return;
    alert(target.id);
  }

  ngAfterContentInit() {
    getStepsFromRawDataAndSaveToMap();
    getDataSourcesFromRawDataAndSaveToMap();
    getPodsFromRawDataAndSaveToMap();

    let visualizationResult: VisualizationData | undefined = getVisualizationData(getAllDataSources(), getAllSteps(), getAllPods());
    if (visualizationResult == undefined) {
      return;
    }
    let visualization: VisualizationData = visualizationResult;
    drawSvg.call(this, visualization);
  }
}
