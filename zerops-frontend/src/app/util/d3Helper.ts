import {D3Edge, D3Node, dataSourceMap, KubernetesGraph, stepMap} from "../definitions/definitions";
import {svgNodeHeight, svgNodeWidth} from "../config/config";
import * as d3 from "d3-selection";

export function drawSvg(kubernetesGraph: KubernetesGraph) {
  let dataSourcesNodes: D3Node[] = kubernetesGraph.dataSourceGraphElements.map(dataSourceGraphElement => dataSourceMap.get(dataSourceGraphElement.uuid))
    .sort((a, b) => {
      if (a.depth < b.depth) {
        return -1;
      }
      if (a.depth > b.depth) {
        return 1;
      }
      return 0;
    })
    .map((dataSource, i) => ({
      id: dataSource.uuid,
      text: dataSource.name + ' | ' + dataSource.labels.map(label => [label.key, label.value].join(':')).join(' | '),
      x: 10 + ((svgNodeWidth + 150) * 2) * dataSource.depth,
      y: 10 + 1.50 * svgNodeHeight * i,
      width: svgNodeWidth,
      height: svgNodeHeight
    }));
  let stepsNodes: D3Node[] = kubernetesGraph.stepGraphElements.map(stepGraphElement => stepMap.get(stepGraphElement.uuid))
    .map((step, i) => ({
      id: step.uuid,
      text: step.name,
      x: 160 + svgNodeWidth,
      y: 10 + 1.50 * svgNodeHeight * i,
      width: svgNodeWidth,
      height: svgNodeHeight
    }));
  let nodes: D3Node[] = [...dataSourcesNodes, ...stepsNodes];
  let edges: D3Edge[] = [];

  kubernetesGraph.dataSourceGraphElements.forEach(dataSourceGraphElement => {
    dataSourceGraphElement.stepGraphElements.forEach(stepGraphElement => {
      edges.push({
        start: dataSourceGraphElement.uuid,
        stop: stepGraphElement
      });
    });
  });
  kubernetesGraph.stepGraphElements.forEach(stepGraphElement => {
    stepGraphElement.outputDataSourceGraphElements.forEach(outputDataSourceGraphElement => {
      edges.push({
        start: stepGraphElement.uuid,
        stop: outputDataSourceGraphElement
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
    .attr('style', 'stroke:#000000; fill:#eeeeee;')
    .attr('width', function (d) {
      return d.width;
    })
    .attr('height', function (d) {
      return d.height;
    })
    .attr('pointer-events', 'visible');
  g.append('text')
    .attr('x', 10)
    .attr('y', 10)
    .attr('dy', '.35em')
    .attr('font-size', 'smaller')
    .text(function (d) {
      return d.text;
    });

  document.getElementById('mysvg').setAttribute('width', '2000');
  document.getElementById('mysvg').setAttribute('height', '20000');
}
