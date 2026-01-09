#!/bin/bash

# Setup script for the learning-core-api project

set -euo pipefail

echo "Starting setup..."

# 0. Go version requirements
REQUIRED_GO_VERSION="1.25.5"

# 0. Validate/install dependencies
required_bins=(go make)
missing=()
for bin in "${required_bins[@]}"; do
    if ! command -v "$bin" >/dev/null 2>&1; then
        missing+=("$bin")
    fi
done

if [ ${#missing[@]} -ne 0 ]; then
    echo "Missing required tools: ${missing[*]}"
    if command -v apt-get >/dev/null 2>&1; then
        echo "Installing dependencies with apt-get..."
        apt-get update
        apt-get install -y make ca-certificates curl tar

        echo "Installing Go ${REQUIRED_GO_VERSION} from official tarball..."
        os="$(uname -s | tr '[:upper:]' '[:lower:]')"
        arch="$(uname -m)"
        case "$arch" in
            x86_64|amd64) arch="amd64" ;;
            aarch64|arm64) arch="arm64" ;;
            *) echo "Unsupported architecture: ${arch}"; exit 1 ;;
        esac
        go_tarball="go${REQUIRED_GO_VERSION}.${os}-${arch}.tar.gz"
        curl -fsSL "https://go.dev/dl/${go_tarball}" -o "/tmp/${go_tarball}"
        rm -rf /usr/local/go
        tar -C /usr/local -xzf "/tmp/${go_tarball}"
        export PATH="/usr/local/go/bin:${PATH}"
    elif command -v brew >/dev/null 2>&1; then
        echo "Installing dependencies with brew..."
        brew install go make
    else
        echo "No supported package manager found (apt-get/brew). Install missing tools and rerun setup."
        exit 1
    fi
fi

# 0c. Verify Go version matches go.mod
go_version_raw="$(go version | awk '{print $3}')"
go_version="${go_version_raw#go}"
if [ "$go_version" != "$REQUIRED_GO_VERSION" ]; then
    echo "Go version mismatch. Required: ${REQUIRED_GO_VERSION}, found: ${go_version}"
    echo "Install Go ${REQUIRED_GO_VERSION} and rerun setup."
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
