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
func (c *Connection) ReadMessage() ([]byte, error) {
	var length uint16
	// Read the 2-byte length prefix.
	if err := binary.Read(c.conn, binary.LittleEndian, &length); err != nil {
		// An io.EOF here is a clean disconnect.
		return nil, err
	}

	// Read the message body of the specified length.
	payload := make([]byte, length)
	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return nil, err
	}

	if c.XTEAEnabled {
		return crypto.DecryptXTEA(payload, c.XTEAKey)
	}

	return payload, nil
}

// Responsible for framing, encrypting (if enabled), and sending a message.
func (c *Connection) WriteMessage(payload []byte) error {
	var finalPayload []byte
	var err error
	if c.XTEAEnabled {
		finalPayload, err = crypto.EncryptXTEA(payload, c.XTEAKey)
		if err != nil {
			return err
		}
	} else {
		finalPayload = payload
	}

	length := uint16(len(finalPayload))
	// Write the 2-byte length prefix.
	if err := binary.Write(c.conn, binary.LittleEndian, length); err != nil {
		return err
	}

	_, err = c.conn.Write(finalPayload)
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
