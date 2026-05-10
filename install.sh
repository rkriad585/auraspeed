#!/bin/bash
# AuraSpeed Installer for Linux/macOS
# Downloads and installs AuraSpeed to ~/.config/neostore/auraspeed/bin/
# and adds it to PATH in the appropriate shell rc file.

set -e

PROJECT_NAME="auraspeed"
VERSION_URL="https://raw.githubusercontent.com/rkriad585/${PROJECT_NAME}/main/.version"
RELEASE_URL="https://github.com/rkriad585/${PROJECT_NAME}/releases/download"

# Get latest version from GitHub
VERSION=$(curl -sL "$VERSION_URL" | tr -d '[:space:]')
if [ -z "$VERSION" ]; then
    VERSION="dev"
fi

echo ">>> Installing ${PROJECT_NAME} ${VERSION}..."
echo ""

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Determine download file name
case "$OS" in
    linux) DOWNLOAD_NAME="${PROJECT_NAME}-linux-${ARCH}" ;;
    darwin) DOWNLOAD_NAME="${PROJECT_NAME}-darwin-${ARCH}" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

DOWNLOAD_URL="${RELEASE_URL}/${VERSION}/${DOWNLOAD_NAME}"

# Determine install directory
if [ -n "$XDG_CONFIG_HOME" ]; then
    BIN_DIR="$XDG_CONFIG_HOME/neostore/${PROJECT_NAME}/bin"
else
    BIN_DIR="$HOME/.config/neostore/${PROJECT_NAME}/bin"
fi

# Create bin directory
mkdir -p "$BIN_DIR"

# Download binary
echo "Downloading ${DOWNLOAD_URL}"
curl -sL -o "${BIN_DIR}/${PROJECT_NAME}" "$DOWNLOAD_URL" || {
    echo "Error: Failed to download ${PROJECT_NAME}"
    exit 1
}

# Make executable
chmod +x "${BIN_DIR}/${PROJECT_NAME}"

# Verify download
if [ ! -f "${BIN_DIR}/${PROJECT_NAME}" ]; then
    echo "Error: Download failed - output file not found"
    exit 1
fi

SIZE=$(du -k "${BIN_DIR}/${PROJECT_NAME}" | cut -f1)
echo "OK   Installed to ${BIN_DIR}/${PROJECT_NAME} (${SIZE} KB)"

# Add to PATH
SHELL_RC=""
case "$SHELL" in
    */bash) SHELL_RC="$HOME/.bashrc" ;;
    */zsh)  SHELL_RC="$HOME/.zshrc" ;;
    */fish) SHELL_RC="$HOME/.config/fish/config.fish" ;;
    *)      SHELL_RC="$HOME/.profile" ;;
esac

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
echo "OK   ${PROJECT_NAME} ${VERSION} installed successfully!"
echo "Run '${PROJECT_NAME} --help' to get started"
