package game

import (
	"errors"
	"goTibia/protocol"
	"log"
)

// ErrUnknownOpcode is returned when we don't have a parser for this ID.
// The proxy uses this signal to just forward the raw bytes.
var ErrUnknownOpcode = errors.New("unknown opcode")
var ErrNotFullyImplemented = errors.New("not fully implemented")

// S2CPacket is a marker interface for any packet received from Server.
type S2CPacket interface {
	// We can add methods here later, e.g., Name() string
}

func ParseS2CPacket(opcode uint8, pr *protocol.PacketReader) (S2CPacket, error) {
	switch opcode {
	case S2CLoginSuccessful:
		return ParseLoginResultMessage(pr)
	case S2CMapDescription:
		return ParseMapDescription(pr)
	case S2CMoveCreature:
		return ParseMoveCreature(pr)
	case S2CPlayerStats:
		log.Println("player stats")
		return nil, ErrUnknownOpcode

	default:
		return nil, ErrUnknownOpcode
	}
}
