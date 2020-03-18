# Bitflow REST API Proxy

TODO Update Readme

The `bitflow-api-proxy` is used to get the current state of the Kubernetes cluster.
It provides endpoints for nodes, pods and for Bitflow resources such as BitflowDataSources and BitflowSteps.

We also provide some metric endpoints, but these require a running metric-server on the kubernetes cluster.

## Deployment
See `deploy-proxy.yaml` for a deployment example.
