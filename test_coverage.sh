#!/bin/bash

echo "=== Running Go tests with coverage ==="

# Clean up previous coverage files
rm -f coverage.out

# Run tests with coverage
echo "Running tests..."
go test -v -coverprofile=coverage.out ./... 2>&1

# Check if tests ran successfully
if [ $? -eq 0 ]; then
    echo "=== Tests completed successfully ==="
    
    # Generate coverage report
    if [ -f coverage.out ]; then
        echo "=== Coverage Report ==="
        go tool cover -func=coverage.out
        
        echo ""
        echo "=== Overall Coverage ==="
        go tool cover -func=coverage.out | tail -1
        
        # Generate HTML report
        go tool cover -html=coverage.out -o coverage.html
        echo "HTML coverage report generated: coverage.html"
    else
        echo "No coverage file generated"
    fi
else
    echo "=== Tests failed ==="
    echo "Trying to run individual test files to identify issues..."
    
    for test_file in *_test.go; do
        if [ -f "$test_file" ]; then
            echo "Testing $test_file..."
            go test -v -run ".*" "$test_file" 2>&1 | head -10
        fi
    done
fi 