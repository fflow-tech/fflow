#! /bin/sh

export GOPROXY=https://goproxy.cn
go install github.com/go-delve/delve/cmd/dlv@latest
dlv --headless --log --listen :9009 --api-version 2 --accept-multiclient debug service/cmd/workflow-app/engine/main.go -- -config.path /app/config/