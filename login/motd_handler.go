package login

import (
	"goTibia/protocol"
	"io"
	"log"
	"strconv"
	"strings"
)

func init() {
	// This single line registers our handler with the global registry.
	S2CHandlers.Register(ServerOpcodeMOTD, &MOTDHandler{})
}

type MOTDHandler struct {
}

func (h *MOTDHandler) Handle(r io.Reader) error {
	data, err := protocol.ReadString(r)
	if err != nil {
		return err
	}

	parts := strings.SplitN(data, "\n", 2)

	if len(parts) != 2 {
		log.Fatalf("Invalid format: expected a number followed by a newline and a message.")
		return nil
	}

	motdId, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("Failed to parse MOTD ID: %v", err)
		return nil
	}

	message := parts[1]

	log.Printf("MOTDID :%d, MOTD: %s", motdId, message)
	return nil
}
