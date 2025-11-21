package handlers

import (
	"fmt"
	"goTibia/protocol"
	"io"
	"log"
)

type MOTDHandler struct {
}

func (h *MOTDHandler) Handle(r io.Reader) error {
	data, err := protocol.ReadString(r)
	if err != nil {
		return err
	}

	format := "%d\n%s"

	var motdId int
	var message string

	_, err = fmt.Sscanf(data, format, &motdId, &message)
	if err != nil {
		return fmt.Errorf("failed to parse MOTD data: %w", err)
	}

	log.Printf("MOTD: %s", message)
	return nil
}
