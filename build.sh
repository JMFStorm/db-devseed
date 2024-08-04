#!/bin/bash

BUILD_DIR="./build"
EXECUTABLE_NAME_LINUX="dbds"
EXECUTABLE_NAME_WINDOWS="dbds.exe"

# Create the build directory if it does not exist
mkdir -p "$BUILD_DIR"

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o "$BUILD_DIR/$EXECUTABLE_NAME_LINUX"
if [ $? -eq 0 ]; then
    echo "Build successful for Linux: $BUILD_DIR/$EXECUTABLE_NAME_LINUX"
else
    echo "Build failed for Linux"
    exit 1
fi

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o "$BUILD_DIR/$EXECUTABLE_NAME_WINDOWS"
if [ $? -eq 0 ]; then
    echo "Build successful for Windows: $BUILD_DIR/$EXECUTABLE_NAME_WINDOWS"
else
    echo "Build failed for Windows"
    exit 1
fi
