#!/bin/sh
# entrypoint.sh

# Check that DATABASE_URL is provided
if [ -z "$DATABASE_URL" ]; then
  echo "DATABASE_URL is not set"
  exit 1
fi

# Apply migrations using the environment variable
atlas migrate apply \
  --dir file://ent/migrate/migrations \
  --url "$DATABASE_URL"
