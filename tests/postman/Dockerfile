# Test-runner image for the CinemaAbyss API test suite (Newman).
#
# Self-contained: newman + reporters are baked in, so the same image works
# locally and in CI without any host-side npm install. The collection,
# environment files and run-tests.js are mounted at runtime (see the CI
# workflow), which keeps the image small and always in sync with the repo.
FROM node:18-alpine

WORKDIR /app

# Install newman and the reporters used by run-tests.js (junit for CI
# artifacts; htmlextra is optional and not required in CI).
RUN npm install --no-audit --no-fund \
      newman@^4.6.1 \
      yargs@^17.6.2 \
      newman-reporter-junitfull@^1.1.1

# NOTE: /app/reports is intentionally NOT created here. In CI the source dir
# is bind-mounted over /app and a separate volume is mounted at /app/reports;
# Docker creates that mountpoint for us. Pre-creating it would be shadowed by
# the bind mount anyway.

# Default command - can be overridden at runtime.
ENTRYPOINT ["node", "run-tests.js"]
CMD ["--environment", "docker"]