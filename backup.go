package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup <source> <destination>",
	Short: "Backup files from source to destination",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		RunBackup(args[0], args[1])
	},
}

func RunBackup(src, dst string) {
	// Count total files first
	var totalFiles int
	filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() {
			totalFiles++
		}
		return nil
	})

	bar := progressbar.Default(int64(totalFiles), "backing up")

	filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(src, path)
		targetPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		srcInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		dstInfo, err := os.Stat(targetPath)
		if err == nil {
			if !srcInfo.ModTime().After(dstInfo.ModTime()) && srcInfo.Size() == dstInfo.Size() {
				_ = bar.Add(1)
				return nil
			}
		}

		srcFile, _ := os.Open(path)
		defer srcFile.Close()
		dstFile, _ := os.Create(targetPath)
		defer dstFile.Close()
		_, err = io.Copy(dstFile, srcFile)

		_ = bar.Add(1)
		return err
	})

	// Cleanup: delete files not in source
	filepath.WalkDir(dst, func(path string, d os.DirEntry, err error) error {
		relPath, _ := filepath.Rel(dst, path)
		srcPath := filepath.Join(src, relPath)
		_, err = os.Stat(srcPath)
		if os.IsNotExist(err) {
			return os.RemoveAll(path)
		}
		return nil
	})
}
