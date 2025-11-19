package main

import (
	"crypto/rsa"
	"fmt"
	"math/big"
	"strconv"
)

// ParseTibiaRSAPublicKey takes a modulus (as a decimal string) and an exponent
// and constructs a valid *rsa.PublicKey.
func ParseTibiaRSAPublicKey(modulusStr, exponentStr string) (*rsa.PublicKey, error) {
	n := new(big.Int)
	n, ok := n.SetString(modulusStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse modulus string")
	}

	e64, err := strconv.ParseInt(exponentStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse exponent string: %w", err)
	}

	return &rsa.PublicKey{N: n, E: int(e64)}, nil
}
