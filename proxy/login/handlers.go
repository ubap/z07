package login

import (
	"goTibia/protocol"
)

type characterListHandler struct {
	GameProxyIP   string
	GameProxyPort uint16
}

// Handle implements the protocol.PacketHandler interface.
func (h *characterListHandler) Handle(payload protocol.PacketPayload) (protocol.PacketPayload, error) {
	return payload, nil
}
