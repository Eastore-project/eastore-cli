package commands

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eastore-project/eastore/pkg/encryption"
	"github.com/urfave/cli/v2"
)

// DecryptCommand returns the CLI command for decrypting data
func DecryptCommand() *cli.Command {
	return &cli.Command{
		Name:  "decrypt",
		Usage: "Decrypt a file that was encrypted using the encrypt command",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Required: true,
				Usage:    "Input encrypted file path",
				EnvVars:  []string{"INPUT_PATH"},
			},
			&cli.StringFlag{
				Name:    "out-dir",
				Usage:   "Output directory for decrypted files",
				Value:   "./decrypted",
				EnvVars: []string{"OUT_DIR"},
			},
			&cli.StringFlag{
				Name:     "key",
				Required: true,
				Usage:    "Hex-encoded derived key for decryption",
				EnvVars:  []string{"DECRYPT_KEY"},
			},
		},
		Action: decryptAction,
	}

}

func decryptAction(cCtx *cli.Context) error {
	inputPath := cCtx.String("input")
	outDir := cCtx.String("out-dir")
	keyHex := cCtx.String("key")

	// Support both raw hex and hex starting with 0x
	if strings.HasPrefix(keyHex, "0x") {
		keyHex = keyHex[2:]
	}

	// Decode the hex-encoded key
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return fmt.Errorf("failed to decode hex key: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Read encrypted file
	encryptedData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}

	// Decrypt the file with the provided key
	decryptedData, err := encryption.DecryptDataWithKey(encryptedData, key)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Write the decrypted file
	outPath := filepath.Join(outDir, "decrypted_"+filepath.Base(inputPath))
	if err := os.WriteFile(outPath, decryptedData, 0644); err != nil {
		return fmt.Errorf("failed to write decrypted file: %w", err)
	}

	fmt.Printf("File decrypted successfully\n")
	fmt.Printf("Encrypted file: %s\n", inputPath)
	fmt.Printf("Decrypted file: %s\n", outPath)

	return nil
}
