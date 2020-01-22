import {AfterContentInit, Component, HostListener} from '@angular/core';
import {
  DataSource,
  dataSourceMap,
  KubernetesGraph,
  KubernetesNode,
  Pod,
  podMap,
  Step,
  stepMap
} from './definitions/definitions';
import {dataSourcesRuntime, podsRuntime, stepDataSourceMatches, stepsRuntime} from "./data/data";
import {drawSvg} from "./util/d3Helper";

function getDepthOfDataSource(elementName: string): number {
  let element: DataSource = dataSourceMap.get(elementName);
  if (!element) {
    return undefined;
  }
  if (!element.creatorPodName) {
    return 0;
  }
  return getDepthOfPod(element.creatorPodName) + 1;
}

function getDepthOfPod(podName: string): number {
  let element: Pod = podMap.get(podName);
  if (!element) {
    return undefined;
  }
  if (element.creatorDataSourceNames.length === 0) {
    return undefined;
  }

  return 1 + element.creatorDataSourceNames.map(dataSourceName => getDepthOfDataSource(dataSourceName)).reduce((p, c) => {
    if (!c) {
      return p;
    }
    return Math.max(p, c)
  });
}

function getDepthOfStep(stepName: string): number {
  let element: Step = stepMap.get(stepName);
  if (!element) {
    return undefined;
  }
  if (element.podNames.length === 0) {
    return 0;
  }
  let depth: number = 0;
  element.podNames.forEach(podName => {
    let podDepth = getDepthOfPod(podName);
    if (podDepth !== undefined && podDepth > depth) {
      depth = podDepth;
    }
  });
  return depth;
}

function getAllDataSources(): DataSource[] {
  return Array.from(dataSourceMap.keys()).map(dataSourceKey => dataSourceMap.get(dataSourceKey));
}

function getAllSteps(): Step[] {
  return Array.from(stepMap.keys()).map(stepKey => stepMap.get(stepKey));
}

function getAllPods(): Pod[] {
  return Array.from(podMap.keys()).map(podKey => podMap.get(podKey));
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
    return dataSourcesRuntime.map(dataSourceRaw => {
      let name: string = dataSourceRaw.metadata.name;
      let creatorPodName: string = dataSourceRaw.metadata.labels['zerops-pod'];
      let dataSource: DataSource = {
        name: name
      };
      if (creatorPodName) {
        dataSource.creatorPodName = creatorPodName;
      }
      return dataSource;
    });
  }

  getDataSourcesFromRawData().forEach(dataSource => dataSourceMap.set(dataSource.name, dataSource));
}

function getPodsFromRawDataAndSaveToMap() {
  function getPodsFromRawData(): Pod[] {
    return podsRuntime
      .filter(podRaw => {
        if (podRaw.metadata.labels['zerops-analysis-step']) {
          return true;
        }
      })
      .map(podRaw => {
        let name: string = podRaw.metadata.name;
        let creatorStepName: string = podRaw.metadata.labels['zerops-analysis-step'];
        let creatorDataSourceName: string = podRaw.metadata.labels['zerops-data-source-name'];
        let creatorDataSourceNames: string[] = [];
        if (creatorDataSourceName) {
          creatorDataSourceNames = [creatorDataSourceName];
        } else {
          creatorDataSourceNames = stepDataSourceMatches[creatorStepName]; // TODO undefined check
        }
        return {
          name: name,
          creatorStepName: creatorStepName,
          creatorDataSourceNames: creatorDataSourceNames
        } as Pod;
      });
  }

  getPodsFromRawData().forEach(pod => {
    stepMap.get(pod.creatorStepName).podNames.push(pod.name);
    podMap.set(pod.name, pod)
  });
}

function generateNodeLayout() {
  let nodeLayout: KubernetesNode[][] = [];

  let maxDataSourceDepth: number = getAllDataSources().map(element => getDepthOfDataSource(element.name)).reduce((p, c) => {
    if (!c) {
      return p;
    }
    return Math.max(p, c)
  });

  let maxStepDepth: number = getAllSteps().map(element => getDepthOfStep(element.name)).reduce((p, c) => {
    if (!c) {
      return p;
    }
    return Math.max(p, c)
  });

  let maxPodDepth: number = getAllPods().map(element => getDepthOfPod(element.name)).reduce((p, c) => {
    if (!c) {
      return p;
    }
    return Math.max(p, c)
  });

  let maxDepth = Math.max(maxDataSourceDepth, maxStepDepth, maxPodDepth);

  for (let i = 0; i <= maxDepth; i++) {
    nodeLayout[i] = [];
  }

  // getAllDataSources().filter(dataSource => !dataSource.creatorPodName).forEach(dataSource => nodeLayout[0].push(dataSource));

  getAllDataSources().forEach(dataSource => nodeLayout[getDepthOfDataSource(dataSource.name)].push(dataSource));
  getAllSteps().forEach(step => nodeLayout[getDepthOfStep(step.name)].push(step));

  // TODO Only put DataSources and steps, this way the height of steps can be calculated by the number of pods inside and pods are bundled inside the steps

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
    getStepsFromRawDataAndSaveToMap();
    getDataSourcesFromRawDataAndSaveToMap();
    getPodsFromRawDataAndSaveToMap();
    let nodeLayout: KubernetesNode[][] = generateNodeLayout();
    let kubernetesGraph: KubernetesGraph = {
      dataSources: getAllDataSources(),
      steps: getAllSteps(),
      pods: getAllPods()
    };
    drawSvg.call(this, kubernetesGraph, nodeLayout);
  }
}
