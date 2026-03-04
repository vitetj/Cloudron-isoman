#!/bin/sh
set -eu

export PORT=8080
export DATA_DIR="/app/data"
export GIN_MODE="release"
export CORS_ORIGINS="${CLOUDRON_APP_ORIGIN:-http://localhost:8080}"

# Cloudron proxyAuth handles global authentication at reverse-proxy level.
# Keep BASIC_AUTH_* only when create-only auth is explicitly enabled.
if [ "${CREATE_ISO_AUTH_ENABLED:-false}" != "true" ]; then
  unset BASIC_AUTH_USERNAME || true
  unset BASIC_AUTH_PASSWORD || true
fi

mkdir -p \
  /app/data/db \
  /app/data/isos \
  /app/data/isos/.tmp

cd /app/code
exec /app/code/server