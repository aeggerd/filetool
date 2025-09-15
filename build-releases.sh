#!/bin/bash

# Release build script for multiple platforms and configurations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

VERSION="${VERSION:-$(date +%Y%m%d)}"
RELEASE_DIR="releases"

echo -e "${BLUE}🚀 Building filetool releases${NC}"
echo -e "📦 Version: ${YELLOW}$VERSION${NC}"
echo ""

# Create release directory
mkdir -p "$RELEASE_DIR"

# Build configurations
declare -A BUILDS=(
    # Platform configurations: "os:arch:cgo:build_type:suffix"
    ["linux-amd64-cli"]="linux:amd64:0:minimal:-linux-amd64-cli"
    ["linux-arm64-cli"]="linux:arm64:0:minimal:-linux-arm64-cli"
    ["windows-amd64-cli"]="windows:amd64:0:minimal:-windows-amd64-cli"
    ["darwin-amd64-cli"]="darwin:amd64:0:minimal:-macos-amd64-cli"
    ["darwin-arm64-cli"]="darwin:arm64:0:minimal:-macos-arm64-cli"
)

# Counter for builds
TOTAL_BUILDS=${#BUILDS[@]}
CURRENT_BUILD=0
SUCCESS_COUNT=0
FAILED_BUILDS=()

echo -e "${BLUE}📊 Building ${TOTAL_BUILDS} CLI configurations...${NC}"
echo ""

for BUILD_NAME in "${!BUILDS[@]}"; do
    CURRENT_BUILD=$((CURRENT_BUILD + 1))
    
    # Parse build configuration
    IFS=':' read -ra CONFIG <<< "${BUILDS[$BUILD_NAME]}"
    OS="${CONFIG[0]}"
    ARCH="${CONFIG[1]}" 
    CGO="${CONFIG[2]}"
    BUILD_TYPE="${CONFIG[3]}"
    SUFFIX="${CONFIG[4]}"
    
    OUTPUT_NAME="filetool$SUFFIX"
    if [ "$OS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo -e "${BLUE}[$CURRENT_BUILD/$TOTAL_BUILDS]${NC} Building ${YELLOW}$BUILD_NAME${NC}..."
    echo -e "   Target: $OS/$ARCH, Type: $BUILD_TYPE, CGO: $CGO"
    
    # Set environment variables for build script
    export GOOS="$OS"
    export GOARCH="$ARCH"
    export BUILD_TYPE="$BUILD_TYPE" 
    export OUTPUT_NAME="$RELEASE_DIR/$OUTPUT_NAME"
    
    # Build with appropriate CGO setting
    if ./build.sh > /dev/null 2>&1; then
        echo -e "   ${GREEN}✅ Success${NC} - $(ls -lh "$RELEASE_DIR/$OUTPUT_NAME" | awk '{print $5}')"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        echo -e "   ${RED}❌ Failed${NC}"
        FAILED_BUILDS+=("$BUILD_NAME")
    fi
    echo ""
done

# Summary
echo -e "${BLUE}📋 Build Summary${NC}"
echo -e "✅ Successful: ${GREEN}$SUCCESS_COUNT/$TOTAL_BUILDS${NC}"

if [ ${#FAILED_BUILDS[@]} -gt 0 ]; then
    echo -e "❌ Failed: ${RED}${#FAILED_BUILDS[@]}${NC}"
    echo -e "${RED}Failed builds:${NC}"
    for BUILD in "${FAILED_BUILDS[@]}"; do
        echo -e "  - $BUILD"
    done
fi

echo ""
echo -e "${BLUE}📁 Release files:${NC}"
ls -lh "$RELEASE_DIR/"

# Create checksums
echo ""
echo -e "${BLUE}🔐 Generating checksums...${NC}"
cd "$RELEASE_DIR"
sha256sum * > SHA256SUMS
echo -e "✅ Checksums saved to ${YELLOW}$RELEASE_DIR/SHA256SUMS${NC}"

echo ""
echo -e "${GREEN}🎉 Release build complete!${NC}"
echo -e "📦 Files available in: ${YELLOW}$RELEASE_DIR/${NC}"