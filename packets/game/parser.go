package game

import (
	"errors"
	"goTibia/protocol"
)

// ErrUnknownOpcode is returned when we don't have a parser for this ID.
// The proxy uses this signal to just forward the raw bytes.
var ErrUnknownOpcode = errors.New("unknown opcode")
var ErrNotFullyImplemented = errors.New("not fully implemented")

// S2CPacket is a marker interface for any packet received from Server.
type S2CPacket interface {
	// We can add methods here later, e.g., Name() string
	// TODO consider adding Opcode() uint8 to identify the packet type
	// Opcode() uint8
}

func ParseS2CPacket(opcode uint8, pr *protocol.PacketReader) (S2CPacket, error) {
	switch opcode {
	case S2CLoginSuccessful:
		return ParseLoginResultMessage(pr)
	case S2CMapDescription:
		return ParseMapDescriptionMsg(pr)
	case S2CMoveCreature:
		return ParseMoveCreature(pr)
	case S2CPing:
		return &PingMsg{}, nil
	case S2CMagicEffect:
		return ParseMagicEffect(pr)
	case S2CAddTileThing:
		return ParseAddTileThingMsg(pr)
	case S2CRemoveTileThing:
		return ParseRemoveTileThing(pr)
	case S2CAddInventoryItem:
		return ParseAddInventoryItemMsg(pr)
	case S2CRemoveInventoryItem:
		return ParseRemoveInventoryItemMsg(pr)
	case S2CWorldLight:
		return ParseWorldLight(pr)
	case S2CCreatureLight:
		return ParseCreatureLight(pr)
	case S2CCreatureHealth:
		return ParseCreatureHealth(pr)
	case S2CPlayerIcons:
		return ParsePlayerIcons(pr)
	case S2CServerClosed:
		return ParseServerClosedMsg(pr)

	default:
		return nil, ErrUnknownOpcode
	}
}
