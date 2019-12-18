#!/bin/bash

source scripts/ci/build.sh
source scripts/ci/push.sh

curl -H "TOKEN: $DEPLOY_TOKEN" -X POST https://$DEPLOY_HOST/svoyak/deploy
