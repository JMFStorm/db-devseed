#!/bin/bash

BUILD_DIR="./build"
EXECUTABLE_NAME="dbds"

mkdir -p "$BUILD_DIR"

go build -o "$BUILD_DIR/$EXECUTABLE_NAME"

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Build successful: $BUILD_DIR/$EXECUTABLE_NAME"
else
    echo "Build failed"
    exit 1
fi
