#!/bin/bash
# AuraSpeed Installer for Linux/macOS
# Downloads and installs AuraSpeed to ~/.config/neostore/auraspeed/bin/

set -e

VERSION_URL="https://raw.githubusercontent.com/rkriad585/auraspeed/main/.version"
RELEASE_URL="https://github.com/rkriad585/auraspeed/releases/download"

# Get latest version from GitHub
VERSION=$(curl -sL "$VERSION_URL" | tr -d '[:space:]')
if [ -z "$VERSION" ]; then
    VERSION="v3.0.1"
fi

echo "AuraSpeed Installer ${VERSION}"
echo "=============================="

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
fi

# Set binary name based on OS
if [ "$OS" = "darwin" ]; then
    BIN_NAME="auraspeed"
else
    BIN_NAME="auraspeed"
fi

# Set download URL
if [ "$OS" = "darwin" ]; then
    DOWNLOAD_URL="${RELEASE_URL}/${VERSION}/auraspeed-darwin-${ARCH}"
else
    DOWNLOAD_URL="${RELEASE_URL}/${VERSION}/auraspeed-linux-${ARCH}"
fi

# Determine config directory
if [ -n "$XDG_CONFIG_HOME" ]; then
    BIN_DIR="$XDG_CONFIG_HOME/neostore/auraspeed/bin"
else
    BIN_DIR="$HOME/.config/neostore/auraspeed/bin"
fi

# Create bin directory
echo "Creating installation directory: $BIN_DIR"
mkdir -p "$BIN_DIR"

# Download binary
echo "Downloading AuraSpeed v$VERSION from: $DOWNLOAD_URL"
curl -L -o "$BIN_DIR/$BIN_NAME" "$DOWNLOAD_URL" || {
    echo "Error: Failed to download AuraSpeed"
    exit 1
}

# Make executable
chmod +x "$BIN_DIR/$BIN_NAME"

# Add to PATH
SHELL_RC=""
case "$SHELL" in
    */bash) SHELL_RC="$HOME/.bashrc" ;;
    */zsh) SHELL_RC="$HOME/.zshrc" ;;
    */fish) SHELL_RC="$HOME/.config/fish/config.fish" ;;
    *) SHELL_RC="$HOME/.profile" ;;
esac

# Check if already in PATH
if ! echo "$PATH" | grep -q "$BIN_DIR"; then
    echo "Adding to PATH..."
    if [ "$SHELL" = "fish" ]; then
        echo "set -gx PATH $BIN_DIR \$PATH" >> "$SHELL_RC"
    else
        echo "export PATH=\"\$PATH:$BIN_DIR\"" >> "$SHELL_RC"
    fi
    echo "PATH updated. Run 'source $SHELL_RC' or restart your terminal."
else
    echo "Already in PATH"
fi

echo ""
echo "Installation complete!"
echo ""
echo "Run 'auraspeed --help' to get started"
echo "To uninstall, run: auraspeed --selfuninstall"