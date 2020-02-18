import {
  currentDataSourcesMap,
  currentPodsMap, currentStepsMap,
  DataSource,
  dataSourceMap,
  Pod,
  podMap,
  Step,
  stepMap
} from "../definitions/definitions";

export function getAllDataSources(): DataSource[] {
  return Array.from(dataSourceMap.keys()).map(dataSourceKey => <DataSource>dataSourceMap.get(dataSourceKey)).filter(dataSource => dataSource != undefined);
}

export function getAllSteps(): Step[] {
  return Array.from(stepMap.keys()).map(stepKey => <Step>stepMap.get(stepKey)).filter(step => step != undefined);
}

export function getAllPods(): Pod[] {
  return Array.from(podMap.keys()).map(podKey => <Pod>podMap.get(podKey)).filter(pod => pod != undefined);
}

export function setCurrentDataSources(dataSources: DataSource[]) {
  dataSources.forEach(dataSource => {
    currentDataSourcesMap.set(dataSource.name, dataSource);
  });
}

export function setCurrentSteps(steps: Step[]) {
  steps.forEach(step => {
    currentStepsMap.set(step.name, step);
  });
}

export function setCurrentPods(pods: Pod[]) {
  pods.forEach(pod => {
    currentPodsMap.set(pod.name, pod);
  });
}
