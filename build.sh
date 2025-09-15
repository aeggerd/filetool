#!/bin/bash

# Build script for filetool with multiple configuration support

set -e

# Default values
GOOS="${GOOS:-linux}"
GOARCH="${GOARCH:-amd64}"
OUTPUT_NAME="${OUTPUT_NAME:-filetool}"
BUILD_TYPE="${BUILD_TYPE:-minimal}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color


# sudo apt install gcc libgl1-mesa-dev xorg-dev
print_help() {
    echo "🔧 Filetool Build Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Environment Variables:"
    echo "  BUILD_TYPE    Build configuration: 'minimal' only (default: minimal)"
    echo "  GOOS         Target OS: linux, windows, darwin (default: linux)"
    echo "  GOARCH       Target architecture: amd64, arm64, 386 (default: amd64)"
    echo "  OUTPUT_NAME  Output binary name (default: filetool)"
    echo ""
    echo "Build Types:"
    echo "  minimal  - CLI only (default and only option)"
    echo ""
    echo "Examples:"
    echo "  ./build.sh                              # Default CLI-only build"
    echo "  GOOS=windows GOARCH=amd64 ./build.sh   # Windows build"
    echo "  GOOS=darwin GOARCH=arm64 ./build.sh    # macOS ARM build"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            print_help
            exit 0
            ;;
        *)
            echo -e "${RED}Error: Unknown parameter $1${NC}"
            print_help
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🔧 Building filetool...${NC}"
echo -e "📦 Configuration: ${YELLOW}$BUILD_TYPE${NC}"
echo -e "🎯 Target: ${YELLOW}$GOOS/$GOARCH${NC}"
echo -e "📁 Output: ${YELLOW}$OUTPUT_NAME${NC}"
echo ""

# Set build flags based on build type
BUILD_TAGS=""
CGO_ENABLED="0"
LDFLAGS="-s -w"

case $BUILD_TYPE in
    minimal)
        echo -e "${GREEN}⚡ Building CLI-only version${NC}"
        BUILD_TAGS="noui,noweb"
        CGO_ENABLED="0"
        ;;
    *)
        echo -e "${RED}Error: Unknown build type '$BUILD_TYPE'${NC}"
        echo "Valid types: minimal"
        exit 1
        ;;
esac

# Add platform-specific output extension
if [ "$GOOS" = "windows" ]; then
    OUTPUT_NAME="${OUTPUT_NAME}.exe"
fi

# Build the application
echo -e "${BLUE}⚙️  Compiling...${NC}"

if [ -n "$BUILD_TAGS" ]; then
    echo -e "🏷️  Build tags: ${YELLOW}$BUILD_TAGS${NC}"
    CGO_ENABLED=$CGO_ENABLED GOOS=$GOOS GOARCH=$GOARCH go build \
        -tags="$BUILD_TAGS" \
        -ldflags="$LDFLAGS" \
        -o "$OUTPUT_NAME"
else
    CGO_ENABLED=$CGO_ENABLED GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="$LDFLAGS" \
        -o "$OUTPUT_NAME"
fi

# Check if build was successful
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Build successful!${NC}"
    echo ""
    echo -e "📊 Binary info:"
    ls -lh "$OUTPUT_NAME"
    echo ""
    echo -e "${BLUE}🚀 Usage:${NC}"
    echo "  ./$OUTPUT_NAME encrypt  # CLI encrypt"
    echo "  ./$OUTPUT_NAME decrypt  # CLI decrypt"
    echo "  ./$OUTPUT_NAME backup   # CLI backup"
else
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi