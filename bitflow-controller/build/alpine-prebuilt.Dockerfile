# teambitflow/bitflow-controller
# Copies pre-built binaries into the container. The binaries are built on the local machine beforehand:
# ./alpine-build.sh
# docker build -t teambitflow/bitflow-controller -f alpine-prebuilt.Dockerfile .
FROM alpine:3.9
RUN apk --no-cache add libstdc++
COPY _output/bin/alpine/bitflow-controller /
ENTRYPOINT ["/bitflow-controller"]
