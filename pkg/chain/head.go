package chain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GetChainHead fetches the current chain head block number
func GetChainHead(ctx context.Context, rpcURL string) (int64, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	defer client.Close()

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch chain head: %w", err)
	}

	return header.Number.Int64(), nil
}
