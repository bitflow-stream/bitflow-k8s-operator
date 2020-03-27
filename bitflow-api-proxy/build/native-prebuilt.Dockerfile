# teambitflow/bitflow-api-proxy
# Build the bitflow-api-proxy before the container:
# ./native-build.sh
# docker build -t teambitflow/bitflow-api-proxy -f native-prebuilt.Dockerfile .
FROM alpine:3.9
RUN apk --no-cache add libstdc++
COPY _output/bin/bitflow-api-proxy /
ENTRYPOINT ["/bitflow-api-proxy"]
