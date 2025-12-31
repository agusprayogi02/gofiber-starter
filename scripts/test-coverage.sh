#!/bin/bash

# Script to run tests with coverage

echo "Running tests with coverage..."

# Run tests and generate coverage profile
go test ./tests/... -v -coverprofile=coverage.out -covermode=atomic

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Display coverage summary
echo ""
echo "Coverage Summary:"
go tool cover -func=coverage.out | tail -n 1

echo ""
echo "Coverage report generated: coverage.html"
echo "Total coverage profile: coverage.out"
