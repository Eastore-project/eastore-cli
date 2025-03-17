package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eastore-project/eastore/pkg/encryption"
	"github.com/urfave/cli/v2"
)

// EncryptCommand returns the CLI command for encrypting data
func EncryptCommand() *cli.Command {
	return &cli.Command{
		Name:  "encrypt",
		Usage: "Encrypt a file using AES with a key derived from wallet signature",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Required: true,
				Usage:    "Input file or folder path",
				EnvVars:  []string{"INPUT_PATH"},
			},
			&cli.StringFlag{
				Name:    "out-dir",
				Usage:   "Output directory for encrypted files",
				Value:   "./encrypted",
				EnvVars: []string{"OUT_DIR"},
			},
		},
		Action: encryptAction,
	}
}

func encryptAction(cCtx *cli.Context) error {
	inputPath := cCtx.String("input")
	outDir := cCtx.String("out-dir")
	privateKey := cCtx.String("private-key")

	// Use the EncryptFile function
	encryptedData, hexKey, err := encryption.EncryptFile(inputPath, privateKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt file: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write the encrypted file
	encryptedFilePath := filepath.Join(outDir, "encrypted_"+filepath.Base(inputPath))
	if err := os.WriteFile(encryptedFilePath, encryptedData, 0644); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	fmt.Printf("File encrypted successfully\n")
	fmt.Printf("Original file: %s\n", inputPath)
	fmt.Printf("Derived key: %s\n", hexKey)
	fmt.Printf("Encrypted file: %s\n", encryptedFilePath)

	return nil
}
