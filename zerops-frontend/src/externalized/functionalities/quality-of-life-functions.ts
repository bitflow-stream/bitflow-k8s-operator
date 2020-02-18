import {
  currentDataSourcesMap,
  currentPodsMap,
  currentStepsMap,
  DataSource,
  dataSourceMap,
  DataSourceStack,
  GraphElement,
  Pod,
  podMap,
  PodStack,
  Step,
  stepMap
} from "../definitions/definitions";

export function getAllGraphElements(): GraphElement[] {
  return [...getAllDataSources().map(dataSource => ({type: 'data-source', dataSource: dataSource} as GraphElement)),
    ...getAllSteps().map(step => ({type: 'step', step: step} as GraphElement)),
    ...getAllPods().map(pod => ({type: 'pod', pod: pod} as GraphElement))];
}

export function getAllCurrentGraphElements(): GraphElement[] {
  return [...getCurrentDataSources().map(dataSource => ({type: 'data-source', dataSource: dataSource} as GraphElement)),
    ...getCurrentSteps().map(step => ({type: 'step', step: step} as GraphElement)),
    ...getCurrentPods().map(pod => ({type: 'pod', pod: pod} as GraphElement))];
}

export function getAllDataSources(): DataSource[] {
  return Array.from(dataSourceMap.keys()).map(dataSourceKey => <DataSource>dataSourceMap.get(dataSourceKey)).filter(dataSource => dataSource != undefined);
}

export function getAllSteps(): Step[] {
  return Array.from(stepMap.keys()).map(stepKey => <Step>stepMap.get(stepKey)).filter(step => step != undefined);
}

export function getAllPods(): Pod[] {
  return Array.from(podMap.keys()).map(podKey => <Pod>podMap.get(podKey)).filter(pod => pod != undefined);
}

export function getCurrentDataSources(): DataSource[] {
  return Array.from(currentDataSourcesMap.keys()).map(dataSourceKey => <DataSource>currentDataSourcesMap.get(dataSourceKey)).filter(dataSource => dataSource != undefined);
}

export function getCurrentSteps(): Step[] {
  return Array.from(currentStepsMap.keys()).map(stepKey => <Step>currentStepsMap.get(stepKey)).filter(step => step != undefined);
}

export function getCurrentPods(): Pod[] {
  return Array.from(currentPodsMap.keys()).map(podKey => <Pod>currentPodsMap.get(podKey)).filter(pod => pod != undefined);
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

export function getDepthOfDataSource(elementName: string): number {
  let element = dataSourceMap.get(elementName);
  if (element == undefined) {
    return 0;
  }
  if (!element.hasCreatorPod) {
    return 0;
  }
  let depth = getDepthOfPod(element.creatorPod.name);
  if (depth == undefined) {
    depth = 0;
  }
  return depth + 1;
}

export function getDepthOfPod(podName: string): number {
  let element = podMap.get(podName);
  if (element == undefined) {
    return 0;
  }
  if (element.creatorDataSources == undefined || element.creatorDataSources.length === 0) {
    return 0;
  }

  return 1 + element.creatorDataSources.map(dataSource => {
    let depth: number = getDepthOfDataSource(dataSource.name);
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

export function getDepthOfStep(stepName: string): number {
  let element = stepMap.get(stepName);
  if (element == undefined) {
    return 0;
  }
  if ((element.podType === 'pod' && (element.pods == undefined || element.pods.length === 0)) ||
    (element.podType === 'pod-stack' && (element.podStack == undefined || element.podStack.pods == undefined || element.podStack.pods.length === 0))) {
    return 0;
  }
  let depth: number = 0;
  let podNames: string[];
  if (element.podType === 'pod') {
    podNames = element.pods.map(pod => pod.name);
  } else if (element.podType === 'pod-stack') {
    podNames = element.podStack.pods.map(pod => pod.name);
  } else {
    return 0;
  }
  podNames.forEach(podName => {
    let podDepth = getDepthOfPod(podName);
    if (podDepth != undefined && podDepth > depth) {
      depth = podDepth;
    }
  });
  return depth;
}

export function getDepthOfDataSourceStack(dataSourceStack: DataSourceStack): number {
  return dataSourceStack.dataSources.map(dataSource => {
    let depth: number | undefined = getDepthOfDataSource(dataSource.name);
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

export function getDepthOfPodStack(podStack: PodStack): number {
  return podStack.pods.map(pod => {
    let depth: number | undefined = getDepthOfPod(pod.name);
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

export function getDepthOfGraphElement(graphElement: GraphElement) {
  switch (graphElement.type) {
    case "step":
      return getDepthOfStep(graphElement.step.name);
    case "data-source":
      return getDepthOfDataSource(graphElement.dataSource.name);
    case "pod":
      return getDepthOfPod(graphElement.pod.name);
    case "data-source-stack":
      return getDepthOfDataSourceStack(graphElement.dataSourceStack);
    case "pod-stack":
      return getDepthOfPodStack(graphElement.podStack);
    default:
      return 0;
  }
}
