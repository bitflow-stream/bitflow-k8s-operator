const stepsRaw = {
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


const dataSourcesRaw = {
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

export {dataSourcesRaw, stepsRaw};
