#!/bin/bash

# Stop execution on fails
set -e

# Set work dir
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
PROJ_DIR="$SCRIPT_DIR/../.."

echo "Running build $PROJ_DIR."

pushd $PROJ_DIR

docker build -f deployments/Dockerfile-session -t alexnav/svoyak-session .
docker tag alexnav/svoyak-session alexnav/svoyak-session:$TRAVIS_BUILD_NUMBER

docker build -f deployments/Dockerfile-application -t alexnav/svoyak-application .
docker tag alexnav/svoyak-application alexnav/svoyak-application:$TRAVIS_BUILD_NUMBER

docker build -f deployments/Dockerfile-chat -t alexnav/svoyak-chat .
docker tag alexnav/svoyak-chat alexnav/svoyak-chat:$TRAVIS_BUILD_NUMBER

docker build -f deployments/Dockerfile-game -t alexnav/svoyak-game .
docker tag alexnav/svoyak-game alexnav/svoyak-game:$TRAVIS_BUILD_NUMBER

popd
