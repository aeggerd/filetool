package main

import (
	"fmt"
	"os"
)

var decryptedFiles []string

func RegisterTempFile(path string) {
	decryptedFiles = append(decryptedFiles, path)
}

func CleanUpDecryptedFiles() {
	for _, file := range decryptedFiles {
		fmt.Println("Removing:", file)
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Failed to remove:", file, err)
		}
	}
}
