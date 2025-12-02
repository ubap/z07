package protocol_test

import (
	"bytes"
	"encoding/binary"
	"goTibia/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

// --- Helper for Mocking an Encodable Packet ---
type MockPacket struct {
	Content string
}

func (mp *MockPacket) Encode(w *protocol.PacketWriter) {
	w.WriteString(mp.Content)
}

func TestConnection_WriteMessage_Unencrypted(t *testing.T) {
	mock := NewMockConn()
	conn := protocol.NewConnection(mock)

	// We are not enabling XTEA, so it should send plain bytes prefixed by length.
	payload := []byte{0xAA, 0xBB, 0xCC}

	err := conn.WriteMessage(payload)
	require.NoError(t, err)

	// Verify Output: [Length U16 (2 bytes)] + [Payload (3 bytes)]
	// Length should be 3 (0x03, 0x00)
	sentBytes := mock.WriteBuf.Bytes()

	require.Len(t, sentBytes, 5)
	require.Equal(t, uint16(3), binary.LittleEndian.Uint16(sentBytes[0:2]), "Length prefix mismatch")
	require.True(t, bytes.Equal(sentBytes[2:], payload))
}

func TestConnection_ReadMessage_Unencrypted(t *testing.T) {
	mock := NewMockConn()
	conn := protocol.NewConnection(mock)

	// Prepare Input: [Length=4] [0xDE, 0xAD, 0xBE, 0xEF]
	payload := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	binary.Write(mock.ReadBuf, binary.LittleEndian, uint16(len(payload)))
	mock.ReadBuf.Write(payload)

	// Act
	readPayload, reader, err := conn.ReadMessage()
	require.NoError(t, err, "ReadMessage failed")

	// Assert Payload
	require.True(t, bytes.Equal(readPayload, payload), "Payload mismatch")

	// Sanity check on reader
	require.Equal(t, byte(0xDE), reader.ReadByte(), "First byte mismatch in PacketReader")
}
