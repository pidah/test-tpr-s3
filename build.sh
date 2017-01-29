#!/bin/bash

set -e

VERSION=${1:-0.0.1}
DOCKER_REPO="pidah/test-master"
DOCKER_IMAGE=${DOCKER_REPO}:${VERSION}

# run tests before build
#echo "Running go test..."
#go test

echo "Building application..."
git tag ${VERSION}
git push origin ${VERSION}
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.7 bash -c "go get -d && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o test-master -v"

echo "Building docker image..."
docker build -t ${DOCKER_IMAGE} .
docker push ${DOCKER_IMAGE}
