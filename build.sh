#!/bin/bash
set -e

# Extract the most recent tag for this build
VERSION=$(git describe --tags)
GOPATH=$(gb env GB_PROJECT_DIR):$(gb env GB_PROJECT_DIR)/vendor

echo "building the core..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/autobus-core core/cmd/autobus-core
echo "building the platform..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/autobus-platform platform/cmd/autobus-platform
echo "building the web API..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/autobus-web web/cmd/autobus-web
