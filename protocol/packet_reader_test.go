package protocol_test

import (
	"bytes"
	"encoding/binary"
	"goTibia/protocol"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPacketReader_HappyPath(t *testing.T) {
	// Prepare test data:
	// Byte: 10
	// Bool: true (1)
	// Uint16: 500
	// Uint32: 100000
	// String: "Hello" (Len 5 + bytes)
	buf := new(bytes.Buffer)
	buf.WriteByte(10)
	buf.WriteByte(1)
	binary.Write(buf, binary.LittleEndian, uint16(500))
	binary.Write(buf, binary.LittleEndian, uint32(100000))

	str := "Hello"
	binary.Write(buf, binary.LittleEndian, uint16(len(str)))
	buf.WriteString(str)

	// Start Reading
	pr := protocol.NewPacketReader(buf.Bytes())

	require.Equal(t, byte(10), pr.ReadByte())
	require.Equal(t, true, pr.ReadBool())
	require.Equal(t, uint16(500), pr.ReadUint16())
	require.Equal(t, uint32(100000), pr.ReadUint32())
	require.Equal(t, "Hello", pr.ReadString())
	require.NoError(t, pr.Err())
}

func TestPacketReader_StickyError(t *testing.T) {
	// Only provide 2 bytes
	data := []byte{0x01, 0x02}
	pr := protocol.NewPacketReader(data)

	// 1. Read successfully
	_ = pr.ReadByte()
	_ = pr.ReadByte()

	// 2. Trigger EOF (Try to read byte from empty reader)
	val := pr.ReadByte()
	require.Equal(t, byte(0), val)
	require.Error(t, io.EOF, pr.Err())

	// 3. Verify Sticky Error (Subsequent calls should return immediately)
	// Even though we are trying to read Uint32, it shouldn't try to parse anything
	// It should just return 0 and keep the original EOF error.
	u32 := pr.ReadUint32()
	require.Equal(t, uint32(0), u32, "Expected 0 on sticky error")
	require.Error(t, io.EOF, pr.Err(), "Error changed! Expected io.EOF")
}

func TestPacketReader_ReadString_EdgeCases(t *testing.T) {
	t.Run("Empty String", func(t *testing.T) {
		// Length 0 (2 bytes)
		data := []byte{0x00, 0x00}
		pr := protocol.NewPacketReader(data)
		require.Equal(t, "", pr.ReadString())
		require.NoError(t, pr.Err())
	})

	t.Run("Not enough bytes for length", func(t *testing.T) {
		data := []byte{0x05} // Need 2 bytes for length, have 1
		pr := protocol.NewPacketReader(data)
		require.Equal(t, "", pr.ReadString())
		require.Error(t, pr.Err(), io.ErrUnexpectedEOF)
	})

	t.Run("Not enough bytes for body", func(t *testing.T) {
		// Length 5, but body is empty
		data := []byte{0x05, 0x00}
		pr := protocol.NewPacketReader(data)
		require.Equal(t, "", pr.ReadString())
		require.Error(t, pr.Err(), io.EOF)
	})
}

func TestPacketReader_SkipAndRemaining(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	pr := protocol.NewPacketReader(data)

	require.Equal(t, 5, pr.Remaining())
	pr.Skip(2)
	require.Equal(t, 3, pr.Remaining())
	require.Equal(t, byte(3), pr.ReadByte())
}

func TestPacketReader_ReadAll(t *testing.T) {
	data := []byte{1, 2, 3, 4}
	pr := protocol.NewPacketReader(data)

	// Read one byte first
	b := pr.ReadByte()
	require.Equal(t, b, byte(1))

	// ReadAll should get the rest
	rest := pr.ReadAll()
	require.Len(t, rest, 3)
	require.True(t, bytes.Equal(rest, []byte{2, 3, 4}))

	// ReadAll on empty should return empty
	empty := pr.ReadAll()
	require.Len(t, empty, 0)
}

func TestPacketReader_PeekUint16(t *testing.T) {
	// Let's create data: [0x01, 0x00, 0x02, 0x00]
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint16(2))

	pr := protocol.NewPacketReader(buf.Bytes())

	// Read the first value (1)
	val1 := pr.ReadUint16()
	require.Equal(t, uint16(1), val1)

	// Peek. We EXPECT to see the next value (2),
	peekVal, err := pr.PeekUint16()
	require.NoError(t, err)

	require.Equal(t, uint16(2), peekVal)

	// Ensure Peek didn't advance cursor
	peekVal = pr.ReadUint16()
	require.Equal(t, uint16(2), peekVal)
}
