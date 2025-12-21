package packets

import (
	"goTibia/internal/game/domain"
	"goTibia/internal/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_readPosition(t *testing.T) {
	input := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	pr := protocol.NewPacketReader(input[:])

	position := readPosition(pr)

	require.Equal(t, domain.Position{X: 0x201, Y: 0x403, Z: 0x5}, position)
}
