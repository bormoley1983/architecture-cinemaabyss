# Test-runner image for the CinemaAbyss API test suite (Newman).
#
# Self-contained: newman + reporters AND the collection/environment files are
# baked in, so the same image works locally and in CI without any host-side
# npm install or source bind mount. Baking the files in (rather than mounting
# them) avoids shadowing /app/node_modules when a volume is mounted at /app.
FROM node:18-alpine

WORKDIR /app

# Install newman and the reporters used by run-tests.js (junit for CI
# artifacts; htmlextra is optional and not required in CI).
RUN npm install --no-audit --no-fund \
      newman@^4.6.1 \
      yargs@^17.6.2 \
      newman-reporter-junitfull@^1.1.1

# Bake in the test assets (runner script, collection, environment files).
COPY run-tests.js ./
COPY *.postman_collection.json ./
COPY *.environment.json ./

# Default command - can be overridden at runtime. Only /app/reports is mounted
# at run time (for the JUnit artifact); everything else lives in the image.
ENTRYPOINT ["node", "run-tests.js"]
CMD ["--environment", "docker"]