#!/bin/bash
set -e

# Extract the most recent tag for this build
VERSION=$(git describe --tags)

# export the GOPATH to somewhere we can actually use the go tool on
GOPATH=$(gb env GB_PROJECT_DIR):$(gb env GB_PROJECT_DIR)/vendor

echo "building the core..."
CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=$VERSION" -a -installsuffix cgo -o bin/autobus-core core/cmd/autobus-core
echo "building the platform..."
CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=$VERSION" -a -installsuffix cgo -o bin/autobus-platform platform/cmd/autobus-platform
echo "building the web API..."
CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=$VERSION" -a -installsuffix cgo -o bin/autobus-web web/cmd/autobus-web
