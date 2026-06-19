#!/bin/bash

set -e

clear
# Generate sql
echo "Building sql queries with sqlc"
sqlc generate
echo "Success"

# Build the Go application
echo "Building..."
if go build -o chirpy-server; then
    echo "Build successful. Running application..."
    echo "----------------------------------------"
    # Run the executable
    ./chirpy-server
else
    echo "Build failed. Execution aborted."
    exit 1
fi