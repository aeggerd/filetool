package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBackupFiles(t *testing.T) {
	// Create temporary source directory
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create test files in source
	testFiles := map[string]string{
		"file1.txt":        "Content 1",
		"file2.txt":        "Content 2", 
		"subdir/file3.txt": "Content 3",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(srcDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test backup operation
	RunBackup(srcDir, dstDir)

	// Verify all files were backed up
	for path, expectedContent := range testFiles {
		dstPath := filepath.Join(dstDir, path)
		
		// Check file exists
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			t.Errorf("Backup file does not exist: %s", path)
			continue
		}

		// Check content matches
		content, err := os.ReadFile(dstPath)
		if err != nil {
			t.Errorf("Failed to read backup file %s: %v", path, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("Backup file %s content mismatch: got %q, want %q", path, string(content), expectedContent)
		}
	}
}

func TestBackupIncrementalSync(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create initial files
	file1 := filepath.Join(srcDir, "file1.txt")
	file2 := filepath.Join(srcDir, "file2.txt")
	
	err := os.WriteFile(file1, []byte("original content"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(file2, []byte("another file"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// First backup
	RunBackup(srcDir, dstDir)

	// Get modification time of backed up file
	dstFile1 := filepath.Join(dstDir, "file1.txt") 
	stat1, err := os.Stat(dstFile1)
	if err != nil {
		t.Fatal(err)
	}
	firstBackupTime := stat1.ModTime()

	// Wait a bit to ensure different modification times
	time.Sleep(10 * time.Millisecond)

	// Modify source file
	err = os.WriteFile(file1, []byte("modified content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Second backup (should update only the modified file)
	RunBackup(srcDir, dstDir)

	// Check that file1 was updated
	content, err := os.ReadFile(dstFile1)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "modified content" {
		t.Errorf("File1 should have been updated, got: %s", string(content))
	}

	// Check that modification time changed (indicating it was copied)
	stat2, err := os.Stat(dstFile1)
	if err != nil {
		t.Fatal(err)
	}
	secondBackupTime := stat2.ModTime()

	if !secondBackupTime.After(firstBackupTime) {
		t.Error("File should have been updated with newer modification time")
	}
}

func TestBackupNonExistentSource(t *testing.T) {
	// RunBackup doesn't return errors - it would panic or handle internally
	// This test is more about documenting expected behavior
	t.Skip("RunBackup doesn't return errors to test")
}

func TestBackupToInvalidDestination(t *testing.T) {
	// RunBackup doesn't return errors to test
	t.Skip("RunBackup doesn't return errors to test")
}

func TestBackupEmptyDirectory(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Backup empty directory
	RunBackup(srcDir, dstDir)

	// Destination should exist
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		t.Error("Destination directory should exist after backup")
	}
}

func TestBackupPreserveDirectoryStructure(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create nested directory structure
	nestedPath := filepath.Join(srcDir, "level1", "level2", "level3")
	err := os.MkdirAll(nestedPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create file in nested directory
	testFile := filepath.Join(nestedPath, "deep_file.txt")
	err = os.WriteFile(testFile, []byte("deep content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Backup
	RunBackup(srcDir, dstDir)

	// Verify nested structure was preserved
	dstFile := filepath.Join(dstDir, "level1", "level2", "level3", "deep_file.txt")
	content, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read backed up nested file: %v", err)
	}

	if string(content) != "deep content" {
		t.Errorf("Nested file content mismatch: got %q, want %q", string(content), "deep content")
	}
}