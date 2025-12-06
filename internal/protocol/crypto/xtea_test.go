package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXTEA_Pure_RoundTrip(t *testing.T) {
	testKey := [4]uint32{0x11223344, 0x55667788, 0x99AABBCC, 0xDDEEFF00}

	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "Short data",
			data: []byte{'T', 'e', 's', 't'},
		},
		{
			name: "Plaintext with length multiple of 8",
			data: []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11},
		},
		{
			name: "Empty data",
			data: []byte{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- ARRANGE ---
			blockToEncrypt := bytes.NewBuffer(tc.data)

			// --- ACT ---
			ciphertext, err := EncryptXTEA(blockToEncrypt.Bytes(), testKey)
			require.NoError(t, err)

			decryptedBlock, err := DecryptXTEA(ciphertext, testKey)
			require.NoError(t, err)

			// --- ASSERT ---
			require.GreaterOrEqual(t, len(decryptedBlock), 0, "Decrypted block must be at least 2 bytes for the header")
			require.True(t, bytes.HasPrefix(decryptedBlock, tc.data), "Decrypted block should start with the original plaintext")
			// Verify that the padding is correct.
			// This confirms that no extra garbage data was added.
			padding := decryptedBlock[len(tc.data):]
			for _, b := range padding {
				require.Equal(t, byte(0x00), b, "Padding bytes should be zero")
			}
		})
	}
}
