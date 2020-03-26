# teambitflow/bitflow-controller
# Build the controller before the container:
# ./native-build.sh
# docker build -t teambitflow/bitflow-controller -f native-prebuilt.Dockerfile .

FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

#FROM alpine:3.9
#RUN apk --no-cache add libstdc++

WORKDIR /bitflow
COPY _output/bin/bitflow-controller /
ENTRYPOINT ["/bitflow-controller"]
