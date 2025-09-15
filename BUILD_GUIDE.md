# 🔧 Filetool Build Guide

## 📦 Build Type

### ⚡ CLI-Only Build
**Features**: Command line interface only, smallest binary size
```bash
./build.sh                     # Default and only option
BUILD_TYPE=minimal ./build.sh  # Explicit minimal build
```

## 🎯 Cross-Platform Builds

### Windows
```bash
GOOS=windows ./build.sh
```

### macOS
```bash
GOOS=darwin ./build.sh
```

### Linux ARM64
```bash
GOARCH=arm64 ./build.sh
```

## 🚀 Usage Examples

### Command Line
```bash
./filetool encrypt /path/to/folder output.enc
./filetool decrypt output.enc
./filetool backup /source /destination
```

## 📊 Build Size

| Build Type | Size  | Features |
|-----------|-------|----------|
| CLI       | ~8MB  | Command line only |

## 🛠️ Development

### Prerequisites
```bash
# Go 1.20+
go version
```

### Build All Platforms
```bash
./build-releases.sh
```
Generates CLI builds for:
- Linux (amd64, arm64)
- Windows (amd64) 
- macOS (amd64, arm64)

## 🐛 Troubleshooting

### Build Fails
```bash
# Clean dependencies
go mod tidy

# Try manual build
go build -o filetool
```

## 📋 Quick Reference

```bash
# Standard build
./build.sh                              # CLI build (default)

# Cross-platform
GOOS=windows ./build.sh                 # Windows
GOOS=darwin ./build.sh                  # macOS

# Custom output
OUTPUT_NAME=my-filetool ./build.sh      # Custom name
```