package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func ReadPassword() string {
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(bytePassword)
}

func ReadLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func ParseInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	return n, err
}
