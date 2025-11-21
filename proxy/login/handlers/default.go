package handlers

import (
	"goTibia/protocol"
	"log"
)

type DefaultPassThroughHandler struct {
	Opcode uint8
}

// Handle implements the PacketHandler interface.
func (h *DefaultPassThroughHandler) Handle(payload protocol.PacketPayload) (protocol.PacketPayload, error) {
	log.Printf("Forwarding unknown opcode 0x%02X", h.Opcode)
	return payload, nil
}
