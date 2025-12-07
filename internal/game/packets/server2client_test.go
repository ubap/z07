package packets_test

import (
	"goTibia/internal/game/packets"
	"goTibia/internal/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMagicEffect(t *testing.T) {
	input := []byte{0x64, 0x01, 0x1A, 0x2B, 0x07, 0x08}
	pr := protocol.NewPacketReader(input[:])

	effect, err := packets.ParseMagicEffect(pr)
	require.NoError(t, err)
	require.IsType(t, &packets.MagicEffect{}, effect)
}
