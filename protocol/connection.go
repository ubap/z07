package protocol

import (
	"encoding/binary"
	"goTibia/protocol/crypto"
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

// NewConnection creates a new protocol-aware connection wrapper.
func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn: conn}
}

func (c *Connection) EnableXTEA(key [4]uint32) {
	c.XTEAEnabled = true
	c.XTEAKey = key
}

// ReadMessage reads a single, complete message from the stream.
// It handles the 2-byte length prefix and returns the message payload.
func (c *Connection) ReadMessage() (*PacketReader, error) {
	// 1. Read the header (2 bytes)
	var header [2]byte
	if _, err := io.ReadFull(c.conn, header[:]); err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint16(header[:])

	// 2. Read the exact payload
	payload := make([]byte, length)
	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return nil, err
	}

	// 3. Decrypt if necessary (Linear flow, no else block)
	if c.XTEAEnabled {
		var err error
		payload, err = crypto.DecryptXTEA(payload, c.XTEAKey)
		if err != nil {
			return nil, err
		}
	}

	return NewPacketReader(payload), nil
}

func (c *Connection) WriteMessage(payload []byte) error {
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

// Close simply closes the underlying network connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}

// RemoteAddr returns the remote network address.
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) RawConn() net.Conn {
	return c.conn
}
