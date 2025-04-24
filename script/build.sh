#!/bin/bash

export GOPROXY=https://proxy.golang.com.cn,direct
BASE_DIR=$(dirname "${BASH_SOURCE[0]}")"/.."
BASE_DIR_PATH=$(cd ${BASE_DIR} && pwd)

BUILD_DIR="build"
DATE=$(date +'%Y_%m_%d_%H%M')
BUILD_FILE="$BUILD_DIR/feature-flag_$DATE.tar.gz"
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}/bin

go build -o ${BUILD_DIR}/bin/feature-flag cmd/main.go

mkdir ${BUILD_DIR}/conf
cp -r conf/* ${BUILD_DIR}/conf/
mkdir ${BUILD_DIR}/script
cp -r script/start.sh ${BUILD_DIR}/script/

COMMIT_ID=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "VERSION = \"$COMMIT_ID\"" >"$BUILD_DIR/version"
echo "BRANCH = \"$BRANCH\"" >>"$BUILD_DIR/version"
echo "DATE = \"$DATE\"" >>"$BUILD_DIR/version"

tar -zcvf $BUILD_FILE $BUILD_DIR >/dev/null
echo "${BASE_DIR_PATH}/${BUILD_FILE}"
