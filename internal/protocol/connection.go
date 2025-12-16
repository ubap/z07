package protocol

import (
	"encoding/binary"
	"fmt"
	"goTibia/internal/protocol/crypto"
	"io"
	"net"
)

// Connection is a wrapper around a raw network connection (net.Conn)
// that understands the Tibia protocol's message framing.
type Connection struct {
	conn        net.Conn
	XTEAEnabled bool
	XTEAKey     [4]uint32
}

// Encodable represents anything that can write itself to a PacketWriter.
type Encodable interface {
	// Encode writes the packet data to the provided writer.
	// It does not return []byte.
	// It does not return error (errors are stored in the PacketWriter state).
	Encode(pw *PacketWriter)
}

// NewConnection creates a new protocol-aware connection wrapper.
func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn: conn}
}

func (c *Connection) EnableXTEA(key [4]uint32) {
	c.XTEAEnabled = true
	c.XTEAKey = key
}

func (c *Connection) ReadMessage() ([]byte, *PacketReader, error) {
	// Read the header (2 bytes)
	var header [2]byte
	if _, err := io.ReadFull(c.conn, header[:]); err != nil {
		return nil, nil, err
	}
	length := binary.LittleEndian.Uint16(header[:])

	// 2. Read the exact payload
	payload := make([]byte, length)
	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return nil, nil, err
	}

	if c.XTEAEnabled {
		var err error
		payload, err = crypto.DecryptXTEA(payload, c.XTEAKey)
		if err != nil {
			return nil, nil, fmt.Errorf("decryption failed: %w", err)
		}
		if len(payload) < 2 {
			return nil, nil, fmt.Errorf("decrypted payload too short: %d bytes", len(payload))
		}
		innerLength := binary.LittleEndian.Uint16(payload[0:2])

		requiredSize := int(innerLength) + 2
		if requiredSize > len(payload) {
			return nil, nil, fmt.Errorf("malformed packet: inner length %d exceeds buffer size %d", innerLength, len(payload))
		}

		payload = payload[2:requiredSize]
	}

	return payload, NewPacketReader(payload), nil
}

func (c *Connection) WriteMessage(payload []byte) error {
	// TODO: This need mutex

	var dataToSend []byte
	var err error

	if c.XTEAEnabled {
		// 1. Prepend the message length to the payload BEFORE encryption.
		// We create a slice sized [2 bytes for length] + [Payload]
		plaintext := make([]byte, 2+len(payload))

		// Write inner length (size of the actual message)
		binary.LittleEndian.PutUint16(plaintext, uint16(len(payload)))
		// Copy payload after the first 2 bytes
		copy(plaintext[2:], payload)

		// 2. Encrypt the combined block (InnerLength + Payload)
		// 'dataToSend' will now contain the encrypted bytes (likely padded)
		dataToSend, err = crypto.EncryptXTEA(plaintext, c.XTEAKey)
		if err != nil {
			return err
		}
	} else {
		// If encryption is off, we just send the raw payload.
		dataToSend = payload
	}

	// 3. Calculate the TCP Frame Length.
	// This tells the receiver how many bytes to read from the socket.
	// - If Encrypted: Length of the Ciphertext (includes padding + inner length).
	// - If Raw: Length of the Payload.
	frameLength := uint16(len(dataToSend))

	// 4. Write the Frame Length (Header)
	// Optimized: Use stack array instead of binary.Write to avoid allocations.
	header := [2]byte{}
	binary.LittleEndian.PutUint16(header[:], frameLength)

	// Write Header
	if _, err := c.conn.Write(header[:]); err != nil {
		return err
	}

	// 5. Write Body (Encrypted blob or Raw payload)
	_, err = c.conn.Write(dataToSend)
	return err
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) RawConn() net.Conn {
	return c.conn
}

func (c *Connection) SendPacket(packet Encodable) error {
	// 1. Create (or grab from pool) a Writer
	writer := NewPacketWriter()

	// 2. Encode logic
	packet.Encode(writer)

	// 3. Check for logical errors during writing
	if err := writer.Err(); err != nil {
		return err
	}

	// 4. Get the raw bytes (Payload only)
	payload, err := writer.GetBytes()
	if err != nil {
		return err
	}

	// 5. Send to connection (Handles Encryption + Length Prefix)
	return c.WriteMessage(payload)
}
