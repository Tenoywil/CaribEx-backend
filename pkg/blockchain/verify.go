package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TransactionVerification contains the result of transaction verification
type TransactionVerification struct {
	TxHash    string `json:"txHash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	ChainID   int64  `json:"chainId"`
	Verified  bool   `json:"verified"`
	IsPending bool   `json:"isPending"`
	Status    uint64 `json:"status"` // 1 = success, 0 = failed
}

// VerifyTransaction validates that a transaction exists, is confirmed, and matches the intended parameters
func VerifyTransaction(txHash string, expectedChainID int64) (*TransactionVerification, error) {
	if client == nil {
		return nil, fmt.Errorf("RPC client not initialized - please configure RPC_URL environment variable and restart the server")
	}

	ctx := context.Background()
	hash := common.HexToHash(txHash)

	// Get transaction details
	tx, isPending, err := client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// Get transaction receipt (only available for confirmed transactions)
	receipt, err := client.TransactionReceipt(ctx, hash)
	if err != nil {
		// Transaction might still be pending
		if isPending {
			return &TransactionVerification{
				TxHash:    txHash,
				From:      "",
				To:        "",
				Value:     "0",
				ChainID:   expectedChainID,
				Verified:  false,
				IsPending: true,
				Status:    0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	// Verify chain ID matches
	actualChainID := tx.ChainId()
	if actualChainID != nil && actualChainID.Int64() != expectedChainID {
		return nil, fmt.Errorf("chain ID mismatch: expected %d, got %d", expectedChainID, actualChainID.Int64())
	}

	// Extract transaction details
	from, err := client.TransactionSender(ctx, tx, receipt.BlockHash, receipt.TransactionIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction sender: %w", err)
	}

	var toAddr string
	if tx.To() != nil {
		toAddr = tx.To().Hex()
	}

	verification := &TransactionVerification{
		TxHash:    txHash,
		From:      from.Hex(),
		To:        toAddr,
		Value:     tx.Value().String(),
		ChainID:   expectedChainID,
		Verified:  receipt.Status == 1,
		IsPending: isPending,
		Status:    receipt.Status,
	}

	// Check for success
	if receipt.Status != 1 {
		return verification, fmt.Errorf("transaction failed on-chain")
	}

	return verification, nil
}

// ValidateChainID checks if the chain ID is in the list of supported networks
func ValidateChainID(chainID int64) bool {
	// Supported networks: Ethereum Mainnet (1), Sepolia (11155111), etc.
	supportedChains := map[int64]bool{
		1:        true, // Ethereum Mainnet
		11155111: true, // Sepolia Testnet
		137:      true, // Polygon Mainnet
		80001:    true, // Mumbai (Polygon Testnet)
		// Note: Goerli (5) was removed as it's deprecated and shut down
	}
	return supportedChains[chainID]
}

// FormatValue converts wei value to a human-readable format
func FormatValue(weiValue string) (string, error) {
	wei := new(big.Int)
	wei, ok := wei.SetString(weiValue, 10)
	if !ok {
		return "", fmt.Errorf("invalid wei value")
	}

	// Convert to ETH (divide by 10^18)
	eth := new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		new(big.Float).SetInt(big.NewInt(1e18)),
	)

	return eth.String(), nil
}
