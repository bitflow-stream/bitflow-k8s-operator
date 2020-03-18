# teambitflow/bitflow-controller
# Build the controller before the container:
# go build -v -o ./build/_output/bin/bitflow-controller ./cmd/manager
FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

WORKDIR /bitflow
COPY build/_output/bin/bitflow-controller .
ENTRYPOINT ["/bitflow/bitflow-controller"]
