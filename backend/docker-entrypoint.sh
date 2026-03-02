#!/bin/sh
set -e

# Data directory (Cloudron localstorage defaults to /app/data)
DATA_DIR="${DATA_DIR:-/app/data}"

# Create directories if they don't exist
mkdir -p "$DATA_DIR/isos" "$DATA_DIR/db"

# Fix ownership - only if needed and don't follow symlinks
# Only change ownership of the directories themselves, not recursively
chown isoman:isoman "$DATA_DIR" "$DATA_DIR/isos" "$DATA_DIR/db"

# Switch to isoman user and execute the main command
exec su-exec isoman "$@"
