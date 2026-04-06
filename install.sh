#!/bin/bash

# Define variables
BINARY_NAME="devtool"
INSTALL_DIR="$HOME/.local/bin"


# Check if mise is installed and run mise install
if command -v mise &> /dev/null; then
    echo "mise found. Installing dependencies..."
    mise install
else
    echo "mise not found. Skipping 'mise install'. Ensure Go is available."
fi

echo "Building $BINARY_NAME..."
CGO_ENABLED=0 go build -o $BINARY_NAME main.go

if [ $? -ne 0 ]; then
    echo "Build failed! Please check your Go code."
    exit 1
fi

echo "Build successful."

if [ ! -d "$INSTALL_DIR" ]; then
    echo "Creating install directory at $INSTALL_DIR..."
    mkdir -p "$INSTALL_DIR"
fi

echo "Installing $BINARY_NAME to $INSTALL_DIR..."
mv "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

if [ $? -eq 0 ]; then
    echo "$BINARY_NAME installed successfully to $INSTALL_DIR"
    case ":$PATH:" in
        *":$INSTALL_DIR:"*)
            echo "You can now run '$BINARY_NAME' from anywhere."
            ;;
        *)
            echo "$INSTALL_DIR is not currently in your PATH."
            echo "Add this line to your shell profile to use '$BINARY_NAME' globally:"
            echo "export PATH=\"$INSTALL_DIR:\$PATH\""
            ;;
    esac
else
    echo "Installation failed!"
    exit 1
fi
