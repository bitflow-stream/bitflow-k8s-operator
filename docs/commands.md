
# Working with Custom Resource Objects

The Kubernetes concept of CRDs (custom resource definitions) allows to define objects that are stored in Kubernetes and can be manipulated in the same way as regular Kubernetes objects like Pods or Services. The file [custom-resource-definition.yaml](helm/templates/custom-resource-definition.yaml) creates two custom resource types: BitflowSource (abbreviation: `bso`) and BitflowStep (abbreviation: `bst`). Examples for the specification of both objects can be found in [example-data-source.yaml](examples/example-data-source.yaml) and [example-analysis-step.yaml](examples/example-analysis-step.yaml).

These objects can be manipulated through the `kubectl` tool, using their abbreviations:
- `kubectl get bst`: List all analysis steps
- `kubectl create -f step.yaml`: Create an analysis step or data source
- `kubectl delete bso <step-name>`: Delete a data source
- `kubectl describe bst <step-name>`: Show details of an analysis step

Read the documentation of `kubectl` for more commands and details.

#### Getting the Bitflow Controller logs
- `kubectl logs $(kubectl get pod -l app.kubernetes.io/name=bitflow-controller -o jsonpath='{.items[0].metadata.name}')`

#### Accessing the Bitflow Controller REST API
- Enable port forwarding:
    - `kubectl port-forward $(kubectl get pod -l app.kubernetes.io/name=bitflow-controller -o jsonpath='{.items[0].metadata.name}') 8000:rest-api`
- Access: `curl http://localhost:8000/status | jq`
- Other paths in the API: `GET /ip`, `POST /triggerUpdate`
