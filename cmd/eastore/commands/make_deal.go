package commands

import (
	"fmt"
	"math/big"

	"github.com/eastore-project/eastore/pkg/chain"
	"github.com/eastore-project/eastore/pkg/contract"
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
				Name:     "outdir",
				Required: true,
				Usage:    "Output directory for CAR files",
				EnvVars:  []string{"OUT_DIR"},
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
		},
		Action: makeDealAction,
	}
}

func makeDealAction(cCtx *cli.Context) error {
	// Create data prep config
	config := &buffer.Config{
		Type:    cCtx.String("buffer-type"),
		ApiKey:  cCtx.String("buffer-api-key"),
		BaseURL: cCtx.String("buffer-url"),
	}

	// Prepare data using our dataprep package
	prepResult, err := dealutils.PrepareData(
		cCtx.String("input"),
		cCtx.String("outdir"),
		config,
	)
	fmt.Println(prepResult)
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
