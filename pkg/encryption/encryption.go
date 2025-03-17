// Package encryption provides functionality for encrypting and decrypting files
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// EncryptData encrypts data using AES with a key derived from signature
// Returns base64-encoded encrypted data and the hex encoded key for logging
func EncryptData(data []byte, signature []byte) ([]byte, string, error) {
	// Derive a deterministic 32-byte key from the signature using SHA-256
	key := deriveKeyFromSignature(signature)
	keyHex := hex.EncodeToString(key)

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create deterministic IV by hashing the key
	iv := deriveIVFromKey(key)

	// Create cipher stream
	stream := cipher.NewCTR(block, iv)

	// Allocate ciphertext buffer
	ciphertext := make([]byte, len(data))

	// Encrypt data
	stream.XORKeyStream(ciphertext, data)

	// Return base64 encoded for string compatibility
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(encoded, ciphertext)
	return encoded, keyHex, nil
}

// deriveKeyFromSignature creates a deterministic 32-byte key from the signature
// using SHA-256, which always produces a 32-byte output (perfect for AES-256)
func deriveKeyFromSignature(signature []byte) []byte {
	hasher := sha256.New()
	hasher.Write(signature)
	return hasher.Sum(nil) // Returns a 32-byte array
}

// deriveIVFromKey creates a deterministic IV from the key
// We need a 16-byte IV for AES, so we'll hash the key again with a prefix
func deriveIVFromKey(key []byte) []byte {
	hasher := sha256.New()
	hasher.Write([]byte("IV"))
	hasher.Write(key)
	return hasher.Sum(nil)[:aes.BlockSize] // First 16 bytes (AES block size)
}
