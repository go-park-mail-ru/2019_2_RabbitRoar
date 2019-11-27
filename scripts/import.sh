#!/bin/bash

SCRIPT_DIR=$(cd $(dirname $) && pwd)
WORK_DIR="$SCRIPT_DIR/../"

echo "Running goimport in $WORK_DIR"

pushd $WORK_DIR
find -name "*.go" | xargs goimports -w -v
popd
