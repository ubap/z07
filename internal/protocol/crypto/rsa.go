package crypto

import (
	"crypto/rsa"
	"fmt"
	"math/big"
)

type RSAKeySet struct {
	ClientPrivateKey    *rsa.PrivateKey
	GameServerPublicKey *rsa.PublicKey
}

var RSA RSAKeySet

// TODO: The following should be loaded from config
const (
	OTPublicRSA     = "109120132967399429278860960508995541528237502902798129123468757937266291492576446330739696001110603907230888610072655818825358503429057592827629436413108566029093628212635953836686562675849720620786279431090218017681061521755056710823876476444260558147179707119674283982419152118103759076030616683978566631413"
	OTPrivateRSA    = "46730330223584118622160180015036832148732986808519344675210555262940258739805766860224610646919605860206328024326703361630109888417839241959507572247284807035235569619173792292786907845791904955103601652822519121908367187885509270025388641700821735345222087940578381210879116823013776808975766851829020659073"
	targetServerRSA = "138358917549655551601135922545920258651079249320630202917602000570926337770168654400102862016157293631277888588897291561865439132767832236947553872456033140205555218536070792283327632773558457562430692973109061064849319454982125688743198270276394129121891795353179249782548271479625552587457164097090236827371"
)

func init() {
	privateKey, err := buildPrivateKeyFromComponents(OTPublicRSA, OTPrivateRSA)
	if err != nil {
		panic(fmt.Sprintf("FATAL: Could not build RSA private key: %v", err))
	}
	RSA.ClientPrivateKey = privateKey

	publicKey, err := buildPublicKeyFromComponents(targetServerRSA)
	if err != nil {
		panic(fmt.Sprintf("FATAL: Could not build RSA public key: %v", err))
	}
	RSA.GameServerPublicKey = publicKey
}

func buildPublicKeyFromComponents(nStr string) (*rsa.PublicKey, error) {
	// 1. Parse the public modulus string into a big.Int.
	n, ok := new(big.Int).SetString(nStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse public modulus (N)")
	}
	return &rsa.PublicKey{
		N: n,
		E: 65537,
	}, nil
}

func buildPrivateKeyFromComponents(nStr, dStr string) (*rsa.PrivateKey, error) {
	// 1. Parse the public modulus string into a big.Int.
	publicKey, err := buildPublicKeyFromComponents(nStr)
	if err != nil {
		return nil, err
	}

	// 2. Parse the private exponent string into a big.Int.
	d, ok := new(big.Int).SetString(dStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse private exponent (D)")
	}

	// 3. Assemble the rsa.PrivateKey struct.
	// We only populate the essential fields for decryption. The public exponent 'E'
	// is almost always 65537, so it's a safe value to hardcode.
	// The Primes and Precomputed fields are left nil.
	privKey := &rsa.PrivateKey{
		PublicKey: *publicKey,
		D:         d,
	}

	return privKey, nil
}

func DecryptRSA(encryptedData []byte) []byte {
	c := new(big.Int).SetBytes(encryptedData)
	m := new(big.Int).Exp(c, RSA.ClientPrivateKey.D, RSA.ClientPrivateKey.N)

	keySize := RSA.ClientPrivateKey.Size()
	decryptedData := make([]byte, keySize)
	m.FillBytes(decryptedData)

	return decryptedData
}

func EncryptRSA(pubKey *rsa.PublicKey, data []byte) ([]byte, error) {
	keySize := pubKey.Size()
	// 0. Input validation
	if len(data) > keySize {
		return nil, fmt.Errorf("crypto failed: data length (%d) is greater than the key size (%d)", len(data), keySize)
	}

	// 1. Create a new buffer of the exact key size.
	paddedPlaintext := make([]byte, keySize)

	// 2. Copy the data to the BEGINNING of the buffer.
	// This creates the required left-aligned, right-padded message.
	copy(paddedPlaintext, data)

	// 3. Encrypt this full, padded block.
	m := new(big.Int).SetBytes(paddedPlaintext)
	e := big.NewInt(int64(pubKey.E))
	c := new(big.Int).Exp(m, e, pubKey.N)

	// 4. Return the ciphertext, ensuring it's also the full key size.
	ciphertext := make([]byte, keySize)
	c.FillBytes(ciphertext)

	return ciphertext, nil
}
