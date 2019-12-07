#!/bin/bash

echo "Deploying stack."
docker-compose -p svoyak -f docker-compose-prod.yml up --build -d
echo "Deployed successfully!"
