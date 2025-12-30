package bot

import (
	"testing"
	"z07/internal/game/packets"
	"z07/internal/protocol"

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

func TestInterceptS2CPacket_LoginQueue(t *testing.T) {
	b := &Bot{}
	loginQueuePacket := []byte{byte(packets.S2CSLoginQueue)}

	// 2. Act
	packet, err := b.InterceptS2CPacket(loginQueuePacket)

	// 3. Assert
	require.NoError(t, err)
	pr := protocol.NewPacketReader(packet)

	opcode := pr.ReadUint8()
	assert.Equal(t, uint8(packets.S2CSLoginQueue), opcode)

	msg, err := packets.ParseLoginQueueMsg(pr)
	require.NoError(t, err)

	assert.Equal(t, "Queue hack active.", msg.Message)
	assert.Equal(t, uint8(1), msg.RetryTimeSeconds)
	require.NoError(t, pr.Err())
}
