#!/bin/bash

curl -H "TOKEN: $DEPLOY_TOKEN" -X POST https://$DEPLOY_HOST/svoyak/deploy
