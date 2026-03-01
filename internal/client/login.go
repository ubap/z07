package client

import (
	"fmt"
	"log"
	"z07/internal/login/packets"
	"z07/internal/protocol"
)

func Login(targetAddr string) error {
	server, err := ConnectToBackend(targetAddr)
	if err != nil {
		return fmt.Errorf("connect backend: %w", err)
	} else {
		log.Printf("connected to backend: %s", targetAddr)
	}

	packet := packets.NewClientCredentialPacket(1, "1")

	if err := server.SendPacket(packet); err != nil {
		server.Close()
		return fmt.Errorf("forward packet: %w", err)
	}

	log.Printf("[Login] Sent credentials, awaiting response...")

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

	log.Printf("[Login] Received response: %x", resultMessage)
	return nil
}
