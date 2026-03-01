package client

import (
	"fmt"
	"log"
	"z07/internal/game/packets"
	"z07/internal/protocol"
)

func Game(targetAddr string) error {
	server, err := ConnectToBackend(targetAddr)
	defer server.Close()
	if err != nil {
		return fmt.Errorf("connect backend: %w", err)
	} else {
		log.Printf("connected to backend: %s", targetAddr)
	}

	packet := packets.NewGameLoginRequest(1, "God", "1")

	if err := server.SendPacket(packet); err != nil {
		server.Close()
		return fmt.Errorf("forward packet: %w", err)
	}

	log.Printf("[Game] Sent credentials, awaiting response...")

	// 5. Enable Encryption
	key := packet.GetXTEAKey()
	server.EnableXTEA(key)

	message, err := server.ReadMessage()
	if err != nil {
		return err
	}
	packetReader := protocol.NewPacketReader(message)
	resultMessage, err := packets.ParseLoginResultMessage(packetReader)
	if err != nil {
		return err
	}

	log.Printf("[Game] Received response: %x", resultMessage)
	return nil
}
