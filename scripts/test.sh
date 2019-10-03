#!/bin/bash

SCRIPT_DIR=$(cd $(dirname $) && pwd)
WORK_DIR="$SCRIPT_DIR/../"

echo "Running goimport in $WORK_DIR"

pushd $WORK_DIR
go test ./... -coverprofile cover.out
echo "------Coverage per function------"
go tool cover -func cover.out
echo "----Cover html representation----"
go tool cover -html cover.out
popd
