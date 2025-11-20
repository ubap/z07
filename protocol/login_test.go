package protocol

import (
	"goTibia/protocol/crypto"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// We want to use the known test RSA keys so we can decrypt the packet
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
	require.NoError(t, err, "Marshalling login packet should not fail")

	loginPacket, err := ParseLoginPacket(marshal)
	require.NoError(t, err, "Failed to parse login packet")

	require.Equal(t, packet, *loginPacket, "Parsed packet does not match original")
}
