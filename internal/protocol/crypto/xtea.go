package crypto

import (
	"encoding/binary"
	"fmt"
)

const (
	xteaRounds = 32
	xteaDelta  = 0x9E3779B9
)

// DecryptXTEA takes an XTEA-encrypted payload and the session key.
// It decrypts the entire ciphertext and returns the decrypted block,
// including any internal headers or padding.
func DecryptXTEA(ciphertext []byte, key [4]uint32) ([]byte, error) {
	// 1. Input Validation: Ciphertext must be a multiple of the block size (8 bytes).
	if len(ciphertext)%8 != 0 {
		return nil, fmt.Errorf("invalid ciphertext length: must be a multiple of 8, but got %d", len(ciphertext))
	}

	// Create a buffer to hold the fully decrypted data.
	plaintext := make([]byte, len(ciphertext))
	keySlice := key[:] // Convert array to slice for the helper

	// 2. Decrypt the ciphertext in 8-byte blocks (ECB mode).
	for i := 0; i < len(ciphertext); i += 8 {
		// Get the current 8-byte block.
		block := ciphertext[i : i+8]

		// Convert the block to two uint32s for the cipher.
		v0 := binary.LittleEndian.Uint32(block[0:4])
		v1 := binary.LittleEndian.Uint32(block[4:8])

		// Decrypt the block using the core XTEA algorithm.
		decryptedV0, decryptedV1 := decryptBlock(v0, v1, keySlice)

		// Put the decrypted uint32s back into our data buffer.
		binary.LittleEndian.PutUint32(plaintext[i:i+4], decryptedV0)
		binary.LittleEndian.PutUint32(plaintext[i+4:i+8], decryptedV1)
	}

	// 3. Return the entire decrypted block.
	return plaintext, nil
}

// decryptBlock performs the core XTEA cipher on a single 64-bit block.
// This helper function remains the same.
func decryptBlock(v0, v1 uint32, key []uint32) (uint32, uint32) {
	// Create a runtime variable for delta to prevent compile-time overflow.
	delta := uint32(xteaDelta)
	sum := delta * xteaRounds

	for i := 0; i < xteaRounds; i++ {
		v1 -= (((v0 << 4) ^ (v0 >> 5)) + v0) ^ (sum + key[(sum>>11)&3])
		sum -= xteaDelta
		v0 -= (((v1 << 4) ^ (v1 >> 5)) + v1) ^ (sum + key[sum&3])
	}

	return v0, v1
}

// EncryptXTEA takes a data payload and the session key. It pads the
// data to a multiple of 8 bytes, then encrypts it and returns the ciphertext.
// It does NOT prepend any length headers.
func EncryptXTEA(plaintext []byte, key [4]uint32) ([]byte, error) {
	// 1. Calculate the padding needed to make the data a multiple of 8.
	paddingBytesNeeded := (8 - (len(plaintext) % 8)) % 8

	// Create a new buffer large enough for the data and padding.
	paddedMessage := make([]byte, len(plaintext)+paddingBytesNeeded)

	// Copy the data into the buffer.
	copy(paddedMessage, plaintext)

	// The rest of the buffer is already zeros, which serves as our padding.

	// 2. Encrypt the fully padded message in 8-byte blocks.
	ciphertext := make([]byte, len(paddedMessage))
	keySlice := key[:]

	for i := 0; i < len(paddedMessage); i += 8 {
		block := paddedMessage[i : i+8]

		v0 := binary.LittleEndian.Uint32(block[0:4])
		v1 := binary.LittleEndian.Uint32(block[4:8])

		encryptedV0, encryptedV1 := encryptBlock(v0, v1, keySlice)

		binary.LittleEndian.PutUint32(ciphertext[i:i+4], encryptedV0)
		binary.LittleEndian.PutUint32(ciphertext[i+4:i+8], encryptedV1)
	}

	return ciphertext, nil
}

// encryptBlock performs the core XTEA cipher on a single 64-bit block.
// This helper function remains the same.
func encryptBlock(v0, v1 uint32, key []uint32) (uint32, uint32) {
	sum := uint32(0)
	for i := 0; i < xteaRounds; i++ {
		v0 += (((v1 << 4) ^ (v1 >> 5)) + v1) ^ (sum + key[sum&3])
		sum += xteaDelta
		v1 += (((v0 << 4) ^ (v0 >> 5)) + v0) ^ (sum + key[(sum>>11)&3])
	}

	return v0, v1
}
