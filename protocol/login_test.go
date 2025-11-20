package protocol

import (
	"goTibia/protocol/crypto"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// We want to use the known test RSA keys for this test.
	crypto.RSA.GameServerPublicKey = &crypto.RSA.ClientPrivateKey.PublicKey

	packet := LoginPacket{
		Protocol:      1,
		ClientOS:      65535,
		ClientVersion: 1234,
		DatSignature:  7,
		SprSignature:  8,
		PicSignature:  9,
		XTEAKey:       [4]uint32{17, 18, 19, 20},
		AccountNumber: 42,
		Password:      "secret",
	}

	marshal, err := packet.Marshal()
	if err != nil {
		t.Fatalf("Error marshalling login packet: %v", err)
	}

	loginPacket, err := ParseLoginPacket(marshal)
	if err != nil {
		t.Fatalf("Error parsing login packet: %v", err)
	}

	require.Equal(t, packet, *loginPacket, "Parsed packet does not match original")
}
