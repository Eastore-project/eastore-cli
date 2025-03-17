package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"strings"

	pkgabi "github.com/eastore-project/eastore/pkg/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type DealClient struct {
	client       *ethclient.Client
	contract     *bind.BoundContract
	contractAddr common.Address
	auth         *bind.TransactOpts
	abi          abi.ABI
	privateKey   *ecdsa.PrivateKey
}

func NewDealClient(rpcURL, contractAddress, privateKey string) (*DealClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(pkgabi.DealClientABI)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	addr := common.HexToAddress(contractAddress)
	contract := bind.NewBoundContract(addr, parsedABI, client, client, client)

	privateKeyECDSA, err := crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	return &DealClient{
		client:       client,
		contract:     contract,
		contractAddr: addr,
		auth:         auth,
		abi:          parsedABI,
		privateKey:   privateKeyECDSA,
	}, nil
}

// SignMessage signs a message with the client's private key and returns the signature
func (c *DealClient) SignMessage(message string) ([]byte, error) {
	// Create hash of the message
	msgHash := crypto.Keccak256Hash([]byte(message))

	// Sign the hash with the private key
	signature, err := crypto.Sign(msgHash.Bytes(), c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	return signature, nil
}
