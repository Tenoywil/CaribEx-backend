package blockchain

import (
	"testing"
)

func TestValidateChainID(t *testing.T) {
	tests := []struct {
		name     string
		chainID  int64
		expected bool
	}{
		{"Ethereum Mainnet", 1, true},
		{"Sepolia", 11155111, true},
		{"Polygon", 137, true},
		{"Mumbai", 80001, true},
		{"Goerli", 5, true},
		{"Invalid Chain", 999999, false},
		{"Zero Chain", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateChainID(tt.chainID)
			if result != tt.expected {
				t.Errorf("ValidateChainID(%d) = %v, expected %v", tt.chainID, result, tt.expected)
			}
		})
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name      string
		weiValue  string
		expectErr bool
	}{
		{"Valid 1 ETH", "1000000000000000000", false},
		{"Valid 0.5 ETH", "500000000000000000", false},
		{"Zero", "0", false},
		{"Invalid string", "invalid", true},
		{"Empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatValue(tt.weiValue)
			if tt.expectErr {
				if err == nil {
					t.Errorf("FormatValue(%s) expected error but got none", tt.weiValue)
				}
			} else {
				if err != nil {
					t.Errorf("FormatValue(%s) unexpected error: %v", tt.weiValue, err)
				}
				if result == "" {
					t.Errorf("FormatValue(%s) returned empty string", tt.weiValue)
				}
			}
		})
	}
}
