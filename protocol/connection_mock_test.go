package protocol_test

import (
	"bytes"
	"net"
	"time"
)

// MockConn satisfies net.Conn but writes to/reads from buffers.
type MockConn struct {
	ReadBuf  *bytes.Buffer // Data the Connection will "Read" (Input)
	WriteBuf *bytes.Buffer // Data the Connection "Writes" (Output)
}

func NewMockConn() *MockConn {
	return &MockConn{
		ReadBuf:  new(bytes.Buffer),
		WriteBuf: new(bytes.Buffer),
	}
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	return m.ReadBuf.Read(b)
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	return m.WriteBuf.Write(b)
}

// Boilerplate to satisfy net.Conn interface
func (m *MockConn) Close() error                       { return nil }
func (m *MockConn) LocalAddr() net.Addr                { return nil }
func (m *MockConn) RemoteAddr() net.Addr               { return nil }
func (m *MockConn) SetDeadline(t time.Time) error      { return nil }
func (m *MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *MockConn) SetWriteDeadline(t time.Time) error { return nil }
