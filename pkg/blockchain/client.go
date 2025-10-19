package blockchain

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client

// InitRPC initializes the Ethereum RPC client
func InitRPC(rpcURL string) error {
	var err error
	client, err = ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum RPC: %w", err)
	}

	// Test connection
	_, err = client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID from RPC: %w", err)
	}

	log.Printf("Successfully connected to Ethereum RPC at %s", rpcURL)
	return nil
}

// GetClient returns the global RPC client instance
func GetClient() *ethclient.Client {
	return client
}

// Close closes the RPC client connection
func Close() {
	if client != nil {
		client.Close()
	}
}
