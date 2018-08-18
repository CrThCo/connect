#!/bin/bash
echo "preparing dev environment, these env variable only available in this script scope."
set -o allexport
source "./.env"
set +o allexport

mkdir -p $FILE_STORAGE_PATH/profile
mkdir -p $FILE_STORAGE_PATH/content
go run *.go