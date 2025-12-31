#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Running Tests with Coverage ===${NC}\n"

# Run tests with coverage
go test ./tests/... -v -coverprofile=coverage.out -covermode=atomic -coverpkg=./handler/...,./service/...,./repository/...,./middleware/...,./helper/...

# Check if tests passed
if [ $? -eq 0 ]; then
    echo -e "\n${GREEN}✓ All tests passed${NC}\n"
    
    # Generate coverage report
    echo -e "${YELLOW}=== Coverage Report ===${NC}\n"
    go tool cover -func=coverage.out
    
    # Get total coverage percentage
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo -e "\n${GREEN}Total Coverage: $COVERAGE${NC}\n"
    
    # Generate HTML report
    echo -e "${YELLOW}=== Generating HTML Report ===${NC}"
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✓ HTML report generated: coverage.html${NC}\n"
    
    # Check if coverage meets threshold (e.g., 60%)
    THRESHOLD=60
    COVERAGE_NUM=$(echo $COVERAGE | sed 's/%//')
    
    if (( $(echo "$COVERAGE_NUM >= $THRESHOLD" | bc -l) )); then
        echo -e "${GREEN}✓ Coverage meets threshold ($THRESHOLD%)${NC}"
        exit 0
    else
        echo -e "${RED}✗ Coverage below threshold: $COVERAGE_NUM% < $THRESHOLD%${NC}"
        exit 1
    fi
else
    echo -e "\n${RED}✗ Tests failed${NC}\n"
    exit 1
fi
