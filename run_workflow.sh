#!/bin/bash

# 检查参数数量
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <workflow-file> <input-file>"
    exit 1
fi

WORKFLOW_FILE=$1
INPUT_FILE=$2

# 编译Go程序
echo "Compiling Go program..."
go build -o workflow-cli service/cmd/workflow-cli/main.go
if [ $? -ne 0 ]; then
    echo "Failed to compile Go program"
    exit 1
fi

# 执行工作流
echo "Executing workflow..."
./workflow-cli -file "$WORKFLOW_FILE" -input "$INPUT_FILE" -config.name "app" -config.type "yaml" -config.path "./"
