package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "filetool",
	Short: "File backup and encryption tool",
}


func main() {
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.Execute()
}

func HandleCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nCleaning up decrypted files...")
		CleanUpDecryptedFiles()
		os.Exit(0)
	}()
}
