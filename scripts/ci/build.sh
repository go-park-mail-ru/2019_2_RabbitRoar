#!/bin/bash

# Stop execution on fails
set -e

# Set work dir
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
PROJ_DIR="$SCRIPT_DIR/../.."

echo "Running build $PROJ_DIR."

pushd $PROJ_DIR

docker build -f deployments/Dockerfile-session -t alexnav/svoyak-session .
docker build -f deployments/Dockerfile-application -t alexnav/svoyak-application .
docker build -f deployments/Dockerfile-chat -t alexnav/svoyak-chat .
docker build -f deployments/Dockerfile-game -t alexnav/svoyak-game .

docker login --username=alexnav --password=$DOCKER_HUB_TOKEN

docker push alexnav/svoyak-session
docker push alexnav/svoyak-application
docker push alexnav/svoyak-chat
docker push alexnav/svoyak-game

popd
