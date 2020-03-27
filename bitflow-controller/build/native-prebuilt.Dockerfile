# teambitflow/bitflow-controller
# Build the controller before the container:
# ./native-build.sh
# docker build -t teambitflow/bitflow-controller -f native-prebuilt.Dockerfile .
FROM alpine:3.9
RUN apk --no-cache add libstdc++
# FROM registry.access.redhat.com/ubi7/ubi-minimal:latest
COPY _output/bin/bitflow-controller /
ENTRYPOINT ["/bitflow-controller"]
