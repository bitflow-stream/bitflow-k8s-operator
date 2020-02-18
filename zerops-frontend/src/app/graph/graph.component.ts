import {AfterContentInit, Component, HostListener, ViewChild} from "@angular/core";
import {ConfigModalComponent} from "./config-modal/config-modal.component";
import {drawSvg} from "../../externalized/util/d3Helper";
import {
  D3Edge,
  D3Node,
  DataSource,
  FrontendData,
  GraphVisualization,
  GraphVisualizationColumn,
  Pod,
  podMap,
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
  getAllSteps,
  getCurrentPods,
  getDepthOfGraphElement,
  setCurrentDataSources,
  setCurrentPods,
  setCurrentSteps
} from "../../externalized/functionalities/quality-of-life-functions";
import {
  svgHorizontalGap,
  svgNodeHeight,
  svgNodeMargin,
  svgNodeWidth, svgPodNodeMargin,
  svgVerticalGap
} from "../../externalized/config/config";

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

function getGraphVisualization() {
  let maxColumnId = getAllCurrentGraphElements().map(element => {
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

  getAllCurrentGraphElements().forEach(element => {
    let depth = getDepthOfGraphElement(element);
    let graphVisualizationColumn: GraphVisualizationColumn = graphVisualization.graphColumns[depth];
    if (element.type === 'step') {
      graphVisualizationColumn.graphElements.push(element);
      let currentPods = getCurrentPods();
      let currentPodsInStep = element.step.pods.filter(pod => currentPods.some(currentPod => currentPod.name === pod.name));
      currentPodsInStep.forEach(pod => graphVisualizationColumn.graphElements.push({type: 'pod', pod: pod}));
    }
    if (element.type === 'data-source') {
      graphVisualizationColumn.graphElements.push(element);
    }
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
