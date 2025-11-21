package login

import (
	"goTibia/protocol"
	"io"
	"log"
)

func init() {
	S2CHandlers.Register(ServerOpcodeDisconnectClient, &DisconnectClientHandler{})
}

type DisconnectClientHandler struct {
}

func (h *DisconnectClientHandler) Handle(r io.Reader) error {

	readString, err := protocol.ReadString(r)
	if err != nil {
		return err

	}
	log.Print("DisconnectClientHandler: " + readString)
	return nil
}
