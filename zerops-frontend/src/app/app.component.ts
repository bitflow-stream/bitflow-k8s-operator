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
import {
  svgHorizontalGap,
  svgNodeHeight,
  svgNodeMargin,
  svgNodeWidth,
  svgPodNodeMargin,
  svgVerticalGap
} from "./config/config";
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

// function getAllPods(): Pod[] {
//   return Array.from(podMap.keys()).map(podKey => <Pod>podMap.get(podKey)).filter(pod => pod != undefined);
// }

function getStepsFromRawDataAndSaveToMap() {
  function getStepsFromRawData(): Step[] {
    return stepsRuntime.map(stepRaw => {
      let name: string = stepRaw.metadata.name;
      return {
        name: name,
        podNames: [],
        podStackId: uuidv4()
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
        creatorPodName: creatorPodName,
        dataSourceStackId: uuidv4()
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

function getVisualizationData(dataSources: DataSource[], steps: Step[]): VisualizationData | undefined {
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

  function initializeKubernetesGraph(maxColumnId: number) {
    for (let i = 0; i <= maxColumnId; i++) {
      kubernetesGraph[i] = [];
    }
  }

  function fillKubernetesGraph(dataSources: DataSource[], steps: Step[]) {
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
      let dataSourceStackId: string = uuidv4();
      dataSourceGroup.forEach(dataSource => {
        dataSource.dataSourceStackId = dataSourceStackId;
      });
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
  }

  function buildVisualizationData(maxColumnId: number): VisualizationData {
    function handleColumn(columnId: number, nodes: D3Node[], edges: D3Edge[]) {
      function handleRow(currentHeight: number, columnId: number, rowId: number, nodes: D3Node[], edges: D3Edge[]): number {
        function handleDataSources(currentHeight: number, dataSources: DataSource[], columnId: number, nodes: D3Node[], edges: D3Edge[]) {
          dataSources.forEach(dataSource => {
            if (dataSource.creatorPodName != undefined) {
              let pod: Pod | undefined = podMap.get(dataSource.creatorPodName);
              if (pod != undefined) {
                let step: Step | undefined = stepMap.get(pod.creatorStepName);
                if (step != undefined) {
                  edges.push({
                    start: step.podStackId,
                    stop: dataSource.dataSourceStackId
                  });
                }
              }
            }
          });

          if (dataSources.length > 1) {
            nodes.push({
              id: dataSources[0].dataSourceStackId,
              text: 'data-source-stack-' + dataSources[0].dataSourceStackId + ' (' + dataSources.length + ')',
              x: columnId * (svgNodeWidth + svgHorizontalGap) + svgNodeMargin,
              y: currentHeight + svgNodeMargin,
              width: svgNodeWidth,
              height: svgNodeHeight,
              type: 'data-source-stack'
            });
          } else if (dataSources.length === 1) {
            let dataSource: DataSource = dataSources[0];
            nodes.push({
              id: dataSources[0].dataSourceStackId,
              text: dataSource.name,
              x: columnId * (svgNodeWidth + svgHorizontalGap) + svgNodeMargin,
              y: currentHeight + svgNodeMargin,
              width: svgNodeWidth,
              height: svgNodeHeight,
              type: 'data-source'
            });
          }
          return currentHeight + svgNodeHeight + svgVerticalGap;
        }

        function handleStep(currentHeight: number, step: Step, columnId: number, nodes: D3Node[]) {
          nodes.push({
            id: step.name,
            text: step.name,
            x: columnId * (svgNodeWidth + svgHorizontalGap) - svgPodNodeMargin + svgNodeMargin,
            y: currentHeight - svgPodNodeMargin + svgNodeMargin,
            width: svgNodeWidth +  2 * svgPodNodeMargin,
            height: svgNodeHeight + 2 * svgPodNodeMargin,
            type: 'step'
          });

          if (step.podNames.length > 1) {
            nodes.push({
              id: step.podStackId,
              text: 'pod-stack-' + uuidv4(),
              x: columnId * (svgNodeWidth + svgHorizontalGap) + svgNodeMargin,
              y: currentHeight + .5 * svgPodNodeMargin + svgNodeMargin,
              width: svgNodeWidth,
              height: svgNodeHeight,
              type: 'pod-stack'
            });
          } else if (step.podNames.length === 1) {
            let pod: Pod | undefined = podMap.get(step.podNames[0]);
            if (pod != undefined) {
              nodes.push({
                id: step.podStackId,
                text: pod.name,
                x: columnId * (svgNodeWidth + svgHorizontalGap) + svgNodeMargin,
                y: currentHeight + .5 * svgPodNodeMargin + svgNodeMargin,
                width: svgNodeWidth,
                height: svgNodeHeight,
                type: 'pod'
              });
            }
          }

          step.podNames.forEach(podName =>  {
            let pod: Pod | undefined = podMap.get(podName);
            if (pod != undefined) {
              pod.creatorDataSourceNames.forEach(creatorDataSourceName => {
                let dataSource: DataSource | undefined = dataSourceMap.get(creatorDataSourceName);
                if (dataSource != undefined) {
                  edges.push({
                    start: dataSource.dataSourceStackId,
                    stop: step.podStackId
                  });
                }
              });
            }
          });

          return currentHeight + svgNodeHeight + svgVerticalGap;
        }

        let currentNode: KubernetesNode = kubernetesGraph[columnId][rowId];
        if (currentNode.dataSources != undefined) {
          currentHeight = handleDataSources(currentHeight, currentNode.dataSources, columnId, nodes, edges);
        } else if (currentNode.step != undefined) {
          currentHeight = handleStep(currentHeight, currentNode.step, columnId, nodes);
        }
        return currentHeight;
      }

      let currentHeight: number = 0;
      for (let rowId = 0; rowId < kubernetesGraph[columnId].length; rowId++) {
        currentHeight = handleRow(currentHeight, columnId, rowId, nodes, edges);
      }
    }

    function filterIdenticalEdges(edges: D3Edge[]): D3Edge[] {
      return edges.filter((edges, index, self) =>
        index === self.findIndex((t) => (
          t.start === edges.start && t.stop === edges.stop
        ))
      )
    }

    let nodes: D3Node[] = [];
    let edges: D3Edge[] = [];

    for (let columnId = 0; columnId <= maxColumnId; columnId++) {
      handleColumn(columnId, nodes, edges);
    }

    edges = filterIdenticalEdges(edges);

    return {
      nodes: nodes,
      edges: edges
    }
  }

  let maxColumnId: number = getMaxColumnId(dataSources, steps);

  initializeKubernetesGraph(maxColumnId);

  fillKubernetesGraph(dataSources, steps);

  return buildVisualizationData(maxColumnId);
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

    let visualizationResult: VisualizationData | undefined = getVisualizationData(getAllDataSources(), getAllSteps());
    if (visualizationResult == undefined) {
      return;
    }
    let visualization: VisualizationData = visualizationResult;
    drawSvg.call(this, visualization);
  }
}
