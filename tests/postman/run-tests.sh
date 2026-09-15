#!/bin/bash

# Script to run Postman tests using Newman

# Function to display usage information
function show_usage {
  echo "Usage: $0 [options]"
  echo ""
  echo "Options:"
  echo "  -e, --environment ENV   Specify environment (local, docker)"
  echo "  -f, --folder FOLDER     Run specific test folder"
  echo "  -r, --reporters LIST    Comma-separated list of reporters"
  echo "  -b, --bail              Stop on first error"
  echo "  -t, --timeout MS        Request timeout in milliseconds"
  echo "  -d, --docker            Run tests in Docker container"
  echo "  -h, --help              Show this help message"
  echo ""
  echo "Examples:"
  echo "  $0 -e local                     # Run all tests against local environment"
  echo "  $0 -e docker -f \"Movies Microservice\"  # Run only Movies tests against Docker environment"
  echo "  $0 -d -e docker                 # Run tests in Docker container against Docker environment"
}

# Default values
ENVIRONMENT="local"
FOLDER=""
REPORTERS="cli,htmlextra,junit"
BAIL=false
TIMEOUT=10000
USE_DOCKER=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -e|--environment)
      ENVIRONMENT="$2"
      shift 2
      ;;
    -f|--folder)
      FOLDER="$2"
      shift 2
      ;;
    -r|--reporters)
      REPORTERS="$2"
      shift 2
      ;;
    -b|--bail)
      BAIL=true
      shift
      ;;
    -t|--timeout)
      TIMEOUT="$2"
      shift 2
      ;;
    -d|--docker)
      USE_DOCKER=true
      shift
      ;;
    -h|--help)
      show_usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      show_usage
      exit 1
      ;;
  esac
done

# Ensure we're in the right directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Check if Node.js and npm are installed if not using Docker
if [ "$USE_DOCKER" = false ]; then
  if ! command -v node &> /dev/null; then
    echo "Error: Node.js is not installed. Please install Node.js or use the -d option to run in Docker."
    exit 1
  fi

  if ! command -v npm &> /dev/null; then
    echo "Error: npm is not installed. Please install npm or use the -d option to run in Docker."
    exit 1
  fi
fi

# Build command arguments
CMD_ARGS="--environment $ENVIRONMENT"

if [ -n "$FOLDER" ]; then
  CMD_ARGS="$CMD_ARGS --folder \"$FOLDER\""
fi

if [ -n "$REPORTERS" ]; then
  CMD_ARGS="$CMD_ARGS --reporters $REPORTERS"
fi

if [ "$BAIL" = true ]; then
  CMD_ARGS="$CMD_ARGS --bail"
fi

if [ -n "$TIMEOUT" ]; then
  CMD_ARGS="$CMD_ARGS --timeout $TIMEOUT"
fi

# Create reports directory if it doesn't exist
mkdir -p reports

# Run tests
if [ "$USE_DOCKER" = true ]; then
  echo "Running tests in Docker container..."
  
  # Build the Docker image
  docker build -t cinemaabyss-api-tests .
  
  # Run the tests in Docker
  docker run --network=cinemaabyss-network \
    -v "$(pwd)/reports:/app/reports" \
    cinemaabyss-api-tests $CMD_ARGS
else
  echo "Running tests locally..."
  
  # Install dependencies if node_modules doesn't exist
  if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install
  fi
  
  # Run the tests
  eval "node run-tests.js $CMD_ARGS"
fi

# Get the exit code
EXIT_CODE=$?

# Display results
if [ $EXIT_CODE -eq 0 ]; then
  echo "✅ All tests passed!"
else
  echo "❌ Some tests failed. Check the reports for details."
fi

exit $EXIT_CODE