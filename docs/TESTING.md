# Test Coverage Configuration

## Minimum Coverage Threshold
COVERAGE_THRESHOLD=70

## Coverage Goals
- **Overall Project**: 70%
- **Handlers**: 80%
- **Services**: 85%
- **Repositories**: 90%
- **Helpers**: 90%

## Running Tests with Coverage

### Quick Test
```bash
go test ./tests/... -v
```

### Test with Coverage Report
```bash
./scripts/test-coverage.sh
```

### Test Specific Package
```bash
go test ./tests/... -v -run TestAuthTestSuite
```

### Generate HTML Coverage Report
```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Coverage Report Details

The coverage report shows:
- **Green**: Lines covered by tests
- **Red**: Lines not covered by tests
- **Gray**: Lines not executed (comments, blank lines)

## Viewing Coverage

After running `./scripts/test-coverage.sh`:
1. Open `coverage.html` in your browser
2. Check the summary at the bottom for overall coverage percentage
3. Click on files to see line-by-line coverage

## CI/CD Integration

Add to your CI/CD pipeline:

```yaml
# GitHub Actions example
- name: Run tests with coverage
  run: |
    go test ./tests/... -coverprofile=coverage.out -covermode=atomic
    
- name: Check coverage threshold
  run: |
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    if (( $(echo "$COVERAGE < 70" | bc -l) )); then
      echo "Coverage $COVERAGE% is below threshold 70%"
      exit 1
    fi
```

## Excluded from Coverage

The following are typically excluded from coverage requirements:
- `main.go` (entry point)
- Generated code
- Third-party integrations
- Mock implementations

## Best Practices

1. **Write tests first** (TDD approach when possible)
2. **Test edge cases** not just happy paths
3. **Use table-driven tests** for similar test cases
4. **Mock external dependencies** (database, APIs, etc.)
5. **Keep tests fast** (use in-memory database)
6. **Run tests before commit** (pre-commit hook recommended)

## Testing Structure

```
tests/
├── setup_test.go          # Test infrastructure
├── auth_test.go           # Authentication tests
├── post_test.go           # Post CRUD tests
└── ...

mocks/
├── repository_mocks.go    # Repository mocks
├── service_mocks.go       # Service mocks
└── ...
```

## Test Naming Convention

- Suite: `{Feature}TestSuite`
- Test: `Test{Method}_{Scenario}`
- Example: `TestLogin_Success`, `TestLogin_InvalidCredentials`

## Running Tests in Watch Mode

For development, use Air (already configured):
```bash
air test
```

## Troubleshooting

### Tests failing with database errors
- Check if SQLite driver is installed: `go get gorm.io/driver/sqlite`
- Ensure database cleanup is working in `TearDownSuite`

### Coverage report not generated
- Check write permissions in project directory
- Ensure `go tool cover` is available in your Go installation

### Tests running slowly
- Use in-memory SQLite (already configured)
- Mock external services
- Run tests in parallel when possible
