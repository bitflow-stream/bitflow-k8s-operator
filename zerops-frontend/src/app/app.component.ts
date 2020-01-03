import {Component} from '@angular/core';
import * as d3 from 'd3-selection';

declare class dataSource {
  key: string;
  value: string;
}

declare class step {
  key: string;
  value: string;
}

declare class d3Node {
  id: string;
  text: string;
  x: number;
  y: number;
  width: number;
  height: number;
}

declare class d3Edge {
  start: string;
  stop: string;
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'zerops-frontend';


  dataSources: dataSource[] = this.getDataSources().items[0].spec.ingest;
  steps: step[] = this.getSteps().items[0].spec.ingest;
  dataSourcesNodes: d3Node[] = this.dataSources.map((dataSource, i) => ({id: 'dataSource:' + dataSource.key, text: dataSource.key + ':' + dataSource.value, x: 10, y: 10 + 150 * i, width: 150, height: 100}));
  stepsNodes: d3Node[] = this.steps.map((step, i) => ({id: 'step:' + step.key, text: step.key + ':' + step.value, x: 210, y: 10 + 150 * i, width: 150, height: 100}));
  nodes: d3Node[] = [...this.dataSourcesNodes, ...this.stepsNodes];
  edges: d3Edge[] = [{start: 'dataSource:collector', stop: 'step:^layer$'}];

  ngAfterContentInit() {
    d3.select('#mysvg');

    const graph = {
      nodes: this.nodes,
      edges: this.edges,
      node: function (id) {
        if (!this.nmap) {
          this.nmap = {};
          for (var i = 0; i < this.nodes.length; i++) {
            var node = this.nodes[i];
            this.nmap[node.id] = node;
          }
        }
        return this.nmap[id];
      },
      mid: function (id) {
        var node = this.node(id);
        var x = node.width / 2.0 + node.x,
          y = node.height / 2.0 + node.y;
        return {x: x, y: y};
      }
    };

    const arcs = d3.select('#mysvg')
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
        return graph.mid(d.start).x;
      })
      .attr('y1', function (d) {
        return graph.mid(d.start).y;
      })
      .attr('x2', function (d) {
        return graph.mid(d.stop).x;
      })
      .attr('y2', function (d) {
        return graph.mid(d.stop).y
      })
      .attr('style', 'stroke:rgb(255,0,0);stroke-width:2')
      .attr('marker-end', 'url(#arrow)');

    var g = d3.select('#mysvg')
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
      .attr('style', 'stroke:#000000; fill:white;')
      .attr('width', function (d) {
          return d.width;
        })
      .attr('height', function (d) {
          return d.height;
        })
      .attr('pointer-events', 'visible')
      ;
    g.append('text')
      .attr("x", 10)
      .attr("y", 10)
      .attr("dy", ".35em")
      .text(function (d) {
        return d.text;
      });
  }


  getSteps() {
    return {
      "apiVersion": "v1",
      "items": [
        {
          "apiVersion": "zerops.com/v1",
          "kind": "ZerOpsStep",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsStep\",\"metadata\":{\"annotations\":{},\"name\":\"example-analysis-all-to-one\"},\"spec\":{\"ingest\":[{\"key\":\"collector\",\"value\":\"bitflow\"},{\"check\":\"regex\",\"key\":\"^layer$\",\"value\":\"^physical$\"},{\"check\":\"regex\",\"key\":\"host\",\"value\":\"wally1.*\"}],\"outputs\":[{\"labels\":{\"data\":\"aggregated-all-physical\"},\"name\":\"phys\",\"url\":\"tcp://:9000\"}],\"template\":{\"metadata\":{\"labels\":{\"app\":\"aggregate-physical-data\"}},\"spec\":{\"containers\":[{\"command\":[\"sh\",\"-c\",\"echo \\\"My IP: $(ip route get 1 | awk '{print $NF;exit}')\\\"\\n/bitflow-pipeline \\\\\\n\\\"{ZEROPS_DATA_SOURCE}\\n-\\u003e csv://:9000\\\"\\n\"],\"image\":\"teambitflow/go-bitflow\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"container\",\"ports\":[{\"containerPort\":9000,\"name\":\"data\"}]}]}},\"type\":\"all-to-one\"}}\n"
            },
            "creationTimestamp": "2019-11-03T17:24:46Z",
            "generation": 1,
            "name": "example-analysis-all-to-one",
            "resourceVersion": "8295",
            "selfLink": "/apis/zerops.com/v1/zerops-steps/example-analysis-all-to-one",
            "uid": "3cf37727-5651-434a-8ae0-15df62fc5086"
          },
          "spec": {
            "ingest": [
              {
                "key": "collector",
                "value": "bitflow"
              },
              {
                "check": "regex",
                "key": "^layer$",
                "value": "^physical$"
              },
              {
                "check": "regex",
                "key": "host",
                "value": "wally1.*"
              }
            ],
            "outputs": [
              {
                "labels": {
                  "data": "aggregated-all-physical"
                },
                "name": "phys",
                "url": "tcp://:9000"
              }
            ],
            "template": {
              "metadata": {
                "labels": {
                  "app": "aggregate-physical-data"
                }
              },
              "spec": {
                "containers": [
                  {
                    "command": [
                      "sh",
                      "-c",
                      "echo \"My IP: $(ip route get 1 | awk '{print $NF;exit}')\"\n/bitflow-pipeline \\\n\"{ZEROPS_DATA_SOURCE}\n-\u003e csv://:9000\"\n"
                    ],
                    "image": "teambitflow/go-bitflow",
                    "imagePullPolicy": "IfNotPresent",
                    "name": "container",
                    "ports": [
                      {
                        "containerPort": 9000,
                        "name": "data"
                      }
                    ]
                  }
                ]
              }
            },
            "type": "all-to-one"
          }
        }
      ],
      "kind": "List",
      "metadata": {
        "resourceVersion": "",
        "selfLink": ""
      }
    };
  }


  getDataSources() {
    return {
      "apiVersion": "v1"
      ,
      "items": [
        {
          "apiVersion": "zerops.com/v1",
          "kind": "ZerOpsStep",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsStep\",\"metadata\":{\"annotations\":{},\"name\":\"example-analysis-all-to-one\"},\"spec\":{\"ingest\":[{\"key\":\"collector\",\"value\":\"bitflow\"},{\"check\":\"regex\",\"key\":\"^layer$\",\"value\":\"^physical$\"},{\"check\":\"regex\",\"key\":\"host\",\"value\":\"wally1.*\"}],\"outputs\":[{\"labels\":{\"data\":\"aggregated-all-physical\"},\"name\":\"phys\",\"url\":\"tcp://:9000\"}],\"template\":{\"metadata\":{\"labels\":{\"app\":\"aggregate-physical-data\"}},\"spec\":{\"containers\":[{\"command\":[\"sh\",\"-c\",\"echo \\\"My IP: $(ip route get 1 | awk '{print $NF;exit}')\\\"\\n/bitflow-pipeline \\\\\\n\\\"{ZEROPS_DATA_SOURCE}\\n-\\u003e csv://:9000\\\"\\n\"],\"image\":\"teambitflow/go-bitflow\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"container\",\"ports\":[{\"containerPort\":9000,\"name\":\"data\"}]}]}},\"type\":\"all-to-one\"}}\n"
            },
            "creationTimestamp": "2019-11-03T17:24:46Z",
            "generation": 1,
            "name": "example-analysis-all-to-one",
            "resourceVersion": "8295",
            "selfLink": "/apis/zerops.com/v1/zerops-steps/example-analysis-all-to-one",
            "uid": "3cf37727-5651-434a-8ae0-15df62fc5086"
          },
          "spec": {
            "ingest": [
              {
                "key": "collector",
                "value": "bitflow"
              },
              {
                "check": "regex",
                "key": "^layer$",
                "value": "^physical$"
              },
              {
                "check": "regex",
                "key": "host",
                "value": "wally1.*"
              },
              {
                "key": "unusedKey",
                "value": "unusedValue"
              }
            ],
            "outputs": [
              {
                "labels": {
                  "data": "aggregated-all-physical"
                },
                "name": "phys",
                "url": "tcp://:9000"
              }
            ],
            "template": {
              "metadata": {
                "labels": {
                  "app": "aggregate-physical-data"
                }
              },
              "spec": {
                "containers": [
                  {
                    "command": [
                      "sh",
                      "-c",
                      "echo \"My IP: $(ip route get 1 | awk '{print $NF;exit}')\"\n/bitflow-pipeline \\\n\"{ZEROPS_DATA_SOURCE}\n-\u003e csv://:9000\"\n"
                    ],
                    "image": "teambitflow/go-bitflow",
                    "imagePullPolicy": "IfNotPresent",
                    "name": "container",
                    "ports": [
                      {
                        "containerPort": 9000,
                        "name": "data"
                      }
                    ]
                  }
                ]
              }
            },
            "type": "all-to-one"
          }
        }
      ]
      ,
      "kind": "List"
      ,
      "metadata": {
        "resourceVersion": "",
        "selfLink": ""
      }
    };
  }
}
