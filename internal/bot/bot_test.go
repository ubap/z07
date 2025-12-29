package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterceptC2SPacket_Look(t *testing.T) {
	b := &Bot{}
	lookPacket := []byte{0x8C, 0x69, 0x7D, 0xE5, 0x7D, 0x07, 0xBA, 0x11, 0x01}

	// 2. Act
	packet, err := b.InterceptC2SPacket(lookPacket)

	// 3. Assert
	require.NoError(t, err)
	assert.Equal(t, uint16(0x11BA), b.lastLookedAt)
	assert.Equal(t, lookPacket, packet)
}
