import {AfterContentInit, Component, HostListener, ViewChild} from "@angular/core";
import {ConfigModalComponent} from "./config-modal/config-modal.component";
import {drawSvg} from "../../externalized/util/d3Helper";
import {
  D3Edge,
  D3Node,
  DataSource,
  FrontendData,
  GraphElement,
  GraphVisualization,
  GraphVisualizationColumn,
  Pod,
  Step
} from "../../externalized/definitions/definitions";
import {
  getAllCurrentGraphElementsWithStacks,
  getAllDataSources,
  getAllPods,
  getAllSteps,
  getCurrentDataSources,
  getCurrentPods,
  getDepthOfGraphElement,
  getGraphElementIncludingDataSource,
  getGraphElementIncludingPod,
  initializeMaps,
  setAllCurrentGraphElementsWithStacks,
  setCurrentGraphElements
} from "../../externalized/functionalities/quality-of-life-functions";
import {
  svgHorizontalGap,
  svgNodeHeight,
  svgNodeMargin,
  svgNodeWidth,
  svgPodNodeMargin,
  svgVerticalGap
} from "../../externalized/config/config";

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

    if ((element.type === 'pod' && element.pod.hasCreatorStep) || element.type === 'pod-stack') {
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

        if (graphElement.dataSource.hasCreatorPod) {
          let creatorGraphElement: GraphElement = getGraphElementIncludingPod(graphElement.dataSource.creatorPod, getAllCurrentGraphElementsWithStacks());
          if (creatorGraphElement.type === 'pod') {
            edges.push({start: creatorGraphElement.pod.name, stop: graphElement.dataSource.name});
          }
          if (creatorGraphElement.type === 'pod-stack') {
            edges.push({start: creatorGraphElement.podStack.stackId, stop: graphElement.dataSource.name});
          }
        }
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


        graphElement.dataSourceStack.dataSources.forEach(dataSource => {
          if (dataSource.hasCreatorPod) {
            let creatorGraphElement: GraphElement = getGraphElementIncludingPod(dataSource.creatorPod, getAllCurrentGraphElementsWithStacks());
            if (creatorGraphElement.type === 'pod') {
              edges.push({start: creatorGraphElement.pod.name, stop: graphElement.dataSourceStack.stackId});
            }
            if (creatorGraphElement.type === 'pod-stack') {
              edges.push({start: creatorGraphElement.podStack.stackId, stop: graphElement.dataSourceStack.stackId});
            }
          }
        });
      }
      if (graphElement.type === 'pod') {
        nodes.push({
          id: graphElement.pod.name,
          text: graphElement.pod.name,
          x: columnId * (svgNodeWidth + svgHorizontalGap) + svgNodeMargin,
          y: 20 + currentHeight + svgNodeMargin,
          width: svgNodeWidth,
          height: svgNodeHeight,
          type: "pod"
        });
        currentHeight += svgNodeHeight + svgVerticalGap;

        graphElement.pod.creatorDataSources.filter(creator => getCurrentDataSources().some(dataSource => dataSource.name === creator.name)).forEach(creator => {
          let creatorGraphElement: GraphElement = getGraphElementIncludingDataSource(creator, getAllCurrentGraphElementsWithStacks());
          if (creatorGraphElement.type === 'data-source') {
            edges.push({start: creatorGraphElement.dataSource.name, stop: graphElement.pod.name});
          }
          if (creatorGraphElement.type === 'data-source-stack') {
            edges.push({start: creatorGraphElement.dataSourceStack.stackId, stop: graphElement.pod.name});
          }
        });

      }
      if (graphElement.type === 'pod-stack') {
        nodes.push({
          id: graphElement.podStack.stackId,
          text: graphElement.podStack.stackId,
          x: columnId * (svgNodeWidth + svgHorizontalGap) + svgNodeMargin,
          y: 20 + currentHeight + svgNodeMargin,
          width: svgNodeWidth,
          height: svgNodeHeight,
          type: "pod-stack"
        });
        currentHeight += svgNodeHeight + svgVerticalGap;

        graphElement.podStack.pods.forEach(pod => {
          pod.creatorDataSources.forEach(creator => {
            let creatorGraphElement: GraphElement = getGraphElementIncludingDataSource(creator, getAllCurrentGraphElementsWithStacks());
            if (creatorGraphElement.type === 'data-source') {
              edges.push({start: creatorGraphElement.dataSource.name, stop: graphElement.podStack.stackId}); // TODO remove duplicate edges (and make edges thicker?)
            }
            if (creatorGraphElement.type === 'data-source-stack') {
              edges.push({start: creatorGraphElement.dataSourceStack.stackId, stop: graphElement.podStack.stackId});
            }
          });
        });
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

  // Filtering identical edges
  edges = edges.filter((edges, index, self) =>
    index === self.findIndex((t) => (
      t.start === edges.start && t.stop === edges.stop
    ))
  );

  return {nodes, edges} as FrontendData;
}

function displayGraph(this: any, dataSources: DataSource[], steps: Step[], pods: Pod[]): void {
  setCurrentGraphElements(dataSources, steps, pods);

  setAllCurrentGraphElementsWithStacks();
  let graphVisualization: GraphVisualization = getGraphVisualization();
  console.log(graphVisualization)
  let frontendData: FrontendData = getFrontendDataFromGraphVisualization(graphVisualization);

  // TODO why is frontendData containing more than a pod?
  console.log(frontendData)
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


  filterGraph(graphElement: GraphElement) {
    displayGraph.call(this, [], [graphElement.pod.creatorStep], [graphElement.pod]);
  }
}
