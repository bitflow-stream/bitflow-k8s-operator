# Introduction

This repository contains the Bitflow controller, which runs within a Kubernetes cluster.
The functionality of the controller is to automatically create and delete data analysis containers based on custom resource definitions of data sources (*BitflowSources*) and data analysis steps (*BitflowSteps*).
The custom resource objects can be managed by users through the Kubernetes API and tools.
If the Bitflow controller is running in the cluster, it will watch for updates to those custom resources and continuously ensure the desired state.

A BitflowSource object contains a URL for accessing the data and a number of string labels that describe arbitrary properties of the data source.
The labels should contain things like the type of monitoring technology used, the monitored component, system layer, and so on.
An example of a data source definition can be found [here](examples/example-data-source.yaml).
TODO Describe how the bitflow-collector creates data sources automatically.
More data sources can be added manually for testing purposes.
When adapting external monitoring tools, something like a proxy process should register a data source and provide the appropriate URL and labels.

Multiple examples for BitflowStep objects can be found [here](examples/).
A BitflowStep object contains a list of ingest selectors that will be matched against the labels of each data source object.
When all ingest selectors evaluate as `true` for the labels of a data source, that data source is considered appropriate for the analysis step.
An optional **`check`** property can be defined for different types of match checks:
- `check: wildcard`: this is the default, when the `check property` is omitted. A wildcard match is performed, i.e. the star character **\*** can be used anywhere in the key or value of the selector to match arbitrary character sequences.
- `check: exact`: perform an exact string match using the `key` and `value` properties of the selector
- `check: regex`: both the `key` and `value` properties are evaluated as Go regexes against the data source labels
- `check: present`: checks whether a label with exactly the given `key` exists. The property `value` must be omitted in the selector
- `check: absent`: checks whether a label with exactly the given `key` does not exists. The property `value` must be omitted in the selector

Currently the selector matching logic is quite limited, no complex boolean expressions are possible (beyond the `AND`-combination of the selector list).

There are two supported types of analysis steps: **one-to-one** and **all-to-one**.
The desired type can be set with the **type: one-to-one** or **type: all-to-one** property in the BitflowStep object. When the property is omitted, the default is one-to-one.
A one-to-one step results in a separate data analysis Pod for every data source that the step matches.
That Pod receives all information necessary for accessing the data source that it should analyse.
An all-to-one step results in a single data analysis Pod for the step, whose task it is to analyse *all* data sources that are matched by the selectors of the step.

The data analysis step object defines the analysis workload in the form of a Kubernetes [Pod template](https://kubernetes.io/docs/concepts/workloads/pods/pod-overview/).
It contains all information necessary to start a data analysis container.
When such a Pod is instantiated by the Bitflow controller, it will perform the following modifications to the template, before submitting it to Kubernetes:

- Set the `.metadata.namespace` field to the controllers namespace
- Set the `.metadata.name` field to a generated name that includes a configurable prefix (optional`-pod-prefix` parameter), the name of the data analysis step, and a UUID suffix
- Add an entry to `.metadata.ownerReferences` that points to the controller Pod. This will cause Kubernetes to automatically delete all Bitflow analysis pods when the controller is removed from the cluster
- Add the following entries to `.metadata.labels`:
    - A number of configurable labels given as parameters to the controller (required `-select` parameter(s), example: `-select controller=bitflow -select env=production`)
    - `bitflow-analysis-step`: The name of the underlying analysis step
    - `bitflow-analysis-type`: The type of the underlying analysis step (`one-to-one` or `all-to-one`)
    - `bitflow-data-source-name`: The name of the used data source (only for `one-to-one` steps)
    - Note: the labels of the data source of a `one-to-one` are not added as a label, because there is a limitation to characters that can be used for labels
- Add the following environment variables to all containers of the Pod (`.spec.containers[*].env`):
    - A number of configurable variables that can be given as parameters to the controller (optional `-extra-env` parameter(s), example: `-extra-env DB_URL=... -extra-env ENV=production`)
    - `BITFLOW_ANALYSIS_STEP`: The name of the analysis step object used as template
    - `BITFLOW_ANALYSIS_TYPE`: The type of the underlying analysis step (`one-to-one` or `all-to-one`)
    - `BITFLOW_DATA_SOURCE`: This has different meanings for different analysis step types:
        - For `one-to-one` steps it contains the data source URL
        - For `all-to-one` steps it contains a HTTP GET URL that returns a JSON encoded list of data source URLs. This list can change, and the analysis step must take care to regularly poll for updates to the data source list
    - `BITFLOW_DATA_SOURCE_NAME`: The name of the used data source (only for `one-to-one` steps)
    - `BITFLOW_DATA_SOURCE_LABELS`: The labels of the data source (only for `one-to-one` steps), in the format `a=b, x=y`
- In addition to setting these environment variables, they are also directly replaced by their values in certain parts of the pod spec. This enables using these variables without having to access the environment variables in the containers. Before the replacements, all variable names are wrapped in curly braces (e.g. `{BITFLOW_DATA_SOURCE}`). This is intended as the main and most convenient way to inject the Bitflow configurations into the analysis pods. The following parts of the pod spec are patched this way:
    - All `command` and `args` elements in all containers (`.spec.containers[*].command[*]` and `.spec.containers[*].args[*]`)
    - All values of all environment variables of all containers (`.spec.containers[*].env[*].value`). This is rarely useful, but allows embedding the special Bitflow values in existing environment variables.

The controller mainly operates on the described Kubernetes objects, but it also has a small REST API with the following paths:
- `GET /ip`: Return the clients IP (used by the Bitflow Collector to reliably determine their external IP, for using it inside the data source URL.
- `POST /triggerUpdate`: Force an update check to see if any analysis Pod should be created or deleted. This is automatically done whenever there are any changes to a custom resource object or managed Pod.
- `GET /status`: Return a JSON-encoded status of all managed resources it contains:
    - Info about the last update check and the last request to Kubernetes
    - A list of all data sources, and what Pods and analysis steps they are handled by
    - A list of all analysis steps and the Pods and data sources associated with them
    - A list of all created analysis Pods, and what data sources and analysis steps they belong to
- `GET /dataSources/:podName`: Return a JSON-encoded list of data source URLs to be used by the given `all-to-one` analysis pod. A complete URL to this REST API path is passed to such Pods in the `BITFLOW_DATA_SOURCE` environment variable (see above).

See below on how to access the REST API and logs of a Bitflow controller running in a Kubernetes cluster.

### Chaining Analysis Steps

BitflowSource can be created manually or programmatically, and will result in data analysis pods according to existing BitflowStep objects.
In order to automatically create entire data analysis pipelines (chains of analysis steps), the Bitflow controller automatically creates BitflowSource objects for data analysis pods, when the underlying BitflowStep object contains the **outputs:** property.
The `outputs:` property is optional and can contain a list of output definitions, that include a name, a URL template, and a list of at least one label.
See the data analysis step definitions in the `examples/` folder for examples of output definitions.
When an analysis pod is instantiated (regardless of the analysis step type), the Bitflow controller creates a BitflowSource object for each element in the `outputs:` list, with the following properties:

- The data source URL (`.spec.url`) is based on the URL template given in the `output` object, except that the host part of the URL is replaced by the cluster-internal IP of the analysis pod. The rest of the URL remains unmodified, and previous values for the host are overwritten. Example: An output URL of `http://:7000/data` could result in a BitflowSource URL value of `http://10.233.122.245:7000/data`, where `10.233.122.245` is the cluster-internal IP address of the analysis pod.
- The name (`.metadata.name`) is set to `output.POD_NAME.OUTPUT_NAME`, where `POD_NAME` is the name of the analysis pod, and `OUTPUT_NAME` is the name of the output element in the `outputs:` list (example: `output.pod123.analysis-results`)
- The namespace (`.metadata.namespace`) is set to the namespace of the analysis pod
- An `ownerReference` is added, which points to the analysis pod. This ensures that the BitflowSource object is automatically deleted when the analysis pod vanishes.
    - Note: because the Kubernetes garbage collection is asynchronous and slow by default, the Bitflow controller also actively deletes BitflowSource objects that reference missing analysis pods
- The following labels are set:
    - All labels that are shared by all data sources that are analysed by the data analysis pod
        - For `one-to-one` analysis steps, this means the labels of the input data source are copied
        - For `all-to-one` analysis steps, only the labels that are *exactly* shared by *all* input data sources, are propagated
    - `bitflow-pipeline-depth` is set to the length of the analysis pipeline that resulted in the current BitflowSource
    - `bitflow-pipeline-N` (where `N` is the current value of `bitflow-pipeline-depth`) is set to the name of the BitflowStep of the underlying analysis pod
        - Example: the following labels on a BitflowSource indicate a two-element analysis pipeline that executes anomaly-detection first, and a root-cause-analysis after that. Both elements correspond to created analysis pods.
            - `bitflow-pipeline-depth: 2`
            - `bitflow-pipeline-1: anomaly-detection`
            - `bitflow-pipeline-2: root-cause-analysis`
    
**Important note:** BitflowSteps cannot recursively analyse their own outputs. This means that each BitflowStep can occur *at most once* in an analysis pipeline.

# Kubernetes Operator SDK

The Bitflow Operator was initialized using the operator-sdk.
I would recommend to use this package for further development, because it simplifies the code generation and build process.
Unfortunately the operator-sdk does not support Windows.

The entrypoint for the operator logic is placed in `bitflow-controller/pkg/controller/bitflow/controller.go`.
The `Reconcile()`-method is called whenever a BitflowStep or a related resource (pod, BitflowSource) was updated.
Keep it mind, that this method is always called with the name of the BitflowStep-resource, even if a pod
was updated, the call will be mapped to the corresponding BitflowStep.
So the focus of reconciling is to check, maintain and restore the desired state for each BitflowStep.

## Extending Bitflow resources
The system currently maintains the two resource types BitflowStep and BitflowSource.
If additional custom resources become necessary, you can initialize them with the following command
- `operator-sdk add api --api-version=bitflow.com --kind=<Resource name>`
Initialize a corresponding controller, if you need to, by running
- `operator-sdk add controller --api-version=bitflow.com --kind=<resource name>`
Important: The operator-sdk provides code generation for openapi, deepcopy and for the crd yaml files.
So whenever you apply changes to the Go type definition files, you should run the following two commands.
- `operator-sdk generate k8s`
- `operator-sdk generate openapi`
Be aware that earlier changes to the generated files will be overridden by these commands.


## About the cache
The cache that is provided by the operator framework as if it was a normal client.
It caches every sort of resource that is requested by the application, it is not limited to those resources that are specifically watched. However the cache does not store all available resources, but only those which are requested. As soon as a resource type is requested that is not yet stored in the cache, a corresponding informer is initialized - see [here](https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/cache/internal/informers_map.go#L146). The cache is then filled with all available objects of that type and your request will return.

The informer that is being used here is the [cache.SharedIndexInformer](https://godoc.org/k8s.io/client-go/tools/cache#NewSharedIndexInformer)

## Leader election
The operator framework prevents two separate pods in the same namespace from running the same controller. The mechanism to prevent this is called leader election, which basically creates a ConfigMap that serves as a lock.
The pod that creates the map will become the leader, the second controller instance will wait until the lock is resolved. See [here](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md) for general information or [here](https://github.com/operator-framework/operator-sdk/blob/master/pkg/leader/leader.go#L38) for specifics 

## Configuration
The configuration is mainly handled by a kubernetes ConfigMap. However some values that are processed on application startup and cannot be changed during execution are passed as environment variables. The following sections will discuss each parameter in detail.

The Log level can be specified as an argument to the operators start command (see ./deploy/operator.yaml). The default Log level is Info. Other allowed values are -v, -q and -qq (Degub, Warn and Error).

### Environment parameters
- **WATCH_NAMESPACE**: The name of the namespace in which the operator is deployed
- **POD_NAME**: The name of the operator pod
- **OPERATOR_NAME**: Registered name of the operator in the Kubernetes Operator framework
- **CONFIG_MAP**: Name of the ConfigMap resource that stores additional configuration parameters
- **POD_IP**: The IP of the pod that runs the operator. This is used when generating the operators API endpoint URL.
- **API_LISTEN_PORT**: The port on which the operator will provide its endpoint
- **BITFLOW_ID_LABELS**: Comma-separated list of key-value labels that will identify the pods managed by the controller 
- **CONCURRENT_RECONCILE**: Max number of reconciles processed concurrently. Default: 1
- **STATISTICS**: If present (any non empty value) the application will gather information about the number of performed reconciles or the avg time of all reconcile. If not needed this feature should disabled, i.e. in production environments

### CongigMap parameters
- **external.source.label**: BitflowSources that are located inside the kubernetes cluster but are not related to any BitflowStep output can specify their host node. This is done by specifiing a label this key and the nodes name as a value. Default: "bitflow-nodename"
- **resource.limit.slots**: Specifies how many pods should initially fit into the given resource limit. Default: "2"
- **resource.limit.slots.grow**: When reaching the resource limit, this parameter specifies by how much the resources will be reduced. Should be more than 1. Default: "2.0"
- **resource.limit**: Resource limit that should be applied to all Bitflow pods.  Specified as a relative value. A value lower than 0 or greater than 1 will disable the resource restrictions. Default: "0.1"
- **resource.limit.node.annotation**: You can specify the resource limit for each node individually by annotating a node with this annotation and the desired limit. Default: "bitflow-resource-limit"
- **extra.env**: Same as the above, but with environment variables.
- **delete.grace.period**: Grace for pods that should be restarted, i.e. for resource adjustments. TO reduce the down time one may reduce the grace period using this parameter. Pod deletions that are triggerd by other reasons (i.e. step or source deletion) will not be effected by this value. Specified in seconds. Default: "30"
- **pod.spawn.period**: All current Bitflow pods are regularly checked for the correct resource limits,so that the overall node resource limit is realized. The interval in which this validation process is performed can be specified by this parameter. However the validation will be performed during one steps Reconcile not independently. A negative value disables the interval, which means that the validation will be performed whenever `Reconcile` is called. Default: "-1"
- **reconcile.heartbeat**: Complements the above validation period. If there was a `Reconcile` performed after the last validtion process and this heartbeat is triggered it will perform a validation should the above period be already expired. This ensures that there will be no `Reconcile` that is left unchecked for an unknown amount of time. Specified in seconds. Default: "120"
- **scheduler**: BitflowSteps come with the `Scheduler` - Attribute. It is not specified or the specified algorithm does not provide a node, then multiple default scheduler may be supplied as a comma separated list. These algorithms are executed from left to right until a node is found. The first node that is found will be used for scheduling. Default: "sourceAffinity,first"
