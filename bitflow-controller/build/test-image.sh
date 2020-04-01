#!/usr/bin/env sh

test $# = 1 || { echo "Need 1 parameter: image tag to test"; exit 1; }
IMAGE="bitflowstream/bitflow-controller"
TAG="$1"

# Sanity check: image starts, outputs valid JSON, and terminates.
2>&1 docker run "$IMAGE:$TAG" -help | tee /dev/stderr | grep "Usage of /bitflow-controller" > /dev/null
