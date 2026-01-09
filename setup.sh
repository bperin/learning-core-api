#!/bin/bash

# Setup script for the learning-core-api project

set -e

echo "Starting setup..."

# 1. Install dependencies
echo "Tidying Go modules..."
go mod tidy

# 2. Setup environment
if [ ! -f .env ]; then
    echo "Creating .env from .env.example..."
    cp .env.example .env
else
    echo ".env already exists, skipping creation."
fi

# 3. Generate SQLC code
echo "Generating persistence layer..."
make sqlc

echo "Setup complete!"
