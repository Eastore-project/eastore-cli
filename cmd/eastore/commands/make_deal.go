package commands

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/eastore-project/eastore/pkg/chain"
	"github.com/eastore-project/eastore/pkg/contract"
	"github.com/eastore-project/eastore/pkg/encryption"
	"github.com/eastore-project/eastore/pkg/types"
	"github.com/eastore-project/fildeal/src/buffer"
	dealutils "github.com/eastore-project/fildeal/src/deal/utils"
	"github.com/ipfs/go-cid"
	"github.com/urfave/cli/v2"
)

const (
	DefaultDuration             = 518400
	DefaultStartEpochHeadOffset = 1000
	DefaultStoragePrice         = "0"
	DefaultProviderCollateral   = "0"
	DefaultClientCollateral     = "0"
	DefaultSkipIPNI             = false
	DefaultRemoveUnsealed       = false
	DefaultVerifiedDeal         = true
	DefaultEncrypted            = false
)

// MakeDealCommand returns the CLI command for making a new deal
func MakeDealCommand() *cli.Command {
	return &cli.Command{
		Name:  "make-deal",
		Usage: "Submit a new deal proposal",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Required: true,
				Usage:    "Input file or folder path",
				EnvVars:  []string{"INPUT_PATH"},
			},
			&cli.StringFlag{
				Name:    "outdir",
				Usage:   "Output directory for CAR files (if not provided, uses temp dir and cleans up after)",
				EnvVars: []string{"OUT_DIR"},
			},
			&cli.StringFlag{
				Name:    "buffer-type",
				Usage:   "Buffer type (lighthouse or local)",
				Value:   "local",
				EnvVars: []string{"BUFFER_TYPE"},
			},
			&cli.StringFlag{
				Name:    "buffer-api-key",
				Usage:   "Buffer service API key",
				EnvVars: []string{"BUFFER_API_KEY"},
			},
			&cli.StringFlag{
				Name:    "buffer-url",
				Usage:   "Buffer service base URL",
				EnvVars: []string{"BUFFER_URL"},
			},
			&cli.Int64Flag{
				Name:    "duration",
				Usage:   "duration of the deal in epochs (default: 518400)",
				Value:   DefaultDuration,
				EnvVars: []string{"DEAL_DURATION"},
			},
			&cli.Int64Flag{
				Name:    "start-epoch-offset",
				Usage:   "start epoch by when the deal should be proved by provider on-chain after current chain head (default: 1000)",
				Value:   DefaultStartEpochHeadOffset,
				EnvVars: []string{"DEAL_START_EPOCH_OFFSET"},
			},
			&cli.Int64Flag{
				Name:    "start-epoch",
				Usage:   "start epoch by when the deal should be proved by provider on-chain (overrides offset)",
				EnvVars: []string{"DEAL_START_EPOCH"},
			},
			&cli.StringFlag{
				Name:    "storage-price",
				Usage:   "storage price in attoFIL per epoch per GiB (default: 0)",
				Value:   DefaultStoragePrice,
				EnvVars: []string{"STORAGE_PRICE"},
			},
			&cli.StringFlag{
				Name:    "provider-collateral",
				Usage:   "deal collateral that storage miner must put in escrow; if empty, the min collateral for the given piece size will be used (default: 0)",
				Value:   DefaultProviderCollateral,
				EnvVars: []string{"PROVIDER_COLLATERAL"},
			},
			// should be calculated for mainnet
			&cli.StringFlag{
				Name:    "client-collateral",
				Usage:   "Client collateral in attoFil",
				Value:   DefaultClientCollateral,
				EnvVars: []string{"CLIENT_COLLATERAL"},
			},
			&cli.BoolFlag{
				Name:    "skip-ipni",
				Usage:   "indicates that deal index should not be announced to the IPNI(Network Indexer) (default: false)",
				Value:   DefaultSkipIPNI,
				EnvVars: []string{"SKIP_IPNI"},
			},
			&cli.BoolFlag{
				Name:    "remove-unsealed",
				Usage:   "indicates that an unsealed copy of the sector in not required for fast retrieval (default: false)",
				Value:   DefaultRemoveUnsealed,
				EnvVars: []string{"REMOVE_UNSEALED"},
			},
			&cli.BoolFlag{
				Name:    "verified-deal",
				Usage:   "whether the deal funds should come from verified client data-cap (default: true)",
				Value:   DefaultVerifiedDeal,
				EnvVars: []string{"VERIFIED_DEAL"},
			},
			&cli.BoolFlag{
				Name:    "encrypted",
				Usage:   "Whether to encrypt the file before making the deal (default: false)",
				Value:   DefaultEncrypted,
				EnvVars: []string{"ENCRYPTED"},
			},
			&cli.StringFlag{
				Name:    "encrypted-out-dir",
				Usage:   "Output directory for encrypted files (if not provided, uses temp dir and cleans up after)",
				EnvVars: []string{"ENCRYPTED_OUT_DIR"},
			},
		},
		Action: makeDealAction,
	}
}

func makeDealAction(cCtx *cli.Context) error {
	inputPath := cCtx.String("input")
	outDir := cCtx.String("outdir")
	isEncrypted := cCtx.Bool("encrypted")
	encryptedOutDir := cCtx.String("encrypted-out-dir")

	// Handle temporary directories
	useTempMain := outDir == ""
	useTempEncrypted := encryptedOutDir == ""
	var tempDirs []string
	var err error

	// Setup main output directory
	if useTempMain {
		outDir, err = os.MkdirTemp("", "eastore-deal-*")
		if err != nil {
			return fmt.Errorf("failed to create temporary directory: %w", err)
		}
		tempDirs = append(tempDirs, outDir)
	} else {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Clean up temporary directories on exit
	defer func() {
		for _, dir := range tempDirs {
			os.RemoveAll(dir)
		}
	}()

	// If encryption is requested, encrypt the file first
	if isEncrypted {
		privateKey := cCtx.String("private-key")

		// Setup encrypted output directory
		if useTempEncrypted {
			encryptedOutDir, err = os.MkdirTemp("", "eastore-encrypt-*")
			if err != nil {
				return fmt.Errorf("failed to create temporary encryption directory: %w", err)
			}
			tempDirs = append(tempDirs, encryptedOutDir)
		} else {
			if err := os.MkdirAll(encryptedOutDir, 0755); err != nil {
				return fmt.Errorf("failed to create encrypted output directory: %w", err)
			}
		}

		// Use the EncryptFile function
		encryptedData, hexKey, err := encryption.EncryptFile(inputPath, privateKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt file: %w", err)
		}

		// Write encrypted file
		encryptedFilePath := filepath.Join(encryptedOutDir, "encrypted_"+filepath.Base(inputPath))
		if err := os.WriteFile(encryptedFilePath, encryptedData, 0644); err != nil {
			return fmt.Errorf("failed to write encrypted file: %w", err)
		}

		fmt.Printf("File encrypted successfully with key: %s\n", hexKey)
		if !useTempEncrypted {
			fmt.Printf("Encrypted file directory: %s\n", encryptedOutDir)
		}

		// Update input path to use encrypted file for the deal
		inputPath = encryptedFilePath
	}

	// Create data prep config
	config := &buffer.Config{
		Type:    cCtx.String("buffer-type"),
		ApiKey:  cCtx.String("buffer-api-key"),
		BaseURL: cCtx.String("buffer-url"),
	}

	// Prepare data using our dataprep package
	prepResult, err := dealutils.PrepareData(
		inputPath,
		outDir,
		config,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare data: %w", err)
	}

	// Calculate epochs
	startEpoch := cCtx.Int64("start-epoch")
	duration := cCtx.Int64("duration")

	if startEpoch == 0 {
		// Fetch chain head and add offset
		head, err := chain.GetChainHead(cCtx.Context, cCtx.String("rpc-url"))
		if err != nil {
			return fmt.Errorf("failed to get chain head: %w", err)
		}
		startEpoch = head + cCtx.Int64("start-epoch-offset")
	}

	endEpoch := startEpoch + duration

	// Create deal request using prep result
	storagePrice, ok := new(big.Int).SetString(cCtx.String("storage-price"), 10)
	if !ok {
		return fmt.Errorf("invalid storage price format")
	}
	providerCollateral, ok := new(big.Int).SetString(cCtx.String("provider-collateral"), 10)
	if !ok {
		return fmt.Errorf("invalid provider collateral format")
	}
	clientCollateral, ok := new(big.Int).SetString(cCtx.String("client-collateral"), 10)
	if !ok {
		return fmt.Errorf("invalid client collateral format")
	}

	c, err := cid.Decode(prepResult.PieceCid)
	if err != nil {
		return fmt.Errorf("failed to decode piece CID: %w", err)
	}

	dealRequest := types.DealRequest{
		PieceCID:             c.Bytes(),
		PieceSize:            prepResult.PieceSize,
		VerifiedDeal:         cCtx.Bool("verified-deal"),
		Label:                prepResult.PayloadCid,
		StartEpoch:           startEpoch,
		EndEpoch:             endEpoch,
		StoragePricePerEpoch: storagePrice,
		ProviderCollateral:   providerCollateral,
		ClientCollateral:     clientCollateral,
		ExtraParamsVersion:   1,
		ExtraParams: types.ExtraParamsV1{
			LocationRef:        prepResult.BufferInfo.URL,
			CarSize:            prepResult.CarSize,
			SkipIPNIAnnounce:   cCtx.Bool("skip-ipni"),
			RemoveUnsealedCopy: cCtx.Bool("remove-unsealed"),
		},
	}

	client, err := contract.NewDealClient(
		cCtx.String("rpc-url"),
		cCtx.String("contract"),
		cCtx.String("private-key"),
	)
	if err != nil {
		return fmt.Errorf("failed to create deal client: %w", err)
	}

	txHash, err := client.MakeDealProposal(cCtx.Context, dealRequest)
	if err != nil {
		return fmt.Errorf("failed to make deal proposal: %w", err)
	}

	fmt.Printf("Deal proposal submitted in transaction: %s\n", txHash.Hex())
	return nil
}
