#!/bin/bash

# Define variables
BINARY_NAME="devtool"
INSTALL_DIR="/usr/local/bin"

echo "Building $BINARY_NAME..."
go build -o $BINARY_NAME main.go

if [ $? -ne 0 ]; then
    echo "Build failed! Please check your Go code."
    exit 1
fi

echo "Build successful."

echo "Installing $BINARY_NAME to $INSTALL_DIR..."
# Check if we have write permission to the install directory
if [ -w "$INSTALL_DIR" ]; then
    mv "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Sudo privileges required to move binary to $INSTALL_DIR"
    sudo mv "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

if [ $? -eq 0 ]; then
    echo "$BINARY_NAME installed successfully to $INSTALL_DIR"
    echo "You can now run '$BINARY_NAME' from anywhere."
else
    echo "Installation failed!"
    exit 1
fi
