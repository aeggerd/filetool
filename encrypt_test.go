package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateEncryptedArchiveWithIndex(t *testing.T) {
	// Create a temporary test directory with files
	testDir := t.TempDir()
	outputDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"file1.txt":        "Hello World",
		"file2.txt":        "Test Content",
		"subdir/file3.txt": "Nested File",
	}

	for path, content := range files {
		fullPath := filepath.Join(testDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test encryption workflow - output to separate directory
	outputFile := filepath.Join(outputDir, "encrypted.enc")
	password := "test-password"
	key := sha256.Sum256([]byte(password))

	// Test the encryption process
	err := CreateEncryptedArchiveWithIndex(testDir, outputFile, key[:])
	if err != nil {
		t.Fatalf("CreateEncryptedArchiveWithIndex() error = %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Encrypted output file was not created")
	}

	// Verify it's a valid zip file
	reader, err := zip.OpenReader(outputFile)
	if err != nil {
		t.Fatalf("Output file is not a valid zip: %v", err)
	}
	defer reader.Close()

	// Check for index file and count encrypted files
	var indexFound bool
	var encryptedFiles []string
	for _, file := range reader.File {
		if file.Name == "index.json" {
			indexFound = true
		} else if strings.HasSuffix(file.Name, ".enc") {
			encryptedFiles = append(encryptedFiles, file.Name)
		}
	}

	if !indexFound {
		t.Error("index.json not found in encrypted archive")
	}

	if len(encryptedFiles) != 3 { // 3 test files
		t.Errorf("Expected 3 encrypted files, found %d: %v", len(encryptedFiles), encryptedFiles)
	}
}

func TestArchiveIndexStructure(t *testing.T) {
	testDir := t.TempDir()
	outputDir := t.TempDir()
	
	testFile := filepath.Join(testDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outputFile := filepath.Join(outputDir, "test.enc")
	password := "password"
	key := sha256.Sum256([]byte(password))
	
	err = CreateEncryptedArchiveWithIndex(testDir, outputFile, key[:])
	if err != nil {
		t.Fatalf("CreateEncryptedArchiveWithIndex() error = %v", err)
	}

	// Open and examine the index
	reader, err := zip.OpenReader(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	var indexFile *zip.File
	for _, file := range reader.File {
		if file.Name == "index.json" {
			indexFile = file
			break
		}
	}

	if indexFile == nil {
		t.Fatal("index.json not found")
	}

	// Read and parse index
	rc, err := indexFile.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()

	indexData, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}

	var index ArchiveIndex
	err = json.Unmarshal(indexData, &index)
	if err != nil {
		t.Fatalf("Failed to parse index.json: %v", err)
	}

	// Verify index structure
	if index.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", index.Version)
	}

	if len(index.Files) != 1 {
		t.Errorf("Expected 1 file in index, got %d", len(index.Files))
		return
	}

	file := index.Files[0]
	if file.RelativePath != "test.txt" {
		t.Errorf("Expected path test.txt, got %s", file.RelativePath)
	}
	if file.OriginalSize != 12 { // "test content" is 12 bytes
		t.Errorf("Expected size 12, got %d", file.OriginalSize)
	}
	if file.ModTime.IsZero() {
		t.Error("ModTime should not be zero")
	}
}

func TestEncryptEmptyFolder(t *testing.T) {
	testDir := t.TempDir()
	outputFile := filepath.Join(testDir, "empty.enc")
	password := "password"
	key := sha256.Sum256([]byte(password))

	err := CreateEncryptedArchiveWithIndex(testDir, outputFile, key[:])
	
	// Should handle empty folders gracefully - creates archive with just index
	if err != nil {
		t.Fatalf("CreateEncryptedArchiveWithIndex() should handle empty folders: %v", err)
	}

	// Verify the archive structure
	reader, err := zip.OpenReader(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	// Should have index.json
	var indexFound bool
	for _, file := range reader.File {
		if file.Name == "index.json" {
			indexFound = true
			break
		}
	}

	if !indexFound {
		t.Error("Empty archive should contain index.json")
	}
}

func TestEncryptNonExistentPath(t *testing.T) {
	nonExistent := "/path/that/does/not/exist"
	outputFile := "/tmp/test.enc"
	key := sha256.Sum256([]byte("password"))

	err := CreateEncryptedArchiveWithIndex(nonExistent, outputFile, key[:])
	if err == nil {
		t.Error("CreateEncryptedArchiveWithIndex() should fail for non-existent path")
	}
}

func TestEncryptInvalidOutputPath(t *testing.T) {
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Try to write to a directory that doesn't exist
	invalidOutput := "/invalid/path/that/does/not/exist/output.enc"
	key := sha256.Sum256([]byte("password"))

	err = CreateEncryptedArchiveWithIndex(testDir, invalidOutput, key[:])
	if err == nil {
		t.Error("CreateEncryptedArchiveWithIndex() should fail for invalid output path")
	}
}

func TestFileIndexEntry(t *testing.T) {
	// Test FileIndexEntry struct marshaling/unmarshaling
	entry := FileIndexEntry{
		RelativePath: "test/file.txt",
		OriginalSize: 1234,
		ZipName:      "test/file.txt.enc",
	}

	// Test JSON marshaling
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Failed to marshal FileIndexEntry: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled FileIndexEntry
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal FileIndexEntry: %v", err)
	}

	// Verify fields
	if unmarshaled.RelativePath != entry.RelativePath {
		t.Errorf("RelativePath mismatch: got %s, want %s", unmarshaled.RelativePath, entry.RelativePath)
	}
	if unmarshaled.OriginalSize != entry.OriginalSize {
		t.Errorf("OriginalSize mismatch: got %d, want %d", unmarshaled.OriginalSize, entry.OriginalSize)
	}
	if unmarshaled.ZipName != entry.ZipName {
		t.Errorf("ZipName mismatch: got %s, want %s", unmarshaled.ZipName, entry.ZipName)
	}
}