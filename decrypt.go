package main

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt <encrypted-file>",
	Short: "Decrypt an archive and interactively extract files",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		HandleCleanup()
		RunDecryptInteractive(args[0])
		if len(decryptedFiles) > 0 {
			fmt.Println("\nTracking decrypted files:")
			for _, file := range decryptedFiles {
				fmt.Println(" -", file)
			}
			fmt.Println("\nPress Ctrl+C to cleanup and exit")
			select {} // Block forever
		}
	},
}

func RunDecryptInteractive(encFile string) {
	fmt.Print("Enter decryption password: ")
	password := ReadPassword()

	// Derive key once for efficiency
	key := sha256.Sum256([]byte(password))

	// Fast index-based file listing
	fmt.Println("Reading archive index...")
	files, err := ListEncryptedFilesFromIndex(encFile)
	if err != nil {
		fmt.Printf("Failed to read index (trying fallback): %v\n", err)
		// Fallback to old method if index is missing/corrupted
		files, err = ListEncryptedFilesFallback(encFile)
		if err != nil {
			fmt.Printf("Failed to open archive: %v\n", err)
			return
		}
	}

	if len(files) == 0 {
		fmt.Println("No encrypted files found in archive")
		return
	}

	fmt.Printf("Found %d files in archive\n", len(files))

	// Create display options for the user with full paths and file sizes
	displayOptions := make([]string, len(files))
	for i, file := range files {
		sizeStr := formatFileSize(file.OriginalSize)
		// Show full relative path to distinguish between files with same names
		displayOptions[i] = fmt.Sprintf("%s (%s)", file.RelativePath, sizeStr)
	}

	// Let user select files
	selected := MultiSelectPrompt("Select files to decrypt:", displayOptions)
	if len(selected) == 0 {
		fmt.Println("No files selected")
		return
	}

	// Create a map for quick lookup (strip size info)
	selectedMap := make(map[string]bool)
	for _, sel := range selected {
		// Extract just the path part (before the size info in parentheses)
		if idx := strings.LastIndex(sel, " ("); idx != -1 {
			sel = sel[:idx]
		}
		selectedMap[sel] = true
	}

	// Open archive once and decrypt selected files efficiently
	if err := DecryptSelectedFilesFromArchive(encFile, files, selectedMap, key[:]); err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
	}
}

// Fast index-based file listing
func ListEncryptedFilesFromIndex(archivePath string) ([]FileIndexEntry, error) {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Look for index file
	var indexFile *zip.File
	for _, f := range reader.File {
		if f.Name == "index.json" {
			indexFile = f
			break
		}
	}

	if indexFile == nil {
		return nil, errors.New("no index file found")
	}

	// Read index
	indexReader, err := indexFile.Open()
	if err != nil {
		return nil, err
	}
	defer indexReader.Close()

	var index ArchiveIndex
	if err := json.NewDecoder(indexReader).Decode(&index); err != nil {
		return nil, err
	}

	return index.Files, nil
}

// Fallback method for archives without index
func ListEncryptedFilesFallback(archivePath string) ([]FileIndexEntry, error) {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var files []FileIndexEntry
	for _, f := range reader.File {
		if strings.HasSuffix(f.Name, ".enc") {
			// Remove .enc extension for display
			originalName := strings.TrimSuffix(f.Name, ".enc")
			files = append(files, FileIndexEntry{
				ZipName:      f.Name,
				RelativePath: originalName,
				OriginalSize: int64(f.UncompressedSize64), // Approximate
			})
		}
	}

	return files, nil
}

// DecryptSelectedFilesFromArchive efficiently decrypts multiple files from an archive
func DecryptSelectedFilesFromArchive(archivePath string, files []FileIndexEntry, selectedMap map[string]bool, key []byte) error {
	// Open archive once
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Calculate total size for progress
	var totalSize int64
	var selectedFiles []FileIndexEntry
	for _, file := range files {
		if selectedMap[file.RelativePath] {
			selectedFiles = append(selectedFiles, file)
			totalSize += file.OriginalSize
		}
	}

	// Create single progress bar for all files
	bar := progressbar.DefaultBytes(totalSize, "decrypting files")

	// Decrypt each selected file
	for _, file := range selectedFiles {
		outputPath := file.RelativePath
		if err := decryptSingleFileFromReader(reader, file, outputPath, key, bar); err != nil {
			fmt.Printf("Failed to decrypt %s: %v\n", file.RelativePath, err)
		} else {
			absPath, _ := filepath.Abs(outputPath)
			RegisterTempFile(absPath)
			fmt.Printf("Decrypted: %s\n", outputPath)
		}
	}

	return nil
}

func DecryptFileFromArchive(archivePath string, fileInfo FileIndexEntry, outputPath, password string) error {
	// Open archive
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Find the file in archive
	var targetFile *zip.File
	for _, f := range reader.File {
		if f.Name == fileInfo.ZipName {
			targetFile = f
			break
		}
	}

	if targetFile == nil {
		return errors.New("file not found in archive")
	}

	// Open the encrypted file from archive
	encryptedReader, err := targetFile.Open()
	if err != nil {
		return err
	}
	defer encryptedReader.Close()

	// Create output directory if needed
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Create progress bar using original file size if available
	fileSize := fileInfo.OriginalSize
	if fileSize == 0 {
		fileSize = int64(targetFile.UncompressedSize64)
	}

	bar := progressbar.DefaultBytes(
		fileSize,
		fmt.Sprintf("decrypting %s", fileInfo.RelativePath),
	)

	// Decrypt the file
	return decryptFromReader(encryptedReader, outputFile, password, bar)
}

// decryptSingleFileFromReader decrypts a single file using an already-opened archive reader
func decryptSingleFileFromReader(reader *zip.ReadCloser, fileInfo FileIndexEntry, outputPath string, key []byte, progressBar *progressbar.ProgressBar) error {
	// Find the file in archive
	var targetFile *zip.File
	for _, f := range reader.File {
		if f.Name == fileInfo.ZipName {
			targetFile = f
			break
		}
	}

	if targetFile == nil {
		return errors.New("file not found in archive")
	}

	// Open the encrypted file from archive
	encryptedReader, err := targetFile.Open()
	if err != nil {
		return err
	}
	defer encryptedReader.Close()

	// Create output directory if needed
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Decrypt the file using pre-derived key
	return decryptFromReaderWithKey(encryptedReader, outputFile, key, progressBar)
}

func decryptFromReader(reader io.Reader, writer io.Writer, password string, progressBar *progressbar.ProgressBar) error {
	// Generate key from password
	key := sha256.Sum256([]byte(password))
	return decryptFromReaderWithKey(reader, writer, key[:], progressBar)
}

// decryptFromReaderWithKey decrypts using a pre-derived key for efficiency
func decryptFromReaderWithKey(reader io.Reader, writer io.Writer, key []byte, progressBar *progressbar.ProgressBar) error {
	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Read IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(reader, iv); err != nil {
		return errors.New("failed to read IV: likely corrupted file or wrong password")
	}

	// Create cipher stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Create progress tracking reader
	progressReader := &ProgressReader{
		Reader: reader,
		Bar:    progressBar,
	}

	// Decrypt
	streamReader := &cipher.StreamReader{S: stream, R: progressReader}
	_, err = io.Copy(writer, streamReader)

	if err != nil {
		return errors.New("decryption failed: likely wrong password")
	}

	return nil
}

// Helper struct for progress tracking during decryption
type ProgressReader struct {
	Reader io.Reader
	Bar    *progressbar.ProgressBar
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	if n > 0 {
		pr.Bar.Add(n)
	}
	return
}

func MultiSelectPrompt(message string, items []string) []string {
	var selected []string
	prompt := &survey.MultiSelect{
		Message: message,
		Options: items,
	}
	survey.AskOne(prompt, &selected)
	return selected
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
