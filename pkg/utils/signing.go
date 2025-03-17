package utils

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// SignMessage signs a message with the provided private key and returns the signature
func SignMessage(privateKey, message string) ([]byte, error) {
	// Parse private key
	privateKeyECDSA, err := crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create Ethereum prefixed message
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)

	// Hash the prefixed message
	msgHash := crypto.Keccak256Hash([]byte(prefixedMessage))

	// Sign the hash with the private key
	signature, err := crypto.Sign(msgHash.Bytes(), privateKeyECDSA)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	// Adjust the 'v' value to Ethereum wallet standard (add 27)
	signature[64] += 27

	return signature, nil
}
