const stepsRaw = {
  "apiVersion": "v1",
  "items": [
    {
      "apiVersion": "zerops.com/v1",
      "kind": "ZerOpsStep",
      "metadata": {
        "name": "o2o 1:1",
      },
      "spec": {
        "ingest": [
          {
            "key": "1",
            "value": "1"
          }
        ],
        "outputs": [
          {
            "labels": {
              "2": "2"
            },
            "name": "phys",
            "url": "tcp://:9000"
          }
        ],
        "template": {},
        "type": "one-to-one"
      }
    },
    {
      "apiVersion": "zerops.com/v1",
      "kind": "ZerOpsStep",
      "metadata": {
        "name": "o2o 2:2 3:3",
      },
      "spec": {
        "ingest": [
          {
            "check": "regex",
            "key": "^(2|3)$",
            "value": "^(2|3)$"
          }
        ],
        "outputs": [
          {
            "labels": {
              "2+3": "2+3"
            },
            "name": "phys",
            "url": "tcp://:9000"
          }
        ],
        "template": {},
        "type": "one-to-one"
      }
    },
    {
      "apiVersion": "zerops.com/v1",
      "kind": "ZerOpsStep",
      "metadata": {
        "name": "a2o 2+3:2+3 4:4",
      },
      "spec": {
        "ingest": [
          {
            "key": "2+3",
            "value": "2+3"
          },
          {
            "key": "4",
            "value": "4"
          }
        ],
        "outputs": [
          {
            "labels": {
              "5": "5"
            },
            "name": "phys",
            "url": "tcp://:9000"
          }
        ],
        "template": {},
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

const dataSourcesRaw = {
  "apiVersion": "v1",
  "items": [
    {
      "apiVersion": "zerops.com/v1",
      "kind": "ZerOpsDataSource",
      "metadata": {
        "labels": {
          "1": "1"
        },
        "name": "1:1",
      }
    },
    {
      "apiVersion": "zerops.com/v1",
      "kind": "ZerOpsDataSource",
      "metadata": {
        "labels": {
          "3": "3"
        },
        "name": "3:3",
      }
    },
    {
      "apiVersion": "zerops.com/v1",
      "kind": "ZerOpsDataSource",
      "metadata": {
        "labels": {
          "4": "4"
        },
        "name": "4:4",
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
