#!/bin/bash

set -e

IMAGE_REPO="feature-flag"
DOCKER_FILE="docker/Dockerfile"

[ -f $DOCKER_FILE ] || {
  echo "Docker file[$DOCKER_FILE] not exists."
  exit 1
}

bash script/build.sh
DATE=$(date +'%Y_%m_%d_%H%M')
docker build -f $DOCKER_FILE -t $IMAGE_REPO:$DATE.
