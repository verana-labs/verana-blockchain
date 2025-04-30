#!/bin/bash
set -e

COVERAGE_DIR=$(mktemp -d)
COVERAGE_PROFILE="$COVERAGE_DIR/coverage.out"
COVERAGE_HTML="coverage.html"

go test -mod=readonly -coverprofile="$COVERAGE_PROFILE" -covermode=atomic ./...
go tool cover -html="$COVERAGE_PROFILE" -o "$COVERAGE_HTML"

echo "Coverage report generated at $COVERAGE_HTML"