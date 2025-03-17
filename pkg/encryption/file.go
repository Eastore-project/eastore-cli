package encryption

import (
	"fmt"
	"os"

	"github.com/eastore-project/eastore/pkg/utils"
)

// EncryptFile encrypts a file using the private key to derive the encryption key
// Returns encrypted data bytes, hex-encoded key string, and error if any
func EncryptFile(inputPath string, privateKey string) ([]byte, string, error) {
	// Calculate file CID for encryption
	fileCID, err := utils.CalculateFileCID(inputPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to calculate file CID: %w", err)
	}
	cidStr := fileCID.String()

	// Sign the message to derive encryption key
	signature, err := utils.SignMessage(privateKey, cidStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to sign message for encryption: %w", err)
	}

	// Read the file
	fileData, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read input file: %w", err)
	}

	// Encrypt the data
	encryptedData, hexKey, err := EncryptData(fileData, signature)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	return encryptedData, hexKey, nil
}
