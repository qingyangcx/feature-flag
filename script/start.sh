#!/bin/bash

set -e
set -x

PROJECT_NAME="feature-flag"

function info() {
    (printf >&2 "[\e[34m\e[1mINFO\e[0m] $*\n")
}

function error() {
    (printf >&2 "[\033[0;31mERROR\e[0m] $*\n")
}

function warning() {
    (printf >&2 "[\033[0;33mWARNING\e[0m] $*\n")
}

function ok() {
    (printf >&2 "[\e[32m\e[1m OK \e[0m] $*\n")
}

MODE=""
VERSION=""

function usage() {
    echo "---------------------------------------"
    echo " Usage: ${SCRIPT_NAME} [options...]"
    echo "  "
    echo " -m <mode>          running mode, prod/uat/prod_unicom."
    echo " -v <version>       software version used, example:2024_04_22_1620, default use lastest."
    echo " -h                 help."
    echo "---------------------------------------"
    exit -1
}

while getopts ':m:v:h' OPT; do
    case $OPT in
    m) MODE="$OPTARG" ;;
    v) VERSION="$OPTARG" ;;
    h) usage ;;
    esac
done
shift $(($OPTIND - 1))

cd "$(dirname "${BASH_SOURCE[0]}")"

BUILD_TAR=$(ls "${PROJECT_NAME}*")
info "BUILD_TAR=${BUILD_TAR}"

if [ -z "$BUILD_TAR" ]; then
    error "fetch BUILD_TAR failed"
    exit -1
fi

tar -xf $BUILD_TAR
cd build
BASE_DIR=$(pwd)

CONFIG_FILE="conf/${MODE}.json"
if [ ! -e "${CONFIG_FILE}" ]; then
    error "config file ${CONFIG_FILE} not exist"
    exit -1
fi

bin/${PROJECT_NAME} --config_file=${CONFIG_FILE}
