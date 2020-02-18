import {AfterContentInit, Component, HostListener, ViewChild} from "@angular/core";
import {ConfigModalComponent} from "./config-modal/config-modal.component";
import {drawSvg} from "../../externalized/util/d3Helper";
import {DataSource, GraphVisualization, Pod, podMap, Step} from "../../externalized/definitions/definitions";
import {
  getDataSourcesFromRawDataAndSaveToMap,
  getPodsAndStepsFromRawDataAndSaveToMap
} from "../../externalized/functionalities/data-aggregation";
import {
  getAllCurrentGraphElements,
  getAllDataSources,
  getAllPods,
  getAllSteps,
  getDepthOfGraphElement,
  setCurrentDataSources,
  setCurrentPods,
  setCurrentSteps
} from "../../externalized/functionalities/quality-of-life-functions";


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
  let numberOfColumns = getAllCurrentGraphElements().map(element => {
    return getDepthOfGraphElement(element);
  }).reduce((p, c) => {
    if (c == undefined) {
      return p;
    }
    return Math.max(p, c)
  });

  console.log(numberOfColumns);

  return undefined;
}

function displayGraph(this: any, dataSources: DataSource[], steps: Step[], pods: Pod[]): void {
  setCurrentGraphElements(dataSources, steps, pods);

  let graphVisualization: GraphVisualization = getGraphVisualization();
  console.log(graphVisualization);
  drawSvg.call(this, {nodes: [], edges: []});
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
