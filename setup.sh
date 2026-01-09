#!/bin/bash

# Setup script for the learning-core-api project

set -euo pipefail

echo "Starting setup..."

# 0. Validate dependencies
required_bins=(go make)
missing=()
for bin in "${required_bins[@]}"; do
    if ! command -v "$bin" >/dev/null 2>&1; then
        missing+=("$bin")
    fi
done
if [ ${#missing[@]} -ne 0 ]; then
    echo "Missing required tools: ${missing[*]}"
    echo "Install them and rerun setup."
    exit 1
fi

# 0b. Validate environment (cloud containers often rely on explicit envs)
required_envs=(DB_URL JWT_SECRET)
missing_envs=()
for env_name in "${required_envs[@]}"; do
    if [ -z "${!env_name:-}" ]; then
        missing_envs+=("$env_name")
    fi
done
if [ ${#missing_envs[@]} -ne 0 ]; then
    echo "Missing required environment variables: ${missing_envs[*]}"
    echo "Set them (or create a .env) before running setup."
    exit 1
fi

# 1. Install dependencies
echo "Tidying Go modules..."
go mod tidy


# 3. Generate SQLC code
echo "Generating persistence layer..."
make sqlc

# 4. Optional migrations (set RUN_MIGRATIONS=1)
if [ "${RUN_MIGRATIONS:-}" = "1" ]; then
    echo "Running migrations..."
    make migrate-up
fi

echo "Setup complete!"
