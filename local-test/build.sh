#!/bin/bash

# Simple build script for development
echo "Building Verana Docker image..."

docker build -f local-test/Dockerfile -t verana:dev .

if [ $? -eq 0 ]; then
    echo "✅ Build successful: verana:dev"
    echo "Image size: $(docker images verana:dev --format 'table {{.Size}}' | tail -1)"
else
    echo "❌ Build failed"
    exit 1
fi