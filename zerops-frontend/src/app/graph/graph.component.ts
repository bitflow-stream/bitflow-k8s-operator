import {AfterContentInit, Component, HostListener, ViewChild} from "@angular/core";
import {ConfigModalComponent} from "./config-modal/config-modal.component";
import {drawSvg} from "../../externalized/util/d3Helper";
import {
  D3Edge,
  D3Node,
  DataSource, DataSourceStack,
  FrontendData, GraphElement,
  GraphVisualization,
  GraphVisualizationColumn,
  Pod,
  podMap, PodStack,
  Step
} from "../../externalized/definitions/definitions";
import {
  getDataSourcesFromRawDataAndSaveToMap,
  getPodsAndStepsFromRawDataAndSaveToMap
} from "../../externalized/functionalities/data-aggregation";
import {
  getAllCurrentGraphElements,
  getAllDataSources,
  getAllPods,
  getAllSteps, getCurrentDataSources,
  getCurrentPods, getCurrentSteps,
  getDepthOfGraphElement,
  setCurrentDataSources,
  setCurrentPods,
  setCurrentSteps
} from "../../externalized/functionalities/quality-of-life-functions";
import {
  maxNumberOfSeparateGraphElements,
  svgHorizontalGap,
  svgNodeHeight,
  svgNodeMargin,
  svgNodeWidth, svgPodNodeMargin,
  svgVerticalGap
} from "../../externalized/config/config";
import {uuidv4} from "../../externalized/util/util";

function addCreatorPodsToDataSources() {
  getAllDataSources().forEach(dataSource => {
    if (dataSource.hasCreatorPod == true) {
      dataSource.creatorPod = podMap.get(dataSource.creatorPod.name);
    }
  });
}

function initializeMaps() {
  getDataSourcesFromRawDataAndSaveToMap();
  getPodsAndStepsFromRawDataAndSaveToMap();
  addCreatorPodsToDataSources();
}

function setCurrentGraphElements(dataSources, steps, pods) {
  setCurrentDataSources(dataSources);
  setCurrentSteps(steps);
  setCurrentPods(pods);
}

function podShouldBeGroupedWithPodStack(pod: Pod, podStack: PodStack) {
  return pod.hasCreatorStep && podStack.hasCreatorStep && pod.creatorStep.name === podStack.creatorStep.name;
}

function dataSourceShouldBeGroupedWithDataSourceStack(dataSource: DataSource, dataSourceStack: DataSourceStack) {
  return (dataSourceStack.outputName === dataSource.outputName) &&
    (
      (!dataSourceStack.hasSourceGraphElement && !dataSource.hasCreatorPod) ||
      (
        (dataSourceStack.hasSourceGraphElement && dataSource.hasCreatorPod) &&
        (
          (dataSourceStack.sourceGraphElement.type === 'pod' && dataSourceStack.sourceGraphElement.pod.name === dataSource.creatorPod.name) ||
          (dataSourceStack.sourceGraphElement.type === 'pod-stack' && dataSourceStack.sourceGraphElement.podStack.pods.some(pod => pod.name === dataSource.creatorPod.name))
        )
      )
    );
}

function getGraphElementIncludingPod(pod: Pod, podGraphElements: GraphElement[]) {
  for (let i = 0; i < podGraphElements.length; i++) {
    let podGraphElement: GraphElement = podGraphElements[i];
    if (podGraphElement.type === 'pod') {
      if (pod.name === podGraphElement.pod.name) {
        return podGraphElement;
      }
    }
    if (podGraphElement.type === 'pod-stack') {
      for (let j = 0; j < podGraphElement.podStack.pods.length; j++) {
        let innerPod: Pod = podGraphElement.podStack.pods[j];
        if (pod.name === innerPod.name) {
          return podGraphElement;
        }
      }
    }
  }
  return undefined;
}

function getAllCurrentGraphElementsWithStacks(): GraphElement[] {
  let podGraphElements: GraphElement[] = [];
  let currentPods: Pod[] = getCurrentPods();

  currentPods.forEach(pod => {
    if (!pod.hasCreatorStep) {
      podGraphElements.push({type: 'pod', pod: pod});
    }
  });

  currentPods.filter(pod => pod.hasCreatorStep).forEach(pod => {
    for (let i = 0; i < podGraphElements.length; i++) {
      let podGraphElement = podGraphElements[i];
      if (podGraphElement.type != "pod-stack") {
        continue;
      }
      if (podShouldBeGroupedWithPodStack(pod, podGraphElement.podStack)) {
        podGraphElement.podStack.pods.push(pod);
        return;
      }
    }
    let podStackGraphElement: GraphElement = {
      type: "pod-stack",
      podStack: {
        stackId: uuidv4(),
        hasCreatorStep: pod.hasCreatorStep,
        creatorStep: pod.creatorStep,
        pods: [pod]
      }
    };
    podGraphElements.push(podStackGraphElement);
    if (pod.hasCreatorStep) {
      pod.creatorStep.podType = 'pod-stack';
      pod.creatorStep.pods = undefined;
      pod.creatorStep.podStack = podStackGraphElement.podStack;
    }
  });

  let dataSourceGraphElements: GraphElement[] = [];
  let currentDataSources: DataSource[] = getCurrentDataSources();
  currentDataSources.forEach(dataSource => {
    if (!dataSource.hasOutputName) {
      dataSourceGraphElements.push({type: "data-source", dataSource: dataSource});
    }
  });
  currentDataSources.filter(dataSource => dataSource.hasOutputName).forEach(dataSource => {
    for (let i = 0; i < dataSourceGraphElements.length; i++) {
      let dataSourceGraphElement = dataSourceGraphElements[i];
      if (dataSourceGraphElement.type != "data-source-stack") {
        continue;
      }
      if (dataSourceShouldBeGroupedWithDataSourceStack(dataSource, dataSourceGraphElement.dataSourceStack)) {
        dataSourceGraphElement.dataSourceStack.dataSources.push(dataSource);
        return;
      }
    }
    dataSourceGraphElements.push({
      type: "data-source-stack",
      dataSourceStack: {
        stackId: uuidv4(),
        hasSourceGraphElement: dataSource.hasCreatorPod,
        sourceGraphElement: dataSource.hasCreatorPod ? getGraphElementIncludingPod(dataSource.creatorPod, podGraphElements) : undefined,
        outputName: dataSource.outputName,
        dataSources: [dataSource]
      }
    });
  });

  let stepGraphElements: GraphElement[] = getCurrentSteps().map(step => ({type: "step", step}));

  return [
    ...dataSourceGraphElements,
    ...podGraphElements,
    ...stepGraphElements
  ]
}

function getGraphVisualization() {
  let maxColumnId = getAllCurrentGraphElementsWithStacks().map(element => {
    return getDepthOfGraphElement(element);
  }).reduce((p, c) => {
    if (c == undefined) {
      return p;
    }
    return Math.max(p, c)
  });

  let graphVisualization: GraphVisualization = {graphColumns: []};

  for (let i = 0; i <= maxColumnId; i++) {
    graphVisualization.graphColumns.push({graphElements: []});
  }

  let currentGraphElementsWithStacks: GraphElement[] = getAllCurrentGraphElementsWithStacks();

  currentGraphElementsWithStacks.forEach(element => {
    let depth = getDepthOfGraphElement(element);
    let graphVisualizationColumn: GraphVisualizationColumn = graphVisualization.graphColumns[depth];

    if (element.type === 'pod' || element.type === 'pod-stack') {
      return;
    }

    graphVisualizationColumn.graphElements.push(element);
    if (element.type === 'step') {
      if (element.step.podType === 'pod') {
        element.step.pods.forEach(pod => {
          graphVisualizationColumn.graphElements.push({type: 'pod', pod});
        });
      }
      if (element.step.podType === 'pod-stack') {
        graphVisualizationColumn.graphElements.push({type: 'pod-stack', podStack: element.step.podStack});
      }
    }






    // if (element.type === 'step') {
    //   graphVisualizationColumn.graphElements.push(element);
      // let currentPods = getCurrentPods();
      // let currentPodsInStep = element.step.pods.filter(pod => currentPods.some(currentPod => currentPod.name === pod.name));
      // if (currentPodsInStep.length <= maxNumberOfSeparateGraphElements) { //TODO wird an einigen Stellen davon ausgegangen, dass es sich um einen Stack handelt?
      //   currentPodsInStep.forEach(pod => graphVisualizationColumn.graphElements.push({type: 'pod', pod: pod}));
      // }
      // else {
      //   element.step.podType = 'pod-stack';
      //   graphVisualizationColumn.graphElements.push({type: "pod-stack", podStack: {stackId: uuidv4(), pods: currentPodsInStep}});
      // }
    // }
    // if (element.type === 'data-source') {
    //   graphVisualizationColumn.graphElements.push(element);
    // }
    // if (element.type === 'data-source-stack') {
    //   graphVisualizationColumn.graphElements.push(element);
    // }
  });

  return graphVisualization;
}

function getFrontendDataFromGraphVisualization(graphVisualization: GraphVisualization) {
  let nodes: D3Node[] = [];
  let edges: D3Edge[] = [];

  graphVisualization.graphColumns.forEach((column, columnId) => {
    let currentHeight = 0;
    column.graphElements.forEach(graphElement => {
      if (graphElement.type === 'data-source') {
        nodes.push({
          id: graphElement.dataSource.name,
          text: graphElement.dataSource.name,
          x: columnId * (svgNodeWidth + svgHorizontalGap) + svgNodeMargin,
          y: currentHeight + svgNodeMargin,
          width: svgNodeWidth,
          height: svgNodeHeight,
          type: "data-source"
        });
        currentHeight += svgNodeHeight + svgVerticalGap;
      }
      if (graphElement.type === 'data-source-stack') {
        nodes.push({
          id: graphElement.dataSourceStack.stackId,
          text: graphElement.dataSourceStack.stackId,
          x: columnId * (svgNodeWidth + svgHorizontalGap) + svgNodeMargin,
          y: currentHeight + svgNodeMargin,
          width: svgNodeWidth,
          height: svgNodeHeight,
          type: "data-source-stack"
        });
        currentHeight += svgNodeHeight + svgVerticalGap;
      }
      if (graphElement.type === 'pod') {
        nodes.push({
          id: graphElement.pod.name,
          text: graphElement.pod.name,
          x: columnId * (svgNodeWidth + svgHorizontalGap) + svgNodeMargin,
          y: currentHeight + svgNodeMargin,
          width: svgNodeWidth,
          height: svgNodeHeight,
          type: "pod"
        });
        currentHeight += svgNodeHeight + svgVerticalGap;
      }
      if (graphElement.type === 'pod-stack') {
        nodes.push({
          id: graphElement.podStack.stackId,
          text: graphElement.podStack.stackId,
          x: columnId * (svgNodeWidth + svgHorizontalGap) + svgNodeMargin,
          y: currentHeight + svgNodeMargin,
          width: svgNodeWidth,
          height: svgNodeHeight,
          type: "pod-stack"
        });
        currentHeight += svgNodeHeight + svgVerticalGap;
      }
      if (graphElement.type === 'step') {
        if (graphElement.step.podType === 'pod') {
          let currentPods = getCurrentPods();
          let currentPodsInStep = graphElement.step.pods.filter(pod => currentPods.some(currentPod => currentPod.name === pod.name));
          nodes.push({
            id: graphElement.step.name,
            text: graphElement.step.name,
            x: columnId * (svgNodeWidth + svgHorizontalGap) - svgPodNodeMargin + svgNodeMargin,
            y: currentHeight - svgPodNodeMargin + svgNodeMargin,
            width: svgNodeWidth + 2 * svgPodNodeMargin,
            height: Math.max(1, currentPodsInStep.length) * (svgNodeHeight + svgVerticalGap),
            type: 'step'
          });
          if (currentPodsInStep.length === 0) {
            currentHeight += svgNodeHeight + svgVerticalGap;
          }
        }
        if (graphElement.step.podType === 'pod-stack') {
          nodes.push({
            id: graphElement.step.name,
            text: graphElement.step.name,
            x: columnId * (svgNodeWidth + svgHorizontalGap) - svgPodNodeMargin + svgNodeMargin,
            y: currentHeight - svgPodNodeMargin + svgNodeMargin,
            width: svgNodeWidth + 2 * svgPodNodeMargin,
            height: svgNodeHeight + svgVerticalGap,
            type: 'step'
          });
        }
      }
    });
  });

  return {nodes, edges} as FrontendData;
}

function displayGraph(this: any, dataSources: DataSource[], steps: Step[], pods: Pod[]): void {
  setCurrentGraphElements(dataSources, steps, pods);

  let graphVisualization: GraphVisualization = getGraphVisualization();
  let frontendData: FrontendData = getFrontendDataFromGraphVisualization(graphVisualization);
  drawSvg.call(this, frontendData);
}

@Component({
  selector: 'app-graph',
  templateUrl: './graph.component.html',
  styleUrls: ['./graph.component.css']
})
export class GraphComponent implements AfterContentInit {
  @ViewChild(ConfigModalComponent, {static: false}) modal: ConfigModalComponent | undefined;

  @HostListener('click', ['$event.target']) onClick(target: any) {
    if (target.closest('rect') == undefined) return;
    this.modal?.openModal(target.id);
  }

  ngAfterContentInit() {
    initializeMaps();

    displayGraph.call(this, getAllDataSources(), getAllSteps(), getAllPods());
  }
}
