import {DataSource, dataSourceMap, Ingest, Label, Output, Pod, podMap, Step, stepMap} from "../definitions/definitions";
import {dataSourcesRuntime, podsRuntime, stepDataSourceMatches, stepsRuntime} from "../data/data";
import {getAllPods} from "./quality-of-life-functions";
import {
  dataSourcesLink,
  matchingSourcesLink,
  podsLink,
  stepsLink,
  useLocalDataSources,
  useLocalMatchingSources,
  useLocalPods,
  useLocalSteps
} from "../config/config";

export async function getStepsFromRawDataAndSaveToMap() {
  async function getRawStepsFromProxy(): Promise<any> {
    if (useLocalSteps) {
      return stepsRuntime;
    }
    return await fetch(stepsLink)
      .then(function (response) {
        return response.json();
      });
  }

  async function getStepsFromRawData(): Promise<Step[]> {
    let stepsRaw = await getRawStepsFromProxy();

    return stepsRaw.map(stepRaw => {
      let name: string = stepRaw.metadata.name;
      let ingests: Ingest[] = stepRaw.spec.ingest;
      let outputs: Output[] = stepRaw.spec.outputs?.map(output => {
        let keys: string[] = Object.keys(output.labels);
        let labels: Label[] = keys.map(key => ({key: key, value: output.labels[key]}));
        return {
          name: output.name,
          url: output.url,
          labels: labels
        }
      });
      let validationError: string = stepRaw.status.validationError;
      let template: string = JSON.stringify(stepRaw.spec.template, null, 2);
      return {
        name: name,
        ingests: ingests,
        outputs: outputs,
        validationError: validationError,
        template: template,
        podType: 'pod',
        pods: getAllPods().filter(pod => pod.creatorStep?.name === name),
        raw: JSON.stringify(stepRaw, null, 2)
      } as Step;
    });
  }

  (await getStepsFromRawData()).forEach(step => stepMap.set(step.name, step));
}

export async function getDataSourcesFromRawDataAndSaveToMap() {
  async function getRawDataSourcesFromProxy(): Promise<any> {
    if (useLocalDataSources) {
      return dataSourcesRuntime;
    }
    return await fetch(dataSourcesLink)
      .then(function (response) {
        return response.json();
      });
  }

  async function getDataSourcesFromRawData(): Promise<DataSource[]> {
    let dataSourcesRaw = await getRawDataSourcesFromProxy();

    return <DataSource[]>dataSourcesRaw.map(dataSourceRaw => {
      let name = dataSourceRaw.metadata.name;
      let labels: Label[] = Object.keys(dataSourceRaw.metadata.labels).map(key => ({
        key: key,
        value: dataSourceRaw.metadata.labels[key]
      }));
      let specUrl: string = dataSourceRaw.spec.url;
      let validationError: string = dataSourceRaw.status.validationError;
      let creatorPodName = dataSourceRaw.metadata.labels['bitflow-pod'];
      let hasCreatorPod = false;
      if (creatorPodName != undefined) {
        hasCreatorPod = true;
      }
      let outputName = dataSourceRaw.metadata.labels['bitflow-pod-output'];
      let hasOutputName = outputName != undefined;
      return {
        name: name,
        labels: labels,
        specUrl: specUrl,
        validationError: validationError,
        hasCreatorPod: hasCreatorPod,
        creatorPod: hasCreatorPod ? {name: creatorPodName, hasCreatorStep: true} : undefined,
        hasOutputName: hasOutputName,
        outputName: outputName,
        createdPods: [],
        raw: JSON.stringify(dataSourceRaw, null, 2)
      };
    });
  }

  (await getDataSourcesFromRawData()).forEach(dataSource => dataSourceMap.set(dataSource.name, dataSource));
}

export async function getPodsAndStepsFromRawDataAndSaveToMap() {
  async function getRawPodsFromProxy(): Promise<any> {
    if (useLocalPods) {
      return podsRuntime;
    }
    return await fetch(podsLink)
      .then(function (response) {
        return response.json();
      });
  }

  async function getMatchingSourcesFromProxy(): Promise<any> {
    if (useLocalMatchingSources) {
      return stepDataSourceMatches;
    }
    return await fetch(matchingSourcesLink)
      .then(function (response) {
        return response.json();
      });
  }

  async function getPodsFromRawData(): Promise<Pod[]> {
    let podsRaw = await getRawPodsFromProxy();

    let matchingSources = await getMatchingSourcesFromProxy();

    return podsRaw
      .map(podRaw => {
        let name: string = podRaw.metadata.name;
        let phase: string = podRaw.status.phase;
        let creatorStepName = podRaw.metadata.labels['bitflow-step-name'];

        let hasCreatorStep = creatorStepName != undefined;

        let creatorDataSourceName = podRaw.metadata.labels['bitflow-source-name'];
        let creatorDataSourceNames: string[];
        if (creatorDataSourceName != undefined) {
          creatorDataSourceNames = [creatorDataSourceName];
        } else if (creatorStepName != undefined) {
          creatorDataSourceNames = matchingSources[creatorStepName].filter(name => name != undefined);
        } else {
          creatorDataSourceNames = [];
        }
        let creatorDataSources = creatorDataSourceNames.map(name => dataSourceMap.get(name)).filter(dataSource => dataSource != undefined);

        return {
          name: name,
          phase: phase,
          hasCreatorStep: hasCreatorStep,
          creatorStep: {name: creatorStepName, podType: 'pod'},
          creatorDataSources: creatorDataSources,
          createdDataSources: [],
          raw: JSON.stringify(podRaw, null, 2)
        } as Pod;
      });
  }

  (await getPodsFromRawData()).forEach(pod => podMap.set(pod.name, pod));

  await getStepsFromRawDataAndSaveToMap();

  getAllPods().forEach(pod => {
    if (pod.hasCreatorStep) {
      pod.creatorStep = stepMap.get(pod.creatorStep.name);
    }
  });
}
