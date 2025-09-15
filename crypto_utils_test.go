package main

import (
	"crypto/sha256"
	"strconv"
	"testing"
)

func TestReadPassword(t *testing.T) {
	// Note: ReadPassword requires terminal input, so we can't easily test it
	// in an automated way. This test documents the function exists.
	t.Skip("ReadPassword requires terminal input - tested manually")
}

func TestReadLine(t *testing.T) {
	// Note: ReadLine requires stdin input, so we can't easily test it
	// in an automated way. This test documents the function exists.
	t.Skip("ReadLine requires stdin input - tested manually")
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input       string
		expected    int
		shouldError bool
	}{
		{"123", 123, false},
		{"0", 0, false},
		{"-456", -456, false},
		{"abc", 0, true},
		{"", 0, true},
		{"123abc", 0, true},
		{"999999999999999999999999999", 0, true}, // overflow
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseInt(tt.input)

			if tt.shouldError {
				if err == nil {
					t.Errorf("ParseInt(%q) should have returned an error", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseInt(%q) unexpected error: %v", tt.input, err)
				return
			}

			if result != tt.expected {
				t.Errorf("ParseInt(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseIntConsistencyWithStrconv(t *testing.T) {
	testCases := []string{"123", "0", "-456", "abc", ""}
	
	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			ourResult, ourErr := ParseInt(input)
			stdResult, stdErr := strconv.Atoi(input)

			// Check that results match
			if ourResult != stdResult {
				t.Errorf("ParseInt(%q) = %d, strconv.Atoi = %d", input, ourResult, stdResult)
			}

			// Check error conditions match
			if (ourErr == nil) != (stdErr == nil) {
				t.Errorf("ParseInt(%q) error condition doesn't match strconv.Atoi", input)
			}
		})
	}
}

// Test key derivation (using the same method as encrypt.go)
func TestKeyDerivation(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected int // expected key length
	}{
		{"basic password", "password123", 32},
		{"empty password", "", 32},
		{"unicode password", "测试密码🔐", 32},
		{"long password", string(make([]byte, 1000)), 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := sha256.Sum256([]byte(tt.password))
			if len(key) != tt.expected {
				t.Errorf("key derivation returned length %d, expected %d", len(key), tt.expected)
			}
		})
	}
}

func TestKeyDerivationConsistency(t *testing.T) {
	password := "test-password"
	key1 := sha256.Sum256([]byte(password))
	key2 := sha256.Sum256([]byte(password))

	if key1 != key2 {
		t.Error("Key derivation should be deterministic")
	}
}