import {AfterContentInit, Component, HostListener} from '@angular/core';
import {DataSource, dataSourceMap, Pod, podMap, Step, stepMap} from './definitions/definitions';
import {dataSourcesRuntime, podsRuntime, stepsRuntime} from "./data/data";

// function getDepthOfDataSource(elementUuid: string): number {
//   // let element: DataSourceGraphElement = dataSourceGraphElementMap.get(elementUuid);
//   // if (!element) {
//   //   return undefined;
//   // }
//   // if (!element.creatorStepGraphElement) {
//   //   return 0;
//   // }
//   // return getDepthOfStep(element.creatorStepGraphElement) + 1;
// }
//
// function getDepthOfStep(elementUuid: string): number {
//   // let element: StepGraphElement = stepGraphElementMap.get(elementUuid);
//   // if (!element) {
//   //   return undefined;
//   // }
//   // if (!element.sourceDataSourceGraphElements) {
//   //   return undefined;
//   // }
//   // let depth: number = 0;
//   // element.sourceDataSourceGraphElements.forEach(sourceElement => {
//   //   let sourceDepth = getDepthOfDataSource(sourceElement);
//   //   if (sourceDepth !== undefined && sourceDepth > depth) {
//   //     depth = sourceDepth;
//   //   }
//   // });
//   // return depth + 1;
// }
//
// function getDepthByUuid(uuid: string): number {
//   // let dataSourceGraphElement: DataSourceGraphElement = dataSourceGraphElementMap.get(uuid);
//   // if (dataSourceGraphElement) {
//   //   return getDepthOfDataSource(dataSourceGraphElement.uuid)
//   // }
//   // let stepGraphElement: StepGraphElement = stepGraphElementMap.get(uuid);
//   // if (stepGraphElement) {
//   //   return getDepthOfStep(stepGraphElement.uuid);
//   // }
//   // return undefined;
// }

function getAllDataSources(): DataSource[] {
  return Array.from(dataSourceMap.keys()).map(dataSourceKey => dataSourceMap.get(dataSourceKey));
}

function getAllSteps(): Step[] {
  return Array.from(stepMap.keys()).map(stepKey => stepMap.get(stepKey));
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
        dataSource['zerops-pod'] = creatorPodName;
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
        if (podRaw.metadata.labels['zerops-analysis-step'] && podRaw.metadata.labels['zerops-data-source-name']) {
          return true;
        }
      })
      .map(podRaw => {
      let name: string = podRaw.metadata.name;
      let creatorStepName: string = podRaw.metadata.labels['zerops-analysis-step'];
      let creatorDataSourceName: string = podRaw.metadata.labels['zerops-data-source-name'];
      return {
        name: name,
        creatorStepName: creatorStepName,
        creatorDataSourceName: creatorDataSourceName
      } as Pod;
    });
  }

  getPodsFromRawData().forEach(pod => {
    stepMap.get(pod.creatorStepName).podNames.push(pod.name);
    podMap.set(pod.name, pod)
  });
}

function generateNodeLayout() {
  // let nodeLayout: string[][] = [];
  //
  // let maxDataSourceDepth: number = kubernetesGraph.dataSources.map(element => getDepthOfDataSource(element.uuid)).reduce((p, c) => {
  //   if (!c) {
  //     return p;
  //   }
  //   return Math.max(p, c)
  // });
  //
  // let maxStepDepth: number = kubernetesGraph.steps.map(element => getDepthOfStep(element.uuid)).reduce((p, c) => {
  //   if (!c) {
  //     return p;
  //   }
  //   return Math.max(p, c)
  // });
  //
  // let maxDepth = Math.max(maxDataSourceDepth, maxStepDepth);
  //
  // for (let i = 0; i <= maxDepth; i++) {
  //   nodeLayout[i] = [];
  // }
  //
  // kubernetesGraph.dataSources.forEach(dataSourceGraphElement => {
  //   let depth: number = getDepthOfDataSource(dataSourceGraphElement.uuid);
  //   nodeLayout[depth].push(dataSourceGraphElement.uuid);
  // });
  // kubernetesGraph.steps.forEach(stepGraphElement => {
  //   let depth: number = getDepthOfStep(stepGraphElement.uuid);
  //   nodeLayout[depth].push(stepGraphElement.uuid);
  // });
  // return nodeLayout;
  return null;
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
    console.log(stepMap);
    console.log(dataSourceMap);
    console.log(podMap);
    let nodeLayout: string[][] = generateNodeLayout();
    // drawSvg.call(this, kubernetesGraph, nodeLayout);
  }
}
