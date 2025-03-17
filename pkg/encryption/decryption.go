package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"regexp"
)

// DecryptData decrypts data that was encrypted with EncryptData
func DecryptData(encryptedData []byte, signature []byte) ([]byte, error) {
	// Derive the same deterministic key from the signature
	key := deriveKeyFromSignature(signature)

	// Use the common decryption function with the derived key
	return DecryptDataWithKey(encryptedData, key)
}

// DecryptDataWithKey decrypts data using a provided key directly
// This allows decryption without needing the original signature
func DecryptDataWithKey(encryptedData []byte, key []byte) ([]byte, error) {
	// Check if the data is actually base64 encoded
	var decoded []byte
	var err error

	// Use a simple check to determine if the data is already base64 encoded
	isBase64 := isBase64Encoded(encryptedData)

	if isBase64 {
		// Decode base64 data
		decoded = make([]byte, base64.StdEncoding.DecodedLen(len(encryptedData)))
		n, err := base64.StdEncoding.Decode(decoded, encryptedData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 data: %w", err)
		}
		decoded = decoded[:n] // Trim to actual size
	} else {
		// Data is already decoded, use as is
		decoded = encryptedData
	}

	// Extract IV from the beginning
	iv := decoded[:aes.BlockSize]
	actualCiphertext := decoded[aes.BlockSize:]

	// Create cipher and decrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create cipher stream
	stream := cipher.NewCTR(block, iv)

	// Allocate plaintext buffer
	plaintext := make([]byte, len(actualCiphertext))

	// Decrypt data
	stream.XORKeyStream(plaintext, actualCiphertext)

	return plaintext, nil
}

// isBase64Encoded checks if data is likely to be base64 encoded
func isBase64Encoded(data []byte) bool {
	// Check if data contains only valid base64 characters
	re := regexp.MustCompile("^[A-Za-z0-9+/]*=*$")
	return re.Match(data)
}
