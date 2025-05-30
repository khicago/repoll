#!/bin/bash

echo "=== Building project ==="
go build . 2>&1

echo "=== Compiling tests ==="
go test -c -o /dev/null 2>&1

echo "=== Running simple test ==="
go test -run TestGetRepoNameFromURL_Basic -v 2>&1

echo "=== Running all tests ==="
go test -v 2>&1 | head -50 