# teambitflow/bitflow-api-proxy
# Build the bitflow-api-proxy before the container (from parent directory):
# go build -v -o ./build/_output/bin/bitflow-api-proxy .
# docker build -t teambitflow/bitflow-api-proxy -f build/cached.Dockerfile .

# FROM registry.access.redhat.com/ubi7/ubi-minimal:latest
FROM alpine:3.9
RUN apk --no-cache add libstdc++
WORKDIR /bitflow
COPY build/_output/bin/bitflow-api-proxy .
ENTRYPOINT ["/bitflow/bitflow-api-proxy"]

