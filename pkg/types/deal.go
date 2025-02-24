package types

import "math/big"

// DealRequest represents a deal proposal request
type DealRequest struct {
	PieceCID             []byte        `abi:"piece_cid"`
	PieceSize            uint64        `abi:"piece_size"`
	VerifiedDeal         bool          `abi:"verified_deal"`
	Label                string        `abi:"label"`
	StartEpoch           int64         `abi:"start_epoch"`
	EndEpoch             int64         `abi:"end_epoch"`
	StoragePricePerEpoch *big.Int      `abi:"storage_price_per_epoch"`
	ProviderCollateral   *big.Int      `abi:"provider_collateral"`
	ClientCollateral     *big.Int      `abi:"client_collateral"`
	ExtraParamsVersion   uint64        `abi:"extra_params_version"`
	ExtraParams          ExtraParamsV1 `abi:"extra_params"`
}

// ExtraParamsV1 contains additional parameters for the deal
type ExtraParamsV1 struct {
	LocationRef        string `abi:"location_ref"`
	CarSize            uint64 `abi:"car_size"`
	SkipIPNIAnnounce   bool   `abi:"skip_ipni_announce"`
	RemoveUnsealedCopy bool   `abi:"remove_unsealed_copy"`
}
