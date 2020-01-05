import {Component, HostListener} from '@angular/core';
import * as d3 from 'd3-selection';

declare class dataSourceStepMatch {
  dataSource: dataSource;
  step: step;
}

declare class dataSourceLabelKeyValuePair {
  key: string;
  value: string;
}

declare class stepKeyValuePair {
  regex: boolean;
  key: string;
  value: string;
}

declare class dataSource {
  name: string;
  labels: dataSourceLabelKeyValuePair[];
  depth?: number;
}

declare class step {
  name: string;
  keyValuePairs: stepKeyValuePair[];
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

function dataSourceLabelMatchesStepKeyValuePair(dataSourceLabel: dataSourceLabelKeyValuePair, stepKeyValuePair: stepKeyValuePair) {
  if (stepKeyValuePair.regex) {
    let keyRegex: RegExp = RegExp(stepKeyValuePair.key);
    if (!keyRegex.test(dataSourceLabel.key)) {
      return false;
    }
    let valueRegex: RegExp = RegExp(stepKeyValuePair.value);
    if (!valueRegex.test(dataSourceLabel.value)) {
      return false;
    }
  }
  else {
    if (dataSourceLabel.key !== stepKeyValuePair.key) {
      return false;
    }
    if (dataSourceLabel.value !== stepKeyValuePair.value) {
      return false;
    }
  }
  return true;
}

function stepMatchesDataSource(step: step, dataSource: dataSource) {
  for (let stepKeyValuePair of step.keyValuePairs) {
    let foundMatchingLabel: boolean = false;
    for (let dataSourceLabel of dataSource.labels) {
      if (dataSourceLabelMatchesStepKeyValuePair(dataSourceLabel, stepKeyValuePair)) {
        foundMatchingLabel = true;
        break;
      }
    }
    if (!foundMatchingLabel) {
      return false;
    }
  }
  return true;
}

function getDataSourceStepMatches(dataSources: dataSource[], steps: step[]) {
  let dataSourceStepMatches: dataSourceStepMatch[] = [];
  for (let dataSource of dataSources) {
    for (let step of steps) {
      if (stepMatchesDataSource(step, dataSource)) {
        dataSourceStepMatches.push({dataSource: dataSource, step: step})
      }
    }
  }
  return dataSourceStepMatches;
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'zerops-frontend';

  @HostListener('click', ['$event.target']) onClick(target) {
    if (!target.closest('rect')) return;
    alert(target.id);
  }

  width: number = 200;
  height: number = 100;

  dataSources: dataSource[] = this.getDataSourcesRaw().items.map(dataSourceRaw => ({
    name: dataSourceRaw.metadata.name,
    labels: Object.keys(dataSourceRaw.metadata.labels).map(dataSourceLabelKey => ({
      key: dataSourceLabelKey,
      value: dataSourceRaw.metadata.labels[dataSourceLabelKey]
    }))
  }));
  steps: step[] = this.getStepsRaw().items.map(stepRaw => ({
    name: stepRaw.metadata.name,
    keyValuePairs: stepRaw.spec.ingest.map(ingest => ({
      regex: ingest.check === 'regex',
      key: ingest.key,
      value: ingest.value
    }))
  }));

  dataSourceStepMatches: dataSourceStepMatch[] = getDataSourceStepMatches(this.dataSources, this.steps);

  dataSourcesNodes: d3Node[] = this.dataSources.map((dataSource, i) => ({
    id: 'dataSource:' + dataSource.name,
    text: dataSource.name + ' | ' + dataSource.labels.map(label => [label.key, label.value].join(':')).join(' | '),
    x: 10,
    y: 10 + 1.50 * this.height * i,
    width: this.width,
    height: this.height
  }));
  stepsNodes: d3Node[] = this.steps.map((step, i) => ({
    id: 'step:' + step.name,
    text: step.name,
    x: 160 + this.width,
    y: 10 + 1.50 * this.height * i,
    width: this.width,
    height: this.height
  }));
  nodes: d3Node[] = [...this.dataSourcesNodes, ...this.stepsNodes];
  edges: d3Edge[] = this.dataSourceStepMatches.map(match => ({start: 'dataSource:' + match.dataSource.name, stop: 'step:' + match.step.name}));

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
        return graph.mid(d.start).x + this.width / 2;
      }.bind(this))
      .attr('y1', function (d) {
        return graph.mid(d.start).y;
      })
      .attr('x2', function (d) {
        return graph.mid(d.stop).x - this.width / 2;
      }.bind(this))
      .attr('y2', function (d) {
        return graph.mid(d.stop).y
      })
      .attr('style', 'stroke:rgb(80,80,80);stroke-width:2');

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

  }


  getStepsRaw() {
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
        },
        {
          "apiVersion": "zerops.com/v1",
          "kind": "ZerOpsStep",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsStep\",\"metadata\":{\"annotations\":{},\"name\":\"example-step-matches-a\"},\"spec\":{\"ingest\":[{\"key\":\"name\",\"value\":\"A\"}],\"outputs\":[{\"labels\":{\"data\":\"aggregated-all-physical\"},\"name\":\"phys\",\"url\":\"tcp://:9000\"}],\"template\":{\"metadata\":{\"labels\":{\"app\":\"aggregate-physical-data\"}},\"spec\":{\"containers\":[{\"command\":[\"sh\",\"-c\",\"echo \\\"My IP: $(ip route get 1 | awk '{print $NF;exit}')\\\"\\n/bitflow-pipeline \\\\\\n\\\"{ZEROPS_DATA_SOURCE}\\n-\\u003e csv://:9000\\\"\\n\"],\"image\":\"teambitflow/go-bitflow\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"container\",\"ports\":[{\"containerPort\":9000,\"name\":\"data\"}]}]}},\"type\":\"all-to-one\"}}\n"
            },
            "creationTimestamp": "2019-11-07T13:35:43Z",
            "generation": 1,
            "name": "example-step-matches-a",
            "resourceVersion": "55004",
            "selfLink": "/apis/zerops.com/v1/zerops-steps/example-step-matches-a",
            "uid": "bc970363-05d8-4ebe-b6f5-797060edffa9"
          },
          "spec": {
            "ingest": [
              {
                "key": "name",
                "value": "A"
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
        },
        {
          "apiVersion": "zerops.com/v1",
          "kind": "ZerOpsStep",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsStep\",\"metadata\":{\"annotations\":{},\"name\":\"example-step-matches-ab\"},\"spec\":{\"ingest\":[{\"check\":\"regex\",\"key\":\"name\",\"value\":\"^(A|B)$\"}],\"outputs\":[{\"labels\":{\"data\":\"aggregated-all-physical\"},\"name\":\"phys\",\"url\":\"tcp://:9000\"}],\"template\":{\"metadata\":{\"labels\":{\"app\":\"aggregate-physical-data\"}},\"spec\":{\"containers\":[{\"command\":[\"sh\",\"-c\",\"echo \\\"My IP: $(ip route get 1 | awk '{print $NF;exit}')\\\"\\n/bitflow-pipeline \\\\\\n\\\"{ZEROPS_DATA_SOURCE}\\n-\\u003e csv://:9000\\\"\\n\"],\"image\":\"teambitflow/go-bitflow\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"container\",\"ports\":[{\"containerPort\":9000,\"name\":\"data\"}]}]}},\"type\":\"all-to-one\"}}\n"
            },
            "creationTimestamp": "2019-11-07T13:35:48Z",
            "generation": 1,
            "name": "example-step-matches-ab",
            "resourceVersion": "55010",
            "selfLink": "/apis/zerops.com/v1/zerops-steps/example-step-matches-ab",
            "uid": "a3eaaa0a-46f2-46bf-8604-220542d4a79d"
          },
          "spec": {
            "ingest": [
              {
                "check": "regex",
                "key": "name",
                "value": "^(A|B)$"
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
        },
        {
          "apiVersion": "zerops.com/v1",
          "kind": "ZerOpsStep",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsStep\",\"metadata\":{\"annotations\":{},\"name\":\"example-step-matches-bc\"},\"spec\":{\"ingest\":[{\"check\":\"regex\",\"key\":\"name\",\"value\":\"^(B|C)$\"}],\"outputs\":[{\"labels\":{\"data\":\"aggregated-all-physical\"},\"name\":\"phys\",\"url\":\"tcp://:9000\"}],\"template\":{\"metadata\":{\"labels\":{\"app\":\"aggregate-physical-data\"}},\"spec\":{\"containers\":[{\"command\":[\"sh\",\"-c\",\"echo \\\"My IP: $(ip route get 1 | awk '{print $NF;exit}')\\\"\\n/bitflow-pipeline \\\\\\n\\\"{ZEROPS_DATA_SOURCE}\\n-\\u003e csv://:9000\\\"\\n\"],\"image\":\"teambitflow/go-bitflow\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"container\",\"ports\":[{\"containerPort\":9000,\"name\":\"data\"}]}]}},\"type\":\"all-to-one\"}}\n"
            },
            "creationTimestamp": "2019-11-07T13:35:51Z",
            "generation": 1,
            "name": "example-step-matches-bc",
            "resourceVersion": "55017",
            "selfLink": "/apis/zerops.com/v1/zerops-steps/example-step-matches-bc",
            "uid": "f7ba7bef-936d-495a-ad05-87657fb5057a"
          },
          "spec": {
            "ingest": [
              {
                "check": "regex",
                "key": "name",
                "value": "^(B|C)$"
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


  getDataSourcesRaw() {
    return {
      "apiVersion": "v1",
      "items": [
        {
          "apiVersion": "zerops.com/v1",
          "kind": "ZerOpsDataSource",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsDataSource\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"A\"},\"name\":\"data-source-a\"},\"spec\":{\"url\":\"tcp://172.17.0.8:9000\"}}\n"
            },
            "creationTimestamp": "2019-11-07T13:34:46Z",
            "generation": 1,
            "labels": {
              "name": "A"
            },
            "name": "data-source-a",
            "resourceVersion": "54932",
            "selfLink": "/apis/zerops.com/v1/zerops-data-sources/data-source-a",
            "uid": "1d724749-13a0-4a6d-9cbe-b3f3fa304030"
          },
          "spec": {
            "url": "tcp://172.17.0.8:9000"
          }
        },
        {
          "apiVersion": "zerops.com/v1",
          "kind": "ZerOpsDataSource",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsDataSource\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"B\"},\"name\":\"data-source-b\"},\"spec\":{\"url\":\"tcp://172.17.0.8:9000\"}}\n"
            },
            "creationTimestamp": "2019-11-07T13:35:27Z",
            "generation": 1,
            "labels": {
              "name": "B"
            },
            "name": "data-source-b",
            "resourceVersion": "54983",
            "selfLink": "/apis/zerops.com/v1/zerops-data-sources/data-source-b",
            "uid": "cff92c32-ed10-4c8d-a7ac-ab95cd569a1a"
          },
          "spec": {
            "url": "tcp://172.17.0.8:9000"
          }
        },
        {
          "apiVersion": "zerops.com/v1",
          "kind": "ZerOpsDataSource",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsDataSource\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"C\"},\"name\":\"data-source-c\"},\"spec\":{\"url\":\"tcp://172.17.0.8:9000\"}}\n"
            },
            "creationTimestamp": "2019-11-07T13:35:30Z",
            "generation": 1,
            "labels": {
              "name": "C"
            },
            "name": "data-source-c",
            "resourceVersion": "54988",
            "selfLink": "/apis/zerops.com/v1/zerops-data-sources/data-source-c",
            "uid": "2cc4438a-3b11-4324-8d95-8816399b3328"
          },
          "spec": {
            "url": "tcp://172.17.0.8:9000"
          }
        },
        {
          "apiVersion": "zerops.com/v1",
          "kind": "ZerOpsDataSource",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsDataSource\",\"metadata\":{\"annotations\":{},\"labels\":{\"collector\":\"bitflow\",\"host\":\"wally166\",\"layer\":\"physical\"},\"name\":\"wally166-collector\"},\"spec\":{\"url\":\"tcp://172.17.0.8:9000\"}}\n"
            },
            "creationTimestamp": "2019-11-03T17:24:37Z",
            "generation": 1,
            "labels": {
              "collector": "bitflow",
              "host": "wally166",
              "layer": "physical"
            },
            "name": "wally166-collector",
            "resourceVersion": "8283",
            "selfLink": "/apis/zerops.com/v1/zerops-data-sources/wally166-collector",
            "uid": "ec285619-1b02-40af-8187-c6bd75d0d5a6"
          },
          "spec": {
            "url": "tcp://172.17.0.8:9000"
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


  // getSteps() {
  //   return {
  //     "apiVersion": "v1",
  //     "items": [
  //       {
  //         "apiVersion": "zerops.com/v1",
  //         "kind": "ZerOpsStep",
  //         "metadata": {
  //           "annotations": {
  //             "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsStep\",\"metadata\":{\"annotations\":{},\"name\":\"example-analysis-all-to-one\"},\"spec\":{\"ingest\":[{\"key\":\"collector\",\"value\":\"bitflow\"},{\"check\":\"regex\",\"key\":\"^layer$\",\"value\":\"^physical$\"},{\"check\":\"regex\",\"key\":\"host\",\"value\":\"wally1.*\"}],\"outputs\":[{\"labels\":{\"data\":\"aggregated-all-physical\"},\"name\":\"phys\",\"url\":\"tcp://:9000\"}],\"template\":{\"metadata\":{\"labels\":{\"app\":\"aggregate-physical-data\"}},\"spec\":{\"containers\":[{\"command\":[\"sh\",\"-c\",\"echo \\\"My IP: $(ip route get 1 | awk '{print $NF;exit}')\\\"\\n/bitflow-pipeline \\\\\\n\\\"{ZEROPS_DATA_SOURCE}\\n-\\u003e csv://:9000\\\"\\n\"],\"image\":\"teambitflow/go-bitflow\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"container\",\"ports\":[{\"containerPort\":9000,\"name\":\"data\"}]}]}},\"type\":\"all-to-one\"}}\n"
  //           },
  //           "creationTimestamp": "2019-11-03T17:24:46Z",
  //           "generation": 1,
  //           "name": "example-analysis-all-to-one",
  //           "resourceVersion": "8295",
  //           "selfLink": "/apis/zerops.com/v1/zerops-steps/example-analysis-all-to-one",
  //           "uid": "3cf37727-5651-434a-8ae0-15df62fc5086"
  //         },
  //         "spec": {
  //           "ingest": [
  //             {
  //               "key": "collector",
  //               "value": "bitflow"
  //             },
  //             {
  //               "check": "regex",
  //               "key": "^layer$",
  //               "value": "^physical$"
  //             },
  //             {
  //               "check": "regex",
  //               "key": "host",
  //               "value": "wally1.*"
  //             }
  //           ],
  //           "outputs": [
  //             {
  //               "labels": {
  //                 "data": "aggregated-all-physical"
  //               },
  //               "name": "phys",
  //               "url": "tcp://:9000"
  //             }
  //           ],
  //           "template": {
  //             "metadata": {
  //               "labels": {
  //                 "app": "aggregate-physical-data"
  //               }
  //             },
  //             "spec": {
  //               "containers": [
  //                 {
  //                   "command": [
  //                     "sh",
  //                     "-c",
  //                     "echo \"My IP: $(ip route get 1 | awk '{print $NF;exit}')\"\n/bitflow-pipeline \\\n\"{ZEROPS_DATA_SOURCE}\n-\u003e csv://:9000\"\n"
  //                   ],
  //                   "image": "teambitflow/go-bitflow",
  //                   "imagePullPolicy": "IfNotPresent",
  //                   "name": "container",
  //                   "ports": [
  //                     {
  //                       "containerPort": 9000,
  //                       "name": "data"
  //                     }
  //                   ]
  //                 }
  //               ]
  //             }
  //           },
  //           "type": "all-to-one"
  //         }
  //       }
  //     ],
  //     "kind": "List",
  //     "metadata": {
  //       "resourceVersion": "",
  //       "selfLink": ""
  //     }
  //   };
  // }


  // getDataSources() {
  //   return {
  //     "apiVersion": "v1"
  //     ,
  //     "items": [
  //       {
  //         "apiVersion": "zerops.com/v1",
  //         "kind": "ZerOpsStep",
  //         "metadata": {
  //           "annotations": {
  //             "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsStep\",\"metadata\":{\"annotations\":{},\"name\":\"example-analysis-all-to-one\"},\"spec\":{\"ingest\":[{\"key\":\"collector\",\"value\":\"bitflow\"},{\"check\":\"regex\",\"key\":\"^layer$\",\"value\":\"^physical$\"},{\"check\":\"regex\",\"key\":\"host\",\"value\":\"wally1.*\"}],\"outputs\":[{\"labels\":{\"data\":\"aggregated-all-physical\"},\"name\":\"phys\",\"url\":\"tcp://:9000\"}],\"template\":{\"metadata\":{\"labels\":{\"app\":\"aggregate-physical-data\"}},\"spec\":{\"containers\":[{\"command\":[\"sh\",\"-c\",\"echo \\\"My IP: $(ip route get 1 | awk '{print $NF;exit}')\\\"\\n/bitflow-pipeline \\\\\\n\\\"{ZEROPS_DATA_SOURCE}\\n-\\u003e csv://:9000\\\"\\n\"],\"image\":\"teambitflow/go-bitflow\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"container\",\"ports\":[{\"containerPort\":9000,\"name\":\"data\"}]}]}},\"type\":\"all-to-one\"}}\n"
  //           },
  //           "creationTimestamp": "2019-11-03T17:24:46Z",
  //           "generation": 1,
  //           "name": "example-analysis-all-to-one",
  //           "resourceVersion": "8295",
  //           "selfLink": "/apis/zerops.com/v1/zerops-steps/example-analysis-all-to-one",
  //           "uid": "3cf37727-5651-434a-8ae0-15df62fc5086"
  //         },
  //         "spec": {
  //           "ingest": [
  //             {
  //               "key": "collector",
  //               "value": "bitflow"
  //             },
  //             {
  //               "check": "regex",
  //               "key": "^layer$",
  //               "value": "^physical$"
  //             },
  //             {
  //               "check": "regex",
  //               "key": "host",
  //               "value": "wally1.*"
  //             },
  //             {
  //               "key": "unusedKey",
  //               "value": "unusedValue"
  //             }
  //           ],
  //           "outputs": [
  //             {
  //               "labels": {
  //                 "data": "aggregated-all-physical"
  //               },
  //               "name": "phys",
  //               "url": "tcp://:9000"
  //             }
  //           ],
  //           "template": {
  //             "metadata": {
  //               "labels": {
  //                 "app": "aggregate-physical-data"
  //               }
  //             },
  //             "spec": {
  //               "containers": [
  //                 {
  //                   "command": [
  //                     "sh",
  //                     "-c",
  //                     "echo \"My IP: $(ip route get 1 | awk '{print $NF;exit}')\"\n/bitflow-pipeline \\\n\"{ZEROPS_DATA_SOURCE}\n-\u003e csv://:9000\"\n"
  //                   ],
  //                   "image": "teambitflow/go-bitflow",
  //                   "imagePullPolicy": "IfNotPresent",
  //                   "name": "container",
  //                   "ports": [
  //                     {
  //                       "containerPort": 9000,
  //                       "name": "data"
  //                     }
  //                   ]
  //                 }
  //               ]
  //             }
  //           },
  //           "type": "all-to-one"
  //         }
  //       }
  //     ]
  //     ,
  //     "kind": "List"
  //     ,
  //     "metadata": {
  //       "resourceVersion": "",
  //       "selfLink": ""
  //     }
  //   };
  // }
}
