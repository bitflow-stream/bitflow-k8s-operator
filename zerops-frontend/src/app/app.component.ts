import {Component} from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'zerops-frontend';

  steps = this.getSteps();
  dataSources = this.getDataSources();

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
