package client

import (
	"fmt"
	"net"
	"time"
	"z07/internal/protocol"
)

func ConnectToBackend(address string) (protocol.Connection, error) {
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("backend unavailable at %s: %w", address, err)
	}
	return protocol.NewConnection(conn), nil
}
