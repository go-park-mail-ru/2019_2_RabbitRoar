#!/bin/bash

SCRIPT_DIR=$(cd $(dirname $) && pwd)
WORK_DIR="$SCRIPT_DIR/../internal/pkg/models"

echo "Running easyjson generation in $WORK_DIR."

pushd $WORK_DIR

for file in ./*
do
  echo "Generation for $file."
  easyjson -all "$file"
done

popd
