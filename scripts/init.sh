#!/bin/bash
set -e

# Install Go tools
export GOPROXY=https://goproxy.cn,direct

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

go install golang.org/x/tools/gopls@latest
go install github.com/cweill/gotests/...@latest
go install github.com/fatih/gomodifytags@latest
go install github.com/josharian/impl@latest
go install github.com/haya14busa/goplay/cmd/goplay@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

go mod tidy

echo ": $(date +%s):0;go mod tidy" >> "$HOME"/.zsh_history
echo ": $(date +%s):0;go run main.go" >> "$HOME"/.zsh_history

# Init Python environment (for test)
uv sync --active