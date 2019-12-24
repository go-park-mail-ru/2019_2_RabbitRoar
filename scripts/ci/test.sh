#!/bin/bash

# Stop execution on fails
set -e

# Set work dir
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
PROJ_DIR="$SCRIPT_DIR/../.."

echo "Running test $PROJ_DIR."

pushd $PROJ_DIR

docker run --rm -it $(docker build -f deployments/Dockerfile-test -q .)

popd
