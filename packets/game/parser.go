package game

import (
	"errors"
	"goTibia/protocol"
)

// ErrUnknownOpcode is returned when we don't have a parser for this ID.
// The proxy uses this signal to just forward the raw bytes.
var ErrUnknownOpcode = errors.New("unknown opcode")

// S2CPacket is a marker interface for any packet received from Server.
type S2CPacket interface {
	// We can add methods here later, e.g., Name() string
}

func ParseS2CPacket(opcode uint8, pr *protocol.PacketReader) (S2CPacket, error) {
	switch opcode {
	case S2CLoginSuccessful:
		return ParseLoginResultMessage(pr)

	default:
		return nil, ErrUnknownOpcode
	}
}
