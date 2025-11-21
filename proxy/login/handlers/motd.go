package handlers

import (
	"goTibia/protocol"
	"log"
)

type MOTDHandler struct {
	Opcode uint8
}

func (h *MOTDHandler) Handle(payload protocol.PacketPayload) (protocol.PacketPayload, error) {
	log.Printf("Motd handler 0x%02X", h.Opcode)
	return payload, nil
}
