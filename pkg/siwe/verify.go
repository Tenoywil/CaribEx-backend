package siwe

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
)

// SIWEMessage represents a parsed SIWE message
type SIWEMessage struct {
	Domain    string
	Address   string
	Statement string
	URI       string
	Version   string
	ChainID   string
	Nonce     string
	IssuedAt  time.Time
}

// VerifySIWEMessage parses, normalizes, and verifies a signed SIWE message.
func VerifySIWEMessage(message, signature string) (bool, SIWEMessage, error) {
	// Normalize message to ensure consistent \n
	message = strings.ReplaceAll(message, "\r\n", "\n")
	message = strings.TrimSpace(message)

	siwe, err := parseSiweMessage(message)
	if err != nil {
		return false, siwe, fmt.Errorf("failed to parse SIWE: %v", err)
	}

	// Decode the signature (remove "0x")
	sigBytes, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return false, siwe, fmt.Errorf("invalid signature format")
	}

	// Ethereum signatures have "v" as the last byte (27/28 or 0/1 offset)
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// Apply EIP-191 prefix
	prefixed := accounts.TextHash([]byte(message))

	pubKey, err := crypto.SigToPub(prefixed, sigBytes)
	if err != nil {
		return false, siwe, fmt.Errorf("failed to recover public key: %v", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey).Hex()

	// Compare recovered address to SIWE message address
	if !strings.EqualFold(recoveredAddr, siwe.Address) {
		return false, siwe, fmt.Errorf("signature mismatch: recovered=%s, expected=%s", recoveredAddr, siwe.Address)
	}

	return true, siwe, nil
}

// parseSiweMessage does minimal parsing of an EIP-4361 message.
func parseSiweMessage(message string) (SIWEMessage, error) {
	var s SIWEMessage
	lines := strings.Split(message, "\n")
	if len(lines) < 6 {
		return s, errors.New("message too short")
	}

	// Line 1: domain + "wants you to sign in..."
	parts := strings.SplitN(lines[0], " wants you to sign in", 2)
	if len(parts) < 1 {
		return s, errors.New("invalid domain line")
	}
	s.Domain = strings.TrimSpace(parts[0])

	// Find the address line (skip empty lines after domain)
	addressLineIndex := 1
	for addressLineIndex < len(lines) && strings.TrimSpace(lines[addressLineIndex]) == "" {
		addressLineIndex++
	}
	if addressLineIndex >= len(lines) {
		return s, errors.New("no address found")
	}
	s.Address = strings.TrimSpace(lines[addressLineIndex])

	// Extract remaining key-value lines
	patterns := map[string]*regexp.Regexp{
		"URI":      regexp.MustCompile(`URI:\s*(.+)`),
		"Version":  regexp.MustCompile(`Version:\s*(.+)`),
		"ChainID":  regexp.MustCompile(`Chain ID:\s*(.+)`),
		"Nonce":    regexp.MustCompile(`Nonce:\s*(.+)`),
		"IssuedAt": regexp.MustCompile(`Issued At:\s*(.+)`),
	}

	for _, line := range lines {
		for key, re := range patterns {
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				switch key {
				case "URI":
					s.URI = matches[1]
				case "Version":
					s.Version = matches[1]
				case "ChainID":
					s.ChainID = matches[1]
				case "Nonce":
					s.Nonce = matches[1]
				case "IssuedAt":
					t, _ := time.Parse(time.RFC3339, strings.TrimSpace(matches[1]))
					s.IssuedAt = t
				}
			}
		}
	}

	return s, nil
}

// VerifySIWE performs complete SIWE verification
func VerifySIWE(message, signature, expectedDomain string) (*SIWEMessage, error) {
	// Use the comprehensive verification function
	isValid, siweMsg, err := VerifySIWEMessage(message, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to verify SIWE message: %w", err)
	}

	if !isValid {
		return nil, fmt.Errorf("signature verification failed")
	}

	// Verify domain matches
	if siweMsg.Domain != expectedDomain {
		return nil, fmt.Errorf("domain mismatch: expected %s, got %s", expectedDomain, siweMsg.Domain)
	}

	// Convert to pointer and return
	return &SIWEMessage{
		Domain:    siweMsg.Domain,
		Address:   siweMsg.Address,
		Statement: siweMsg.Statement,
		URI:       siweMsg.URI,
		Version:   siweMsg.Version,
		ChainID:   siweMsg.ChainID,
		Nonce:     siweMsg.Nonce,
		IssuedAt:  siweMsg.IssuedAt,
	}, nil
}
