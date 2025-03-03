#!/bin/bash
set -e

# Base URL for the binaries
BASE_URL="https://criblpreptool.s3.amazonaws.com/cribl-storage-tool/latest"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
fi

# Map OS names
if [ "$OS" = "darwin" ]; then
    BIN_PATH="darwin-$ARCH/cribl-storage-tool"
    INSTALL_DIR="/usr/local/bin"
    if [ -d "/opt/homebrew/bin" ] && [ "$ARCH" = "arm64" ]; then
        # Apple Silicon Macs with Homebrew often use this path
        INSTALL_DIR="/opt/homebrew/bin"
    fi
elif [ "$OS" = "linux" ]; then
    BIN_PATH="linux-$ARCH/cribl-storage-tool"
    INSTALL_DIR="/usr/local/bin"
elif [ "$OS" = "freebsd" ]; then
    BIN_PATH="freebsd-amd64/cribl-storage-tool"
    INSTALL_DIR="/usr/local/bin"
elif [[ "$OS" == "mingw"* ]] || [[ "$OS" == "cygwin"* ]] || [ "$OS" = "windows" ]; then
    BIN_PATH="windows-amd64/cribl-storage-tool.exe"
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
    echo "Windows detected, will install to $INSTALL_DIR"
else
    echo "Unsupported operating system: $OS"
    exit 1
fi

# Full URL to download
DOWNLOAD_URL="$BASE_URL/$BIN_PATH"
TMP_FILE="/tmp/cribl-storage-tool-tmp"

echo "Detected system: $OS $ARCH"
echo "Downloading from: $DOWNLOAD_URL"

# Download the binary
curl -L "$DOWNLOAD_URL" -o "$TMP_FILE"

# Make it executable
chmod +x "$TMP_FILE"

# For macOS, remove the quarantine attribute
if [ "$OS" = "darwin" ]; then
    echo "Removing quarantine attribute for macOS"
    xattr -d com.apple.quarantine "$TMP_FILE" 2>/dev/null || true
fi

# Install the binary
if [ "$OS" = "windows" ] || [[ "$OS" == "mingw"* ]] || [[ "$OS" == "cygwin"* ]]; then
    # On Windows, we don't use sudo
    mv "$TMP_FILE" "$INSTALL_DIR/cribl-storage-tool.exe"
    echo "Installed to $INSTALL_DIR/cribl-storage-tool.exe"

    # Add to PATH if not already there
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo "You may need to add $INSTALL_DIR to your PATH."
    fi
else
    # Use sudo on Unix systems
    echo "Installing to $INSTALL_DIR/cribl-storage-tool"
    sudo mv "$TMP_FILE" "$INSTALL_DIR/cribl-storage-tool"
    echo "Installation complete!"
fi

# Verify installation
if [ "$OS" = "windows" ] || [[ "$OS" == "mingw"* ]] || [[ "$OS" == "cygwin"* ]]; then
    if [ -f "$INSTALL_DIR/cribl-storage-tool.exe" ]; then
        echo "Cribl Storage Tool has been successfully installed."
    else
        echo "Installation failed."
    fi
else
    if command -v cribl-storage-tool >/dev/null 2>&1; then
        echo "Cribl Storage Tool has been successfully installed."
        echo "You can now run 'cribl-storage-tool' from anywhere."
    else
        echo "Installation seems to have succeeded, but the tool is not in your PATH."
        echo "You may need to restart your terminal or add $INSTALL_DIR to your PATH."
    fi
fi