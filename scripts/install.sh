#!/bin/bash

# netmon installation script
# Works on Linux and macOS

set -e

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Version information
VERSION=${1:-"latest"}
REPO="zzzzseong/netmon"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="netmon"

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    case $OS in
        linux|darwin)
            PLATFORM="${OS}-${ARCH}"
            ;;
        *)
            echo -e "${RED}❌ Unsupported operating system: $OS${NC}"
            exit 1
            ;;
    esac
    
    echo -e "${BLUE}📦 Detected platform: ${PLATFORM}${NC}"
}

# Get latest version
get_latest_version() {
    if [ "$VERSION" = "latest" ]; then
        echo -e "${BLUE}🔍 Checking for latest version...${NC}"
        VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi
    
    # Remove v prefix
    VERSION_NUM=$(echo $VERSION | sed 's/^v//')
    echo -e "${GREEN}✅ Version: ${VERSION}${NC}"
}

# Download binary
download_binary() {
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/netmon-${PLATFORM}.tar.gz"
    
    echo -e "${BLUE}⬇️  Downloading: ${DOWNLOAD_URL}${NC}"
    
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    trap "rm -rf $TEMP_DIR" EXIT
    
    # Download
    if ! curl -L -f -o "${TEMP_DIR}/netmon.tar.gz" "$DOWNLOAD_URL"; then
        echo -e "${RED}❌ Download failed. Please check if the version exists.${NC}"
        exit 1
    fi
    
    # Extract
    echo -e "${BLUE}📦 Extracting...${NC}"
    tar -xzf "${TEMP_DIR}/netmon.tar.gz" -C "${TEMP_DIR}"
    
    # Check install directory
    if [ ! -d "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}⚠️  Install directory does not exist. Creating...${NC}"
        sudo mkdir -p "$INSTALL_DIR"
    fi
    
    # Install binary
    echo -e "${BLUE}📥 Installing to: ${INSTALL_DIR}/${BINARY_NAME}${NC}"
    sudo mv "${TEMP_DIR}/netmon" "${INSTALL_DIR}/${BINARY_NAME}"
    sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    echo -e "${GREEN}✅ Installation complete!${NC}"
}

# Verify installation
verify_installation() {
    # Refresh command hash to ensure the new binary is found
    hash -r 2>/dev/null || true
    
    # Try to find the binary
    if command -v $BINARY_NAME &> /dev/null; then
        INSTALLED_VERSION=$($BINARY_NAME version 2>/dev/null | head -n 1 || echo "unknown")
        echo -e "${GREEN}✅ Installation verified: ${INSTALLED_VERSION}${NC}"
        echo -e "${GREEN}🎉 netmon has been successfully installed!${NC}"
        echo -e "${BLUE}💡 You can now use: ${BINARY_NAME} ls${NC}"
        echo -e "${BLUE}💡 Or run: ${BINARY_NAME} --help${NC}"
    elif [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        # Binary exists but not in PATH (shouldn't happen with /usr/local/bin)
        echo -e "${GREEN}✅ Installation complete!${NC}"
        echo -e "${YELLOW}⚠️  Note: You may need to open a new terminal or run: export PATH=\$PATH:${INSTALL_DIR}${NC}"
        echo -e "${BLUE}💡 You can also run directly: ${INSTALL_DIR}/${BINARY_NAME} ls${NC}"
    else
        echo -e "${RED}❌ Installation failed. Binary not found.${NC}"
        exit 1
    fi
}

# Main execution
main() {
    echo -e "${BLUE}"
    echo "╔════════════════════════════════════╗"
    echo "║      netmon Installer              ║"
    echo "╚════════════════════════════════════╝"
    echo -e "${NC}"
    
    detect_platform
    get_latest_version
    download_binary
    verify_installation
}

# Run script
main

