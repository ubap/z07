package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_EncryptAndDecrypt(t *testing.T) {
	// 1. Arrange
	inputData := []byte("Lorem Ipsum is simply dummy text of the printing and typesetting industry.")

	// 2. Act
	encryptedData, err := EncryptRSA(&RSA.ClientPrivateKey.PublicKey, inputData)
	require.NoError(t, err, "Encryption should not fail")

	decryptedData := DecryptRSA(encryptedData)

	// 3. Assert
	// The decrypted data is a full block. We must verify that our original
	// input data is at the BEGINNING of this block.

	// require.True is a clear and readable way to assert this.
	require.True(t, bytes.HasPrefix(decryptedData, inputData), "Decrypted block should start with the original data")

	// Optional: You can also log the data to see it visually.
	t.Log("Successfully verified that the decrypted block is left-aligned.")
}

func TestEncryptDecrypt_RoundTrip_WithLeadingZero(t *testing.T) {
	// 1. Arrange: Create a data that is guaranteed to start with a zero byte.
	originalPlaintext := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}

	// 2. Act: Encrypt and then decrypt the data.
	ciphertext, err := EncryptRSA(&RSA.ClientPrivateKey.PublicKey, originalPlaintext)
	require.NoError(t, err, "Encryption should not fail")

	decryptedBlock := DecryptRSA(ciphertext)

	// 3. Assert (Improved)

	// 3a. Construct the exact byte block we expect to receive after decryption.
	// It should be a block of the full key size, with our originalPlaintext
	// at the beginning, and the rest filled with zeros.
	keySize := RSA.ClientPrivateKey.Size()
	expectedBlock := make([]byte, keySize) // Creates a slice of all zeros
	copy(expectedBlock, originalPlaintext) // Copies our data to the start

	// 3b. Perform a single, powerful assertion to compare the entire block.
	// This verifies the prefix (our data) and the suffix (the zero-padding).
	require.Equal(t, expectedBlock, decryptedBlock, "Decrypted block should match the original data with correct right-padding")

	t.Log("Successfully verified that a message with a leading zero is correctly left-aligned and padded.")
}

// TestEncryptRSA_IsLeftAligned confirms that our EncryptRSA function correctly
// creates a left-aligned, right-padded block, mimicking the behavior of the C++ client.
func TestEncryptRSA_IsLeftAligned(t *testing.T) {
	// 1. Arrange: Define a short data message.
	// A short message is crucial because a long one might fill the entire block,
	// hiding the padding behavior.
	originalPlaintext := []byte{0x00, 0xAA, 0xBB, 0xCC, 0xDD}

	// 2. Act: Encrypt the message using the function we want to test.
	ciphertext, err := EncryptRSA(&RSA.ClientPrivateKey.PublicKey, originalPlaintext)
	require.NoError(t, err, "EncryptRSA should not produce an error")

	decryptedBlock := DecryptRSA(ciphertext)

	// 3. Assert: Verify the structure of the decrypted block.

	// Assertion 3a: The decrypted block should have the full key size.
	keySize := RSA.ClientPrivateKey.Size()
	require.Equal(t, keySize, len(decryptedBlock), "Decrypted block must have the full key size")

	// Assertion 3b: The original data must be at the BEGINNING of the block.
	// This is the core test for left-alignment.
	require.True(t, bytes.HasPrefix(decryptedBlock, originalPlaintext), "Decrypted block should start with the original data")

	// Assertion 3c: The bytes immediately following the data should be zeros.
	// This confirms that it was correctly right-padded.
	paddingStart := len(originalPlaintext)
	expectedPadding := make([]byte, keySize-paddingStart) // A slice of all zeros.

	actualPadding := decryptedBlock[paddingStart:]
	require.Equal(t, expectedPadding, actualPadding, "The remainder of the block should be zero-byte padding")

	t.Log("Test passed: EncryptRSA correctly produces a left-aligned, right-padded block.")
}
