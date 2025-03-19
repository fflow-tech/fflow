#!/bin/bash

# 检查参数数量
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <workflow-file> <input-file> [--build]"
    echo "  --build: Optional flag to compile the Go program before execution"
    exit 1
fi

WORKFLOW_FILE=$1
INPUT_FILE=$2
BUILD_FLAG=false

# 检查是否有第三个参数，并且是否为 --build
if [ "$#" -eq 3 ] && [ "$3" = "--build" ]; then
    BUILD_FLAG=true
fi

# 如果指定了编译标志，则编译Go程序
if [ "$BUILD_FLAG" = true ]; then
    echo "Compiling Go program..."
    go build -o workflow-cli service/cmd/workflow-cli/main.go
    if [ $? -ne 0 ]; then
        echo "Failed to compile Go program"
        exit 1
    fi
fi

# 执行工作流
echo "Executing workflow..."
./workflow-cli -file "$WORKFLOW_FILE" -input "$INPUT_FILE" -config.name "app" -config.type "yaml" -config.path "./"
