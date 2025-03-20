#! /bin/sh

export GOPROXY=https://goproxy.cn

go run service/cmd/workflow-app/engine/main.go -config.path /app/config/