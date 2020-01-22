import {
  D3Edge,
  D3Node,
  dataSourceMap,
  KubernetesGraph,
  KubernetesNode, Pod, podMap,
  Step,
  stepMap
} from "../definitions/definitions";
import {svgNodeHeight, svgNodeWidth} from "../config/config";
import * as d3 from "d3-selection";

function getNodeLayoutColumnByName(nodeLayout: KubernetesNode[][], name: string): number {
  let columnId: number = null;
  nodeLayout.forEach((column, currentColumnId) => {
    column.forEach(rowElement => {
      if (rowElement.name === name) {
        columnId = currentColumnId;
      }
    })
  });
  return columnId;
}

function getNodeLayoutRowByName(nodeLayout: KubernetesNode[][], name: string): number {
  function isStep(kubernetesNode: KubernetesNode): kubernetesNode is Step {
    return stepMap.get(kubernetesNode.name) !== undefined;
  }

  let rowId: number = null;
  nodeLayout.forEach((column) => {
    let row: number = 0;
    column.forEach((rowElement) => {
      if (rowElement.name === name) {
        rowId = row;
        return;
      }
      if (isStep(rowElement)) {
        row += rowElement.podNames.length;
      } else {
        row++;
      }
    })
  });
  return rowId;
}

export function drawSvg(kubernetesGraph: KubernetesGraph, nodeLayout: KubernetesNode[][]) {
  let dataSourceNodes: D3Node[] = kubernetesGraph.dataSources.map(dataSource => dataSourceMap.get(dataSource.name))
    .map(dataSource => ({
      id: dataSource.name,
      text: dataSource.name, //, labels[' + dataSource.labels.map(label => [label.key, label.value].join(':')).join(' | ') + ']',
      x: 20 + (svgNodeWidth + 150) * getNodeLayoutColumnByName(nodeLayout, dataSource.name),
      y: 20 + 1.50 * svgNodeHeight * getNodeLayoutRowByName(nodeLayout, dataSource.name),
      width: svgNodeWidth,
      height: svgNodeHeight,
      type: 'data-source'
    }));
  let stepNodes: D3Node[] = kubernetesGraph.steps.map(step => stepMap.get(step.name))
    .map((step, i) => ({
      id: step.name,
      text: step.name,
      x: 10 + (svgNodeWidth + 150) * getNodeLayoutColumnByName(nodeLayout, step.name),
      y: 10 + 1.50 * svgNodeHeight * getNodeLayoutRowByName(nodeLayout, step.name),
      width: svgNodeWidth + 20,
      height: Math.max(svgNodeHeight * step.podNames.length + Math.max(svgNodeHeight / 2 * (step.podNames.length - 1), 0), svgNodeHeight) + 20,
      type: 'step'
    }));
  let podNodes: D3Node[] = [];
  stepNodes.forEach(stepNode => {
    let step: Step = stepMap.get(stepNode.id);
    step.podNames.forEach((podName, i) => {
      podNodes.push({
        id: podName,
        text: podName,
        x: stepNode.x + 10,
        y: stepNode.y + (svgNodeHeight + 50) * i + 30,
        width: svgNodeWidth,
        height: svgNodeHeight - 20,
        type: 'pod'
      });
    });
  });

  let nodes: D3Node[] = [...dataSourceNodes, ...stepNodes, ...podNodes];
  let edges: D3Edge[] = [];

  kubernetesGraph.dataSources.filter(dataSource => dataSource.creatorPodName).forEach(dataSource => {
    edges.push({
      start: dataSource.creatorPodName,
      stop: dataSource.name
    });
  });

  kubernetesGraph.pods.filter(pod => pod.creatorDataSourceNames.length > 0).forEach(pod => {
    pod.creatorDataSourceNames.forEach(dataSourceName => {
      edges.push({
        start: dataSourceName,
        stop: pod.name
      });
    });
  });

  const graph = {
    nodes: nodes,
    edges: edges,
    node: function (id) {
      if (!this.nmap) {
        this["nmap"] = {};
        for (let i = 0; i < this.nodes.length; i++) {
          let node = this.nodes[i];
          this.nmap[node.id] = node;
        }
      }
      return this.nmap[id];
    },
    mid: function (id) {
      let node = this.node(id);
      let x = node.width / 2.0 + node.x,
        y = node.height / 2.0 + node.y;
      return {x: x, y: y};
    }
  };

  let g = d3.select('#mysvg')
    .selectAll('g')
    .data(graph.nodes)
    .enter()
    .append('g')
    .attr('id', function (d) {
      return d.id;
    })
    .attr('transform', function (d) {
      return 'translate(' + d.x + ',' + d.y + ')';
    });
  g.append('rect')
    .attr('id', function (d) {
      return d.id;
    })
    .attr('x', 0)
    .attr('y', 0)
    .attr('style', function (d) {
      if (d.type === 'data-source') {
        return 'stroke:#000000; fill:#eeeeee;';
      }
      if (d.type === 'pod') {
        return 'stroke:#000000; fill:#add8e6;';
      }
      return 'stroke:#000000; fill:#ffaa1d;';
    })
    .attr('width', function (d) {
      return d.width;
    })
    .attr('height', function (d) {
      return d.height;
    })
    .attr('pointer-events', 'visible');

  d3.select('#mysvg')
    .selectAll('line')
    .data(graph.edges)
    .enter()
    .insert('line')
    .attr('data-start', function (d) {
      return d.start;
    })
    .attr('data-stop', function (d) {
      return d.stop;
    })
    .attr('x1', function (d) {
      return graph.mid(d.start).x + svgNodeWidth / 2;
    }.bind(this))
    .attr('y1', function (d) {
      return graph.mid(d.start).y;
    })
    .attr('x2', function (d) {
      return graph.mid(d.stop).x - svgNodeWidth / 2;
    }.bind(this))
    .attr('y2', function (d) {
      return graph.mid(d.stop).y
    })
    .attr('style', 'stroke:rgb(80,80,80);stroke-width:2');

  g.append('text')
    .attr('x', 10)
    .attr('y', 10)
    .attr('dy', '.35em')
    .attr('font-size', 'smaller')
    .text(function (d) {
      return d.text;
    });

  document.getElementById('mysvg').setAttribute('width', '20000');
  document.getElementById('mysvg').setAttribute('height', '20000');
}
