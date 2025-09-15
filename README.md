# 🔐 Filetool - Secure File Encryption & Backup CLI

[![Build and Release](https://github.com/aeggerd/filetool/actions/workflows/release.yml/badge.svg)](https://github.com/aeggerd/filetool/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aeggerd/filetool)](https://goreportcard.com/report/github.com/aeggerd/filetool)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A secure, fast, and user-friendly command-line tool for file encryption, decryption, and backup operations. Built with Go, featuring AES-256 encryption with individual file IVs, interactive decryption, and cross-platform support.

## ✨ Features

- 🔒 **AES-256 Encryption**: Each file encrypted with unique IV for maximum security
- 📁 **Indexed Archives**: Fast file listing without decryption via JSON index
- 🎯 **Interactive Decryption**: Multi-select UI for choosing which files to decrypt  
- 📊 **Progress Tracking**: Real-time progress bars for all operations
- 🧹 **Auto Cleanup**: Temporary files automatically removed on exit
- 🔄 **Incremental Backup**: Only copies modified files for efficiency
- 🌐 **Cross-Platform**: Linux, Windows, macOS (Intel & Apple Silicon)
- ⚡ **CLI-Only**: No GUI dependencies, perfect for servers and automation

## 🚀 Quick Start

### Download

Get the latest release for your platform:

```bash
# Linux x86_64
wget https://github.com/aeggerd/filetool/releases/latest/download/filetool-linux-amd64
chmod +x filetool-linux-amd64
sudo mv filetool-linux-amd64 /usr/local/bin/filetool

# macOS Intel
wget https://github.com/aeggerd/filetool/releases/latest/download/filetool-macos-amd64
chmod +x filetool-macos-amd64
sudo mv filetool-macos-amd64 /usr/local/bin/filetool

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/aeggerd/filetool/releases/latest/download/filetool-windows-amd64.exe" -OutFile "filetool.exe"
```

### Build from Source

```bash
git clone git@github.com:aeggerd/filetool.git
cd filetool
./build.sh
```

## 📖 Usage

### Encrypt Files

Encrypt a folder and all its contents:

```bash
filetool encrypt /path/to/folder output.enc
```

The tool will:
- Prompt for a password (hidden input)
- Encrypt each file individually with AES-256-CFB
- Create an indexed archive with metadata
- Show real-time progress

### Decrypt Files

Decrypt and extract files interactively:

```bash
filetool decrypt output.enc
```

Features:
- Lists all files in the archive instantly (no decryption needed)
- Interactive multi-select interface
- Choose destination directory
- Automatic cleanup of temporary files

### Backup Files

Incremental file backup with progress tracking:

```bash
filetool backup /source/folder /backup/destination
```

- Only copies newer or modified files
- Preserves directory structure
- Real-time progress display
- Removes files from destination that no longer exist in source

## 🔧 Installation

### Package Managers

#### Homebrew (macOS/Linux)
```bash
# Coming soon
brew install filetool
```

#### Chocolatey (Windows)
```bash
# Coming soon  
choco install filetool
```

### Manual Installation

1. Download the appropriate binary from [Releases](https://github.com/aeggerd/filetool/releases)
2. Verify integrity with SHA256SUMS
3. Make executable and move to PATH:

```bash
# Verify download (Linux/macOS)
sha256sum -c SHA256SUMS

# Install
chmod +x filetool-*
sudo mv filetool-* /usr/local/bin/filetool
```

## 🛡️ Security

### Encryption Details

- **Algorithm**: AES-256 in CFB (Cipher Feedback) mode
- **Key Derivation**: SHA-256 hash of password
- **IV Generation**: Cryptographically secure random IV per file
- **Archive Format**: ZIP with no compression (files already encrypted)
- **Index**: Unencrypted JSON metadata for fast file listing

### Security Best Practices

- Each file gets a unique initialization vector (IV)
- No password storage or logging
- Temporary decrypted files auto-cleaned on exit
- Memory is not explicitly cleared (Go garbage collector handles this)
- No network communication

### Threat Model

**Protects Against**:
- Unauthorized file access
- Data theft from backups
- Casual inspection of sensitive data

**Does NOT Protect Against**:
- Keyloggers capturing passwords
- Memory dumps while decrypted files exist
- Sophisticated state-level attacks
- Weak passwords (use strong, unique passwords)

## 🏗️ Development

### Prerequisites

- Go 1.21+
- Git

### Building

```bash
# Clone repository
git clone git@github.com:aeggerd/filetool.git
cd filetool

# Build for current platform
./build.sh

# Build for all platforms
./build-releases.sh
```

### Testing

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -cover ./...

# Run specific test
go test -v -run TestEncrypt
```

### Project Structure

```
├── main.go              # CLI entry point and commands
├── encrypt.go           # Encryption logic and archive creation
├── decrypt.go           # Interactive decryption
├── backup.go            # File backup with progress tracking
├── crypto_utils.go      # Shared utilities (password input, etc.)
├── cleanup.go           # Temporary file management
├── completion.go        # Shell completion support
├── *_test.go           # Comprehensive unit tests
├── build.sh            # Build script
└── .github/workflows/  # CI/CD pipeline
```

## 📊 Performance

### Benchmarks

- **Encryption Speed**: ~100MB/s (depends on storage)
- **Archive Creation**: Optimized for SSDs with no compression
- **Memory Usage**: Minimal, streams large files
- **Binary Size**: ~8MB (static binary, no dependencies)

### File Size Limits

- **Maximum file size**: Limited by available disk space
- **Maximum files per archive**: No artificial limit
- **Archive size**: Limited by ZIP64 format (~9 exabytes)

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for new functionality
4. Ensure all tests pass (`go test ./...`)
5. Commit changes (`git commit -m 'Add amazing feature'`)
6. Push to branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Style

- Follow Go conventions (`gofmt`, `go vet`)
- Write comprehensive tests
- Document public functions
- Keep functions focused and small

## 📝 Changelog

### v1.0.0 (Latest)
- 🎉 Initial release
- ✅ AES-256 file encryption
- ✅ Interactive decryption
- ✅ Incremental backup
- ✅ Cross-platform support
- ✅ Comprehensive test coverage
- ✅ CI/CD pipeline

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Survey](https://github.com/AlecAivazis/survey) - Interactive prompts
- [ProgressBar](https://github.com/schollz/progressbar) - Progress indication
- [Go Team](https://golang.org/) - Amazing programming language

## 🐛 Support

- 📖 **Documentation**: Check this README and `--help` commands
- 🐞 **Bug Reports**: [Open an issue](https://github.com/aeggerd/filetool/issues)
- 💡 **Feature Requests**: [Open an issue](https://github.com/aeggerd/filetool/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/aeggerd/filetool/discussions)

---

<p align="center">
  <strong>⭐ If you find this project helpful, please give it a star! ⭐</strong>
</p>