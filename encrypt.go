package main

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt <source-folder> <output-file>",
	Short: "Encrypt files individually and pack into an encrypted archive",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		RunEncrypt(args[0], args[1])
	},
}

// FileIndexEntry represents metadata for each encrypted file
type FileIndexEntry struct {
	RelativePath string    `json:"relative_path"`
	OriginalSize int64     `json:"original_size"`
	ModTime      time.Time `json:"mod_time"`
	ZipName      string    `json:"zip_name"`
}

// ArchiveIndex contains metadata about all files in the archive
type ArchiveIndex struct {
	Files     []FileIndexEntry `json:"files"`
	CreatedAt time.Time        `json:"created_at"`
	Version   string           `json:"version"`
}

func RunEncrypt(src, dst string) {
	fmt.Print("Enter encryption password: ")
	password := ReadPassword()

	// Derive key once for efficiency
	key := sha256.Sum256([]byte(password))

	// Create the encrypted archive with index
	err := CreateEncryptedArchiveWithIndex(src, dst, key[:])
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}

	fmt.Println("Encrypted archive saved to:", dst)
}

func CreateEncryptedArchiveWithIndex(srcDir, archivePath string, key []byte) error {
	// Create the archive file
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	// Create zip writer with no compression for faster access
	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	// Collect file information first
	var fileEntries []FileIndexEntry
	var totalSize int64

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Normalize path separators
		relPath = filepath.ToSlash(relPath)
		zipName := relPath + ".enc"

		fileEntries = append(fileEntries, FileIndexEntry{
			RelativePath: relPath,
			OriginalSize: info.Size(),
			ModTime:      info.ModTime(),
			ZipName:      zipName,
		})

		totalSize += info.Size()
		return nil
	})
	if err != nil {
		return err
	}

	// Create and add index file first (unencrypted for fast access)
	index := ArchiveIndex{
		Files:     fileEntries,
		CreatedAt: time.Now(),
		Version:   "1.0",
	}

	if err := addIndexToArchive(zipWriter, index); err != nil {
		return fmt.Errorf("failed to create index: %v", err)
	}

	// Create progress bar
	bar := progressbar.DefaultBytes(totalSize, "encrypting files")

	// Encrypt all files using pre-derived key
	for _, entry := range fileEntries {
		sourcePath := filepath.Join(srcDir, filepath.FromSlash(entry.RelativePath))
		if err := addEncryptedFileToArchive(zipWriter, sourcePath, entry.ZipName, key, bar); err != nil {
			return fmt.Errorf("failed to encrypt %s: %v", entry.RelativePath, err)
		}
	}

	return nil
}

func addIndexToArchive(zipWriter *zip.Writer, index ArchiveIndex) error {
	// Create index entry with no compression
	indexHeader := &zip.FileHeader{
		Name:   "index.json",
		Method: zip.Store, // No compression for faster access
	}
	indexHeader.SetModTime(time.Now())

	indexWriter, err := zipWriter.CreateHeader(indexHeader)
	if err != nil {
		return err
	}

	// Write index as JSON
	encoder := json.NewEncoder(indexWriter)
	encoder.SetIndent("", "  ")
	return encoder.Encode(index)
}

func addEncryptedFileToArchive(zipWriter *zip.Writer, sourcePath, zipName string, key []byte, progressBar *progressbar.ProgressBar) error {
	// Create zip entry with no compression for faster access
	fileHeader := &zip.FileHeader{
		Name:   zipName,
		Method: zip.Store, // No compression - we're already encrypting
	}
	fileHeader.SetModTime(time.Now())

	zipEntry, err := zipWriter.CreateHeader(fileHeader)
	if err != nil {
		return err
	}

	// Encrypt file directly into zip entry
	return encryptFileToWriter(sourcePath, zipEntry, key, progressBar)
}

func encryptFileToWriter(inputPath string, writer io.Writer, key []byte, progressBar *progressbar.ProgressBar) error {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Create cipher using pre-derived key
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Generate random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	// Write IV to output first
	if _, err := writer.Write(iv); err != nil {
		return err
	}

	// Create cipher stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Create progress tracking writer
	progressWriter := &ProgressWriter{
		Writer: writer,
		Bar:    progressBar,
	}

	// Encrypt and write
	streamWriter := &cipher.StreamWriter{S: stream, W: progressWriter}
	_, err = io.Copy(streamWriter, inputFile)

	return err
}

// Helper struct for progress tracking during encryption
type ProgressWriter struct {
	Writer io.Writer
	Bar    *progressbar.ProgressBar
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.Writer.Write(p)
	if n > 0 {
		pw.Bar.Add(n)
	}
	return
}
