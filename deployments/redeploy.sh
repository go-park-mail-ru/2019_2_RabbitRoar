#!/bin/bash

echo "Deploying stack."
docker-compose -p svoyak -f docker-compose-prod.yml up -d
echo "Deployed successfully!"
