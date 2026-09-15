# CinemaAbyss API Tests

This directory contains Postman tests for the CinemaAbyss microservices architecture. The tests are designed to be run using Newman, the command-line collection runner for Postman.

## Structure

- `CinemaAbyss.postman_collection.json` - The main Postman collection containing all API tests
- `local.environment.json` - Environment variables for running tests against locally running services
- `docker.environment.json` - Environment variables for running tests against services in Docker containers
- `run-tests.js` - Node.js script to run the tests using Newman
- `package.json` - Node.js package configuration with dependencies and scripts

## Test Coverage

The tests cover the following services:

1. **Monolith Service**
   - User management (create, get)
   - Movie management (create, get)
   - Payment processing (create, get)
   - Subscription management (create, get)

2. **Movies Microservice**
   - Health check
   - Movie management (create, get)

3. **Events Microservice**
   - Health check
   - Event publishing (movie, user, payment events)

4. **Proxy Service**
   - Health check
   - Proxying requests to other services

## Prerequisites

- Node.js (v14 or later)
- npm (v6 or later)

## Installation

```bash
# Navigate to the tests directory
cd tests/postman

# Install dependencies
npm install
```

## Running Tests

### Basic Usage

```bash
# Run all tests against the local environment (default)
npm test

# Run all tests against the Docker environment
npm run test:docker
```

### Running Specific Test Folders

```bash
# Run only Monolith Service tests
npm run test:monolith

# Run only Movies Microservice tests
npm run test:movies

# Run only Events Microservice tests
npm run test:events

# Run only Proxy Service tests
npm run test:proxy
```

### Advanced Usage

The `run-tests.js` script supports several command-line options:

```bash
node run-tests.js --environment <env> --folder <folder> --reporters <reporters> --bail --timeout <ms>
```

Options:
- `--environment`, `-e`: Environment to run tests against (default: 'local')
- `--collection`, `-c`: Collection to run (default: 'CinemaAbyss')
- `--folder`, `-f`: Specific folder in the collection to run
- `--reporters`, `-r`: Reporters to use, comma-separated (default: 'cli,htmlextra,junit')
- `--bail`, `-b`: Stop on first error (default: false)
- `--timeout`, `-t`: Request timeout in ms (default: 10000)

Example:
```bash
node run-tests.js --environment docker --folder "Movies Microservice" --reporters cli,htmlextra --bail
```

## Test Reports

After running the tests, HTML and JUnit XML reports will be generated in the `reports` directory. These reports can be used for CI/CD integration and documentation.

## CI/CD Integration

These tests can be integrated into CI/CD pipelines. Here's an example of how to run them in a GitHub Actions workflow:

```yaml
- name: Run API Tests
  run: |
    cd tests/postman
    npm install
    npm run test:docker
```

## Troubleshooting

If you encounter issues running the tests:

1. Ensure all services are running and accessible
2. Check the environment configuration in the environment JSON files
3. Verify that the API endpoints match those in the collection
4. Increase the timeout value if requests are timing out