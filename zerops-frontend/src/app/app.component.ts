import {Component, HostListener} from '@angular/core';
import * as d3 from 'd3-selection';

let dataSourceMap: Map<string, DataSource> = new Map();
let stepMap: Map<string, Step> = new Map();

// let amountOfNodesByDepth: number[] = [];

enum AnalysisType {
  ALL_TO_ONE = 'all-to-one',
  ONE_TO_ONE = 'one-to-one'
}

declare class KubernetesGraph {
  dataSourceGraphElements: DataSourceGraphElement[];
  stepGraphElements: StepGraphElement[];
}

declare class DataSourceGraphElement {
  uuid: string;
  stepGraphElements: string[];
  creatorStepGraphElement: string;
}

declare class StepGraphElement {
  uuid: string;
  outputDataSourceGraphElements: string[];
  sourceDataSourceGraphElements: string[];
}

declare class DataSourcesAndStepsAndMatches {
  dataSources: DataSource[];
  steps: Step[];
  matches: DataSourceStepMatch[];
}

declare class DataSourceStepMatch {
  dataSource: string;
  step: string;
}

declare class DataSourceLabelKeyValuePair {
  key: string;
  value: string;
}

declare class StepKeyValuePair {
  regex: boolean;
  key: string;
  value: string;
}

declare class KeyValuePair {
  key: string;
  value: string;
}

declare class DataSource {
  uuid: string;
  name: string;
  labels: DataSourceLabelKeyValuePair[];
  depth;
}

declare class Step {
  uuid: string;
  name: string;
  keyValuePairs: StepKeyValuePair[];
  type: string;
  outputLabelsArray: KeyValuePair[][];
}

declare class D3Node {
  id: string;
  text: string;
  x: number;
  y: number;
  width: number;
  height: number;
}

declare class D3Edge {
  start: string;
  stop: string;
}

function uuidv4(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

function dataSourceLabelMatchesStepKeyValuePair(dataSourceLabel: DataSourceLabelKeyValuePair, stepKeyValuePair: StepKeyValuePair): boolean {
  if (stepKeyValuePair.regex) {
    let keyRegex: RegExp = RegExp(stepKeyValuePair.key);
    if (!keyRegex.test(dataSourceLabel.key)) {
      return false;
    }
    let valueRegex: RegExp = RegExp(stepKeyValuePair.value);
    if (!valueRegex.test(dataSourceLabel.value)) {
      return false;
    }
  } else {
    if (dataSourceLabel.key !== stepKeyValuePair.key) {
      return false;
    }
    if (dataSourceLabel.value !== stepKeyValuePair.value) {
      return false;
    }
  }
  return true;
}

function stepMatchesDataSource(step: Step, dataSource: DataSource): boolean {
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

function getDataSourceStepMatches(dataSources: string[], steps: string[]): DataSourceStepMatch[] {
  let dataSourceStepMatches: DataSourceStepMatch[] = [];
  for (let dataSource of dataSources) {
    for (let step of steps) {
      if (stepMatchesDataSource(stepMap.get(step), dataSourceMap.get(dataSource))) {
        dataSourceStepMatches.push({dataSource: dataSource, step: step})
      }
    }
  }
  return dataSourceStepMatches;
}

function createKubernetesGraph(dataSources: DataSource[], steps: Step[]): KubernetesGraph {
  let dataSourceGraphElements: DataSourceGraphElement[] = [];
  dataSources.forEach(dataSource => dataSourceGraphElements.push({
    uuid: dataSource.uuid,
    stepGraphElements: [],
    creatorStepGraphElement: null
  }));

  let stepGraphElements: StepGraphElement[] = [];
  steps.forEach(step => stepGraphElements.push({
    uuid: step.uuid,
    outputDataSourceGraphElements: [],
    sourceDataSourceGraphElements: []
  }));

  let dataSourceStepMatches: DataSourceStepMatch[] = getDataSourceStepMatches(dataSources.map(dataSource => dataSource.uuid), steps.map(step => step.uuid));
  dataSourceStepMatches.forEach(match => {
    let dataSourceGraphElement: DataSourceGraphElement = dataSourceGraphElements.filter(dataSourceGraphElement => dataSourceGraphElement.uuid === match.dataSource)[0];
    let stepGraphElement: StepGraphElement = stepGraphElements.filter(stepGraphElement => stepGraphElement.uuid === match.step)[0];

    dataSourceGraphElement.stepGraphElements.push(stepGraphElement.uuid);
    stepGraphElement.sourceDataSourceGraphElements.push(dataSourceGraphElement.uuid);
  });

  let dataSourceStepCombinationsWhichAlreadyCreatedNewDataSources: KeyValuePair[] = [];

  stepGraphElements.forEach(stepGraphElement => {
    let currentStep: Step = stepMap.get(stepGraphElement.uuid);
    if (currentStep.type === AnalysisType.ONE_TO_ONE) {
      stepGraphElement.sourceDataSourceGraphElements.forEach(sourceDataSourceGraphElement => {
        currentStep.outputLabelsArray.forEach(outputLabels => {
          let currentDataSource: DataSource = dataSourceMap.get(sourceDataSourceGraphElement);
          let outputDataSource: DataSource = {
            uuid: uuidv4(),
            name: currentDataSource.name + '->' + currentStep.name,
            labels: [...currentDataSource.labels, ...outputLabels], // TODO prevent doubles, overwrite old ones
            depth: currentDataSource.depth + 1
          };
          dataSourceMap.set(outputDataSource.uuid, outputDataSource);
          // incrementNumberOfNodesInDepth(outputDataSource.depth);
          dataSourceGraphElements.push({
            uuid: outputDataSource.uuid,
            stepGraphElements: [],
            creatorStepGraphElement: currentStep.uuid
          });
          stepGraphElements.filter(stepGraphElement => stepGraphElement.uuid === currentStep.uuid)[0].outputDataSourceGraphElements.push(outputDataSource.uuid);
        });
      });
    } else if (currentStep.type === AnalysisType.ALL_TO_ONE) {
      let sourceDataSources: DataSource[] = [];
      stepGraphElement.sourceDataSourceGraphElements.forEach(sourceDataSource => sourceDataSources.push(dataSourceMap.get(sourceDataSource)));
      let labelIntersection: KeyValuePair[] = [];
      if (sourceDataSources.length > 0) {
        sourceDataSources[0].labels.forEach(label => {
          let labelIsPartOfEverySourceDataSource: boolean = true;
          for (let sourceDataSource of sourceDataSources) {
            if (sourceDataSource.labels[label.key] === undefined || sourceDataSource.labels[label.key] !== label.value) {
              labelIsPartOfEverySourceDataSource = false;
              break;
            }
          }
          if (labelIsPartOfEverySourceDataSource) {
            labelIntersection.push(label);
          }
        });
      }
      currentStep.outputLabelsArray.forEach(outputLabels => {
        let uuid: string = uuidv4();
        let outputDataSource: DataSource = {
          uuid: uuid,
          name: '(all-to-one-' + uuid + ')->' + currentStep.name,
          labels: [...labelIntersection, ...outputLabels], // TODO prevent doubles, overwrite old ones
          depth: Math.max(...sourceDataSources.map(dataSource => dataSource.depth))  + 1
        };

        dataSourceMap.set(outputDataSource.uuid, outputDataSource);
        // incrementNumberOfNodesInDepth(outputDataSource.depth);
        dataSourceGraphElements.push({
          uuid: outputDataSource.uuid,
          stepGraphElements: [],
          creatorStepGraphElement: currentStep.uuid
        });
        stepGraphElements.filter(stepGraphElement => stepGraphElement.uuid === currentStep.uuid)[0].outputDataSourceGraphElements.push(outputDataSource.uuid);

      })
    }
  });


  return {
    dataSourceGraphElements: dataSourceGraphElements,
    stepGraphElements: stepGraphElements
  };
}

// function incrementNumberOfNodesInDepth(depth: number): void {
//   if (amountOfNodesByDepth[depth] === undefined || amountOfNodesByDepth[depth] === null) {
//     amountOfNodesByDepth[depth] = 1;
//   } else {
//     amountOfNodesByDepth[depth] = amountOfNodesByDepth[depth] + 1;
//   }
// }

function fillDataSourceMap(dataSources: DataSource[]): void {
  dataSources.forEach(dataSource => {
    dataSourceMap.set(dataSource.uuid, dataSource);
      // incrementNumberOfNodesInDepth(dataSource.depth);
  });
}

function fillStepMap(steps: Step[]): void {
  steps.forEach(step => stepMap.set(step.uuid, step));
}

function getAllDataSources(): DataSource[] {
  return Array.from(dataSourceMap.keys()).map(dataSourceKey => dataSourceMap.get(dataSourceKey));
}

function getAllSteps(): Step[] {
  return Array.from(stepMap.keys()).map(stepKey => stepMap.get(stepKey));
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


  ngAfterContentInit() {

    let initialDataSources: DataSource[] = this.getDataSourcesRaw().items.map(dataSourceRaw => ({
      name: dataSourceRaw.metadata.name,
      labels: Object.keys(dataSourceRaw.metadata.labels).map(dataSourceLabelKey => ({
        key: dataSourceLabelKey,
        value: dataSourceRaw.metadata.labels[dataSourceLabelKey]
      })),
      uuid: uuidv4(),
      depth: 0
    }));
    fillDataSourceMap(initialDataSources);
    let initialSteps: Step[] = this.getStepsRaw().items.map(stepRaw => ({
      name: stepRaw.metadata.name,
      keyValuePairs: stepRaw.spec.ingest.map(ingest => ({
        regex: ingest.check === 'regex',
        key: ingest.key,
        value: ingest.value
      })),
      uuid: uuidv4(),
      type: stepRaw.spec.type,
      outputLabelsArray: stepRaw.spec.outputs.map(output => {
        return Object.keys(output.labels).map(labelKey => ({
          key: labelKey,
          value: output.labels[labelKey]
        }));
      })
    }));
    fillStepMap(initialSteps);

    let kubernetesGraph: KubernetesGraph = createKubernetesGraph(getAllDataSources(), getAllSteps());
    // kubernetesGraph = createKubernetesGraph(getAllDataSources(), getAllSteps());

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
        x: 10 + ((this.width + 150) * 2) * dataSource.depth,
        y: 10 + 1.50 * this.height * i,
        width: this.width,
        height: this.height
      }));
    let stepsNodes: D3Node[] = kubernetesGraph.stepGraphElements.map(stepGraphElement => stepMap.get(stepGraphElement.uuid))
      .map((step, i) => ({
        id: step.uuid,
        text: step.name,
        x: 160 + this.width,
        y: 10 + 1.50 * this.height * i,
        width: this.width,
        height: this.height
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

    d3.select('#mysvg');

    const graph = {
      nodes: nodes,
      edges: edges,
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

    document.getElementById('mysvg').setAttribute('width', '2000');
    document.getElementById('mysvg').setAttribute('height', '20000');

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
        },
        {
          "apiVersion": "zerops.com/v1",
          "kind": "ZerOpsStep",
          "metadata": {
            "annotations": {
              "kubectl.kubernetes.io/last-applied-configuration": "{}"
            },
            "creationTimestamp": "2019-11-07T13:35:48Z",
            "generation": 1,
            "name": "example-step-matches-ab",
            "resourceVersion": "55010",
            "selfLink": "/apis/zerops.com/v1/zerops-steps/example-step-one-to-one-1",
            "uid": "a3eaaa0a-46f2-46bf-8604-220542d4a78d"
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
            "type": "one-to-one"
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
              "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"zerops.com/v1\",\"kind\":\"ZerOpsDataSource\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"A\"},\"name\":\"data-source-a\"},\"spec\":{\"url\":\"tcp://172.17.0.8:9000\"}}\n"
            },
            "creationTimestamp": "2019-11-07T13:34:46Z",
            "generation": 1,
            "labels": {
              "name": "A"
            },
            "name": "data-source-a-(2)",
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
}
