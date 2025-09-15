# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based CLI tool called "filetool" that provides secure file backup, encryption, and decryption capabilities. The tool encrypts individual files using AES-256 encryption and packages them into zip archives with metadata indexing for fast access.

## Build & Development Commands

```bash
# Build the project (CLI-only)
./build.sh

# Or build manually
go build -o filetool

# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o filetool

# Run the tool
./filetool --help

# Test basic functionality
./filetool backup <source> <dest>
./filetool encrypt <source-folder> <output-file>
./filetool decrypt <encrypted-file>
```

## Architecture Overview

### Core Components

- **main.go**: Entry point with Cobra CLI setup, defines root command and subcommands
- **encrypt.go**: File encryption logic with AES-256-CFB, creates indexed zip archives
- **decrypt.go**: Interactive decryption with file selection and automatic cleanup
- **backup.go**: Simple file backup with progress tracking and incremental sync
- **crypto_utils.go**: Shared cryptographic utilities and password handling
- **cleanup.go**: Temporary file management and cleanup on exit
- **completion.go**: Shell completion support

### Key Features

1. **Indexed Archives**: Creates `index.json` in archives for fast file listing without decryption
2. **Interactive Decryption**: Multi-select UI for choosing which files to decrypt
3. **Progress Tracking**: Real-time progress bars during encryption/decryption operations
4. **Automatic Cleanup**: Tracks and removes decrypted files on Ctrl+C
5. **Incremental Backup**: Only copies modified files during backup operations

### Data Structures

- **FileIndexEntry**: Metadata for encrypted files (path, size, modification time)
- **ArchiveIndex**: Container for all file metadata with versioning
- **ProgressWriter/ProgressReader**: Wrapper structs for progress tracking during I/O

### Encryption Details

- Uses AES-256 in CFB mode
- Password-derived keys using SHA-256
- Random IV per file stored at beginning of encrypted data
- No compression in zip archives (files already encrypted)

### Dependencies

- `github.com/spf13/cobra`: CLI framework
- `github.com/AlecAivazis/survey/v2`: Interactive prompts
- `github.com/schollz/progressbar/v3`: Progress tracking
- `golang.org/x/term`: Secure password input

## File Organization

```
├── main.go              # CLI entry point and command structure
├── encrypt.go           # Encryption logic and archive creation
├── decrypt.go           # Decryption with interactive file selection
├── backup.go            # File backup functionality
├── crypto_utils.go      # Shared cryptographic utilities
├── cleanup.go           # Temporary file management
├── completion.go        # Shell completion
├── build.sh            # CLI-only build script
└── export.sh           # Development utility for code export
```

## Development Notes

- The tool is designed for security-focused file operations
- All encrypted files use individual IVs for security
- Archive format includes unencrypted index for performance
- Temporary decrypted files are automatically tracked and cleaned up
- Progress bars provide user feedback during long operations
- The codebase follows defensive security practices with proper error handling