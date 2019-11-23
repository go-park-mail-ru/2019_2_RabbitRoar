#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="$SCRIPT_DIR/../internal/pkg/models"

echo "Running easyjson generation in $WORK_DIR."

pushd $WORK_DIR

echo "Removing existing generation."
rm -f -v *_easyjson.go
echo "Done!"

echo "Generating easyjson for models."
go generate .
echo "Done!"

popd
