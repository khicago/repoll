#!/bin/bash

# 检查覆盖率的脚本
echo "Running coverage tests..."

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out -v ./...

# 检查退出码
if [ $? -eq 0 ]; then
    echo "Tests passed, generating coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    go tool cover -func=coverage.out
else
    echo "Tests failed, checking for compilation errors..."
    go build .
fi 