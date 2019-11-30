#!/bin/bash

echo "Recreating container $1."
docker-compose -p svoyak -f docker-compose-prod.yml up --build -d --no-deps $1
