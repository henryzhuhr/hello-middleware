#!/bin/bash

# Init Python environment
uv sync

# Install Go tools
export GOPROXY=https://goproxy.cn,direct

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
# for Go1.23
go install golang.org/x/tools/gopls@v0.18.1
go install github.com/cweill/gotests/...@v1.6.0
go install github.com/fatih/gomodifytags@v1.17.0
go install github.com/josharian/impl@v1.4.0
go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0
go install github.com/go-delve/delve/cmd/dlv@v1.25.0
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

# go-zero toolkit
go install github.com/zeromicro/go-zero/tools/goctl@v1.8.4

go mod tidy

