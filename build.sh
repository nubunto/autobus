#!/bin/bash
set -e

# Extract the most recent tag for this build
VERSION=$(git describe --tags)

# Build the files, without cache, with the version above
GOOS=linux gb build -f -F -ldflags "-X main.Version=$VERSION"
