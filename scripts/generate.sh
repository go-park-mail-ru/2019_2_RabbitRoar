#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="$SCRIPT_DIR/../"

echo "Running easyjson generation in $WORK_DIR."

pushd $WORK_DIR


echo "Removing existing easyjson generation."
find . -name "*_easyjson.go" | xargs rm -vf
echo "Done!"

echo "Generating easyjson for models."
go generate ./internal/pkg/models
echo "Done!"

echo "Removing existing protobuf generation."
find . -name "*.pb.go" | xargs rm -vf
echo "Done!"

echo "Generating grpc interface."
go generate ./internal/pkg/session/delivery/grpc
echo "Done!"

popd
