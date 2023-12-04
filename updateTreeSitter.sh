#!/bin/bash

# Script to update tree-sitter
# Exits immediately if a command exits with a non-zero status.
set -e

sitter_version=v0.20.8

# Function to download and set up tree-sitter
function download_sitter() {
    # Clean up the working directory
    rm -rf ./vendor
    rm -rf ./tmpts*
    # Tidy and vendor Go modules
    go mod tidy
    go mod vendor  

    # Remove any existing tree-sitter directory
    rm -rf tmpts
    # Clone the specified version of tree-sitter
    git clone -b $1 https://github.com/tree-sitter/tree-sitter.git tmpts --depth 1

    # Modify include paths in the source and header files for compatibility
    sed -i.bak 's/"tree_sitter\//"/g' tmpts/lib/src/*.c tmpts/lib/src/*.h
    sed -i.bak 's/"unicode\//"/g' tmpts/lib/src/unicode/*.h tmpts/lib/src/*.h

    # Create a new directory for modified tree-sitter files
    mkdir tmpts2

    # Copy the necessary tree-sitter files to the new directory
    cp tmpts/lib/include/tree_sitter/*.h ./tmpts2
    cp tmpts/lib/src/*.c ./tmpts2
    cp tmpts/lib/src/*.h ./tmpts2
    cp tmpts/lib/src/unicode/*.h ./tmpts2
    # Remove the lib.c file to avoid "duplicate symbols" errors
    rm ./tmpts2/lib.c

    # Clean up the original tree-sitter clone
    rm -rf tmpts

    # Copy modified files to the vendor directory
    cp -f tmpts2/*.c tmpts2/*.h ./vendor/github.com/smacker/go-tree-sitter/
    # Remove the temporary directory
    rm -rf tmpts2
}

# Call the function with the desired tree-sitter version
download_sitter $sitter_version
