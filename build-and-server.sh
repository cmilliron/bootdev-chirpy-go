#!/bin/bash

clear
# Build the Go application
echo "Building..."
if go build -o out; then
    echo "Build successful. Running application..."
    echo "----------------------------------------"
    # Run the executable
    ./out
else
    echo "Build failed. Execution aborted."
    exit 1
fi