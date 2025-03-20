#!/bin/bash

echo "Building fflow-cli ..."
go build -o fflow-cli service/cmd/workflow-cli/main.go
if [ $? -ne 0 ]; then
    echo "Failed to build fflow-cli"
    exit 1
fi

echo "Build fflow-cli success"
exit 0