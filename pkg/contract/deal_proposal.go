package contract

import (
	"context"
	"fmt"

	"github.com/eastore-project/eastore/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

// MakeDealProposal sends a deal proposal to the smart contract
func (d *DealClient) MakeDealProposal(ctx context.Context, deal types.DealRequest) (common.Hash, error) {
	d.auth.Context = ctx

	tx, err := d.contract.Transact(d.auth, "makeDealProposal", deal)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction: %w", err)
	}
	return tx.Hash(), nil
}
