# bitflowstream/bitflow-api-proxy
# docker build -t bitflowstream/bitflow-api-proxy -f bitflow-api-proxy/build/multi-stage/alpine-full.Dockerfile .
FROM golang:1.14.1-alpine as build
ENV GO111MODULE=on
RUN apk --no-cache add git gcc g++ musl-dev openssh-client mercurial
WORKDIR /build

# Copy go.mod first and download dependencies, to enable the Docker build cache
COPY bitflow-api-proxy/go.mod bitflow-api-proxy/go.mod
COPY bitflow-controller/go.mod bitflow-controller/go.mod
RUN sed -i $(find -name go.mod) -e '\_//.*gitignore$_d' -e '\_#.*gitignore$_d' && \
    cd bitflow-api-proxy && \
    go mod download

# Copy rest of the source code and build
# Delete go.sum files and clean go.mod files form local 'replace' directives
COPY . .
RUN find -name go.sum -delete && \
    sed -i $(find -name go.mod) -e '\_//.*gitignore$_d' -e '\_#.*gitignore$_d' && \
    cd bitflow-api-proxy && \
    go build -o /bitflow-api-proxy .

FROM alpine:3.11.5
RUN apk --no-cache add libstdc++
COPY --from=build /bitflow-api-proxy /
ENTRYPOINT ["/bitflow-api-proxy"]
