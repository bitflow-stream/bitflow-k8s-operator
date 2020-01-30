import {VisualizationData} from "../definitions/definitions";
import {svgNodeWidth} from "../config/config";
import * as d3 from "d3-selection";

// function getNodeLayoutColumnByName(nodeLayout: KubernetesNode[][], name: string): number {
//   let columnId: number = null;
//   nodeLayout.forEach((column, currentColumnId) => {
//     column.forEach(rowElement => {
//       if (rowElement.name === name) {
//         columnId = currentColumnId;
//       }
//     })
//   });
//   return columnId;
// }
//
// function getNodeLayoutRowByName(nodeLayout: KubernetesNode[][], name: string): number {
//   function isStep(kubernetesNode: KubernetesNode): kubernetesNode is Step {
//     return stepMap.get(kubernetesNode.name) !== undefined;
//   }
//
//   let rowId: number = null;
//   nodeLayout.forEach((column) => {
//     let row: number = 0;
//     column.forEach((rowElement) => {
//       if (rowElement.name === name) {
//         rowId = row;
//         return;
//       }
//       if (isStep(rowElement)) {
//         row += rowElement.podNames.length;
//       } else {
//         row++;
//       }
//     })
//   });
//   return rowId;
// }

export function drawSvg(this: any, visualization: VisualizationData) {
  const graph: any = {
    nodes: visualization.nodes,
    edges: visualization.edges,
    node: function (id: any): any {
      if (!this.nmap) {
        this["nmap"] = {};
        for (let i = 0; i < this.nodes.length; i++) {
          let node = this.nodes[i];
          this.nmap[node.id] = node;
        }
      }
      return this.nmap[id];
    },
    mid: function (id: any) {
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
    .attr('id', function (d: any) {
      return d.id;
    })
    .attr('transform', function (d: any) {
      return 'translate(' + d.x + ',' + d.y + ')';
    });
  g.append('rect')
    .attr('id', function (d: any) {
      return d.id;
    })
    .attr('x', 0)
    .attr('y', 0)
    .attr('style', function (d: any) {
      if (d.type === 'data-source') {
        return 'stroke:#000000; fill:#eeeeee;';
      }
      if (d.type === 'pod') {
        return 'stroke:#000000; fill:#add8e6;';
      }
      return 'stroke:#000000; fill:#ffaa1d;';
    })
    .attr('width', function (d: any) {
      return d.width;
    })
    .attr('height', function (d: any) {
      return d.height;
    })
    .attr('pointer-events', 'visible');

  d3.select('#mysvg')
    .selectAll('line')
    .data(graph.edges)
    .enter()
    .insert('line')
    .attr('data-start', function (d: any) {
      return d.start;
    })
    .attr('data-stop', function (d: any) {
      return d.stop;
    })
    .attr('x1', function (d: any) {
      return graph.mid(d.start).x + svgNodeWidth / 2;
    }.bind(this))
    .attr('y1', function (d: any) {
      return graph.mid(d.start).y;
    })
    .attr('x2', function (d: any) {
      return graph.mid(d.stop).x - svgNodeWidth / 2;
    }.bind(this))
    .attr('y2', function (d: any) {
      return graph.mid(d.stop).y
    })
    .attr('style', 'stroke:rgb(80,80,80);stroke-width:2');

  g.append('text')
    .attr('x', 10)
    .attr('y', 10)
    .attr('dy', '.35em')
    .attr('font-size', 'smaller')
    .text(function (d: any) {
      return d.text;
    });

  document.getElementById('mysvg')?.setAttribute('width', '20000');
  document.getElementById('mysvg')?.setAttribute('height', '20000');
}
