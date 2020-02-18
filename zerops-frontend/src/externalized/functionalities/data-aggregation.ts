import {DataSource, dataSourceMap, Pod, podMap, Step, stepMap} from "../definitions/definitions";
import {dataSourcesRuntime, podsRuntime, stepDataSourceMatches, stepsRuntime} from "../data/data";
import {getAllPods} from "./quality-of-life-functions";

export function getStepsFromRawDataAndSaveToMap() {
  function getStepsFromRawData(): Step[] {
    return stepsRuntime.map(stepRaw => {
      let name = stepRaw.metadata.name;
      return {
        name: name,
        podType: 'pod',
        pods: getAllPods().filter(pod => pod.creatorStep?.name === name)
      } as Step;
    });
  }

  getStepsFromRawData().forEach(step => stepMap.set(step.name, step));
}

export function getDataSourcesFromRawDataAndSaveToMap() {
  function getDataSourcesFromRawData(): DataSource[] {
    return <DataSource[]>dataSourcesRuntime.map(dataSourceRaw => {
      let name = dataSourceRaw.metadata.name;
      let creatorPodName = dataSourceRaw.metadata.labels['zerops-pod'];
      let hasCreatorPod = false;
      let creatorPod;
      if (creatorPodName != undefined) {
        hasCreatorPod = true;
        creatorPod = podMap.get(creatorPodName);
      }
      let outputName = dataSourceRaw.metadata.labels['zerops-output'];
      let hasOutputName = outputName != undefined;
      return {
        name: name,
        hasCreatorPod: hasCreatorPod,
        creatorPod: creatorPod,
        hasOutputName: hasOutputName,
        outputName: outputName
      };
    });
  }

  getDataSourcesFromRawData().forEach(dataSource => dataSourceMap.set(dataSource.name, dataSource));
}

export function getPodsFromRawDataAndSaveToMap() {
  function getPodsFromRawData(): Pod[] {
    return podsRuntime
      .map(podRuntime => {
        let name: string = podRuntime.metadata.name;
        let creatorStepName = podRuntime.metadata.labels['zerops-analysis-step'];

        let creatorStep;
        if (creatorStepName != undefined) {
          creatorStep = stepMap.get(creatorStepName);
        }
        let hasCreatorStep = creatorStep != undefined;

        let creatorDataSourceName = podRuntime.metadata.labels['zerops-data-source-name'];
        let creatorDataSourceNames: string[];
        if (creatorDataSourceName != undefined) {
          creatorDataSourceNames = [creatorDataSourceName];
        } else if (creatorStepName != undefined) {
          creatorDataSourceNames = stepDataSourceMatches[creatorStepName].filter(name => name != undefined);
        } else {
          creatorDataSourceNames = [];
        }
        let creatorDataSources = creatorDataSourceNames.map(name => dataSourceMap.get(name)).filter(dataSource => dataSource != undefined);

        return {
          name: name,
          hasCreatorStep: hasCreatorStep,
          creatorStep: creatorStep,
          creatorDataSources: creatorDataSources
        } as Pod;
      });
  }

  getPodsFromRawData().forEach(pod => podMap.set(pod.name, pod));
}
