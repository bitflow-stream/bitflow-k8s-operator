import {D3Edge, D3Node, stepMap, dataSourceMap, podMap} from "../definitions/definitions";
import {svgNodeHeight, svgNodeWidth} from "../config/config";
import * as d3 from "d3-selection";

function getNodeLayoutColumnByUuid(nodeLayout: string[][], uuid: string): number {
  let columnId: number = null;
  nodeLayout.forEach((column, currentColumnId) => {
    column.forEach(rowElement => {
      if (rowElement === uuid) {
        columnId = currentColumnId;
      }
    })
  });
  return columnId;
}

function getNodeLayoutRowByUuid(nodeLayout: string[][], uuid: string): number {
  let rowId: number = null;
  nodeLayout.forEach((column) => {
    column.forEach((rowElement, currentRowId) => {
      if (rowElement === uuid) {
        rowId = currentRowId;
      }
    })
  });
  return rowId;
}

export function drawSvg(nodeLayout: string[][]) {
  // let dataSourcesNodes: D3Node[] = kubernetesGraph.dataSources.map(dataSourceGraphElement => dataSourceMap.get(dataSourceGraphElement.uuid))
  //   .map(dataSource => ({
  //     id: dataSource.uuid,
  //     text: 'name[' + dataSource.name + '], labels[' + dataSource.labels.map(label => [label.key, label.value].join(':')).join(' | ') + ']',
  //     x: 10 + (svgNodeWidth + 150) * getNodeLayoutColumnByUuid(nodeLayout, dataSource.uuid),
  //     y: 10 + 1.50 * svgNodeHeight * getNodeLayoutRowByUuid(nodeLayout, dataSource.uuid),
  //     width: svgNodeWidth,
  //     height: svgNodeHeight,
  //     type: 'data-source'
  //   }));
  // let stepsNodes: D3Node[] = kubernetesGraph.steps.map(stepGraphElement => stepMap.get(stepGraphElement.uuid))
  //   .map((step, i) => ({
  //     id: step.uuid,
  //     text: step.name,
  //     x: 10 + (svgNodeWidth + 150) * getNodeLayoutColumnByUuid(nodeLayout, step.uuid),
  //     y: 10 + 1.50 * svgNodeHeight * getNodeLayoutRowByUuid(nodeLayout, step.uuid),
  //     width: svgNodeWidth,
  //     height: svgNodeHeight,
  //     type: 'step'
  //   }));
  // let nodes: D3Node[] = [...dataSourcesNodes, ...stepsNodes];
  // let edges: D3Edge[] = [];
  //
  // kubernetesGraph.dataSources.forEach(dataSourceGraphElement => {
  //   dataSourceGraphElement.stepGraphElements.forEach(stepGraphElement => {
  //     edges.push({
  //       start: dataSourceGraphElement.uuid,
  //       stop: stepGraphElement
  //     });
  //   });
  // });
  // kubernetesGraph.steps.forEach(stepGraphElement => {
  //   stepGraphElement.outputDataSourceGraphElements.forEach(outputDataSourceGraphElement => {
  //     edges.push({
  //       start: stepGraphElement.uuid,
  //       stop: outputDataSourceGraphElement
  //     });
  //   });
  // });
  //
  // const graph = {
  //   nodes: nodes,
  //   edges: edges,
  //   node: function (id) {
  //     if (!this.nmap) {
  //       this["nmap"] = {};
  //       for (let i = 0; i < this.nodes.length; i++) {
  //         let node = this.nodes[i];
  //         this.nmap[node.id] = node;
  //       }
  //     }
  //     return this.nmap[id];
  //   },
  //   mid: function (id) {
  //     let node = this.node(id);
  //     let x = node.width / 2.0 + node.x,
  //       y = node.height / 2.0 + node.y;
  //     return {x: x, y: y};
  //   }
  // };
  //
  // d3.select('#mysvg')
  //   .selectAll('line')
  //   .data(graph.edges)
  //   .enter()
  //   .insert('line')
  //   .attr('data-start', function (d) {
  //     return d.start;
  //   })
  //   .attr('data-stop', function (d) {
  //     return d.stop;
  //   })
  //   .attr('x1', function (d) {
  //     return graph.mid(d.start).x + svgNodeWidth / 2;
  //   }.bind(this))
  //   .attr('y1', function (d) {
  //     return graph.mid(d.start).y;
  //   })
  //   .attr('x2', function (d) {
  //     return graph.mid(d.stop).x - svgNodeWidth / 2;
  //   }.bind(this))
  //   .attr('y2', function (d) {
  //     return graph.mid(d.stop).y
  //   })
  //   .attr('style', 'stroke:rgb(80,80,80);stroke-width:2');
  //
  // let g = d3.select('#mysvg')
  //   .selectAll('g')
  //   .data(graph.nodes)
  //   .enter()
  //   .append('g')
  //   .attr('id', function (d) {
  //     return d.id;
  //   })
  //   .attr('transform', function (d) {
  //     return 'translate(' + d.x + ',' + d.y + ')';
  //   });
  // g.append('rect')
  //   .attr('id', function (d) {
  //     return d.id;
  //   })
  //   .attr('x', 0)
  //   .attr('y', 0)
  //   .attr('style', function(d) {
  //     if (d.type === 'data-source') {
  //       return 'stroke:#000000; fill:#eeeeee;';
  //     }
  //     return 'stroke:#000000; fill:#ffaa1d;';
  //   })
  //   .attr('width', function (d) {
  //     return d.width;
  //   })
  //   .attr('height', function (d) {
  //     return d.height;
  //   })
  //   .attr('pointer-events', 'visible');
  // g.append('text')
  //   .attr('x', 10)
  //   .attr('y', 10)
  //   .attr('dy', '.35em')
  //   .attr('font-size', 'smaller')
  //   .text(function (d) {
  //     return d.text;
  //   });
  //
  // document.getElementById('mysvg').setAttribute('width', '20000');
  // document.getElementById('mysvg').setAttribute('height', '20000');
}
