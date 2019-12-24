#!/bin/bash

# Set work dir
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
PROJ_DIR="$SCRIPT_DIR/.."

echo "Running test $PROJ_DIR."

pushd $PROJ_DIR

go test ./... -coverprofile cover.out
echo "------Coverage per function------"
go tool cover -func cover.out
echo "----Cover html representation----"
go tool cover -html cover.out

popd
