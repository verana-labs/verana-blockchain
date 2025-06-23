#!/bin/bash

echo "Cleaning up Verana development environment..."

# Stop and remove containers
docker stop val1 val2 val3 2>/dev/null || true
docker rm val1 val2 val3 2>/dev/null || true

# Remove data directories
rm -rf val1 val2 val3

echo "✅ Cleanup complete!"