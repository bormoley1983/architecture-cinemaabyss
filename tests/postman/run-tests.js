const newman = require('newman');
const fs = require('fs');
const path = require('path');
const yargs = require('yargs/yargs');
const { hideBin } = require('yargs/helpers');

// Parse command line arguments
const argv = yargs(hideBin(process.argv))
  .option('environment', {
    alias: 'e',
    description: 'Environment to run tests against',
    type: 'string',
    default: 'local'
  })
  .option('collection', {
    alias: 'c',
    description: 'Collection to run',
    type: 'string',
    default: 'CinemaAbyss'
  })
  .option('folder', {
    alias: 'f',
    description: 'Specific folder in the collection to run',
    type: 'string'
  })
  .option('reporters', {
    alias: 'r',
    description: 'Reporters to use (comma-separated)',
    type: 'string',
    default: 'cli,htmlextra,junit'
  })
  .option('bail', {
    alias: 'b',
    description: 'Stop on first error',
    type: 'boolean',
    default: false
  })
  .option('timeout', {
    alias: 't',
    description: 'Request timeout in ms',
    type: 'number',
    default: 10000
  })
  .help()
  .alias('help', 'h')
  .argv;

// Create reports directory if it doesn't exist
const reportsDir = path.join(__dirname, 'reports');
if (!fs.existsSync(reportsDir)) {
  fs.mkdirSync(reportsDir, { recursive: true });
}

// Configure Newman run
const collectionPath = path.join(__dirname, `${argv.collection}.postman_collection.json`);
const environmentPath = path.join(__dirname, `${argv.environment}.environment.json`);

// Validate files exist
if (!fs.existsSync(collectionPath)) {
  console.error(`Collection file not found: ${collectionPath}`);
  process.exit(1);
}

if (!fs.existsSync(environmentPath)) {
  console.error(`Environment file not found: ${environmentPath}`);
  process.exit(1);
}

// Parse reporters
const reporters = argv.reporters.split(',').map(r => r.trim());

// Configure Newman options
const newmanOptions = {
  collection: require(collectionPath),
  environment: require(environmentPath),
  reporters: reporters,
  reporter: {
    htmlextra: {
      export: path.join(reportsDir, `report-${argv.environment}-${new Date().toISOString().replace(/:/g, '-')}.html`),
      template: 'default',
      showOnlyFails: false,
      noSyntaxHighlighting: false,
      testPaging: true,
      browserTitle: "CinemaAbyss API Test Report",
      title: "CinemaAbyss API Test Report",
      titleSize: 1,
      omitHeaders: false
    },
    junit: {
      export: path.join(reportsDir, `junit-report-${argv.environment}-${new Date().toISOString().replace(/:/g, '-')}.xml`)
    }
  },
  bail: argv.bail,
  timeoutRequest: argv.timeout,
  delayRequest: 100 // Small delay between requests
};

// Add folder option if specified
if (argv.folder) {
  newmanOptions.folder = argv.folder;
}

// Run Newman
console.log(`Running tests against ${argv.environment} environment...`);
newman.run(newmanOptions, function (err, summary) {
  if (err) { 
    console.error('Error running Newman:', err);
    process.exit(1);
  }
  
  // Log results
  console.log('Newman run completed!');
  
  const failureCount = summary.run.failures.length;
  console.log(`Total requests: ${summary.run.stats.requests.total}`);
  console.log(`Failed requests: ${summary.run.stats.requests.failed}`);
  console.log(`Total assertions: ${summary.run.stats.assertions.total}`);
  console.log(`Failed assertions: ${summary.run.stats.assertions.failed}`);
  
  // Exit with appropriate code
  process.exit(failureCount > 0 ? 1 : 0);
});