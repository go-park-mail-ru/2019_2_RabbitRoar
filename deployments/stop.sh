#!/bin/bash

echo "Stopping stack."
docker-compose -p svoyak -f docker-compose-prod.yml down
echo "Stopped successfully!"
