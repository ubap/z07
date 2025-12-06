package proxy

import (
	"fmt"
	"goTibia/protocol"
	"log"
)

type XTEAPacket interface {
	protocol.Encodable
	GetXTEAKey() [4]uint32
}

// InitSession handles the Client -> Proxy -> Server flow for initial handshake packets.
func InitSession[T XTEAPacket](
	logPrefix string,
	client *protocol.Connection,
	targetAddr string,
	parser func(*protocol.PacketReader) (T, error),
) (T, *protocol.Connection, error) {

	var empty T // Zero value for error returns

	// 1. Read Initial Message
	_, packetReader, err := client.ReadMessage()
	if err != nil {
		return empty, nil, fmt.Errorf("read initial message: %w", err)
	}

	// 2. Parse (using the specific parser provided)
	packet, err := parser(packetReader)
	if err != nil {
		return empty, nil, fmt.Errorf("parse initial packet: %w", err)
	}

	// 3. Connect to Backend
	server, err := ConnectToBackend(targetAddr)
	if err != nil {
		return empty, nil, fmt.Errorf("connect backend: %w", err)
	}

	// 4. Forward Packet
	if err := server.SendPacket(packet); err != nil {
		server.Close()
		return empty, nil, fmt.Errorf("forward packet: %w", err)
	}

	log.Printf("[%s] Session established, forwarding to backend.", logPrefix)

	// 5. Enable Encryption
	key := packet.GetXTEAKey()
	server.EnableXTEA(key)
	client.EnableXTEA(key)

	// Return the parsed packet (in case we need data from it) and the open connection
	return packet, server, nil
}
