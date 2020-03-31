## Building
When building the project you might experience an error, stating that the `hg`-command is not found.
This refers to a missing installation of mercurial - a source control tool.
### Using operator-sdk
Just run
- `operator-sdk build <repo>:<image tag>`

### Without operator-sdk 
- `go build -v -o ./bitflow-controller/build/_output/bin/bitflow-controller ./bitflow-controller/cmd/manager`
- `docker build --tag=bitflowstream/bitflow-controller -f ./bitflow-controller/build/Dockerfile .`

### Without using operator-sdk
You can certainly add resources and controller files by yourself, but pay attention to how the existing resources
a integrated in the system. For example, don`t forget to properly register your new resources with the kubernetes
scheme, otherwise kubernetes will have problems to get or create your resources.

I do not know how to trigger the code generation without the above operator-sdk commands.
