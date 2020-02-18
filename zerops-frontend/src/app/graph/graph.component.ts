import {AfterContentInit, Component, HostListener, ViewChild} from "@angular/core";
import {ConfigModalComponent} from "./config-modal/config-modal.component";
import {drawSvg} from "../../externalized/util/d3Helper";
import {DataSource, Pod, Step} from "../../externalized/definitions/definitions";
import {
  getDataSourcesFromRawDataAndSaveToMap,
  getPodsFromRawDataAndSaveToMap,
  getStepsFromRawDataAndSaveToMap
} from "../../externalized/functionalities/data-aggregation";
import {
  getAllDataSources,
  getAllPods,
  getAllSteps,
  setCurrentDataSources,
  setCurrentPods,
  setCurrentSteps
} from "../../externalized/functionalities/quality-of-life-functions";


function initializeMaps() {
  getPodsFromRawDataAndSaveToMap();
  getStepsFromRawDataAndSaveToMap();
  getDataSourcesFromRawDataAndSaveToMap();
}

function displayGraph(this: any, dataSources: DataSource[], steps: Step[], pods: Pod[]): void {
  setCurrentDataSources(dataSources);
  setCurrentSteps(steps);
  setCurrentPods(pods);

  // let graphVisualization: GraphVisualization = {graphColumns: []}
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
