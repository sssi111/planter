#!/bin/bash

# This script runs tests for the Planter backend application
# It provides options for running specific tests, generating coverage reports,
# and running tests with race detection

# Set default values
VERBOSE=false
COVERAGE=false
RACE=false
PACKAGE="./..."

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -v|--verbose)
      VERBOSE=true
      shift
      ;;
    -c|--coverage)
      COVERAGE=true
      shift
      ;;
    -r|--race)
      RACE=true
      shift
      ;;
    -p|--package)
      PACKAGE="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [-v|--verbose] [-c|--coverage] [-r|--race] [-p|--package PACKAGE]"
      exit 1
      ;;
  esac
done

# Build the test command
TEST_CMD="go test"

if [ "$VERBOSE" = true ]; then
  TEST_CMD="$TEST_CMD -v"
fi

if [ "$RACE" = true ]; then
  TEST_CMD="$TEST_CMD -race"
fi

if [ "$COVERAGE" = true ]; then
  TEST_CMD="$TEST_CMD -coverprofile=coverage.out"
fi

TEST_CMD="$TEST_CMD $PACKAGE"

# Run the tests
echo "Running tests with command: $TEST_CMD"
eval $TEST_CMD

# Generate coverage report if requested
if [ "$COVERAGE" = true ]; then
  echo "Generating coverage report..."
  go tool cover -html=coverage.out -o coverage.html
  echo "Coverage report generated at coverage.html"
fi

echo -e "\nTest examples:"
echo "1. Run all tests:"
echo "   ./scripts/run_tests.sh"
echo "2. Run all tests with verbose output:"
echo "   ./scripts/run_tests.sh -v"
echo "3. Run all tests with coverage report:"
echo "   ./scripts/run_tests.sh -c"
echo "4. Run tests for a specific package:"
echo "   ./scripts/run_tests.sh -p ./internal/services"
echo "5. Run tests with race detection:"
echo "   ./scripts/run_tests.sh -r"
echo "6. Combine options:"
echo "   ./scripts/run_tests.sh -v -c -p ./internal/services"