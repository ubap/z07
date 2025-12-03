package game

import (
	"goTibia/protocol"
	"goTibia/types"
)

type LoginResponse struct {
	ClientDisconnected bool
	PlayerId           uint32
	BeatDuration       uint16
	CanReportBugs      bool
}

type MapDescription struct {
	Pos types.Position
}

type MoveCreatureMsg struct {
	// The destination is always present
	ToPos types.Position

	// --- CONDITIONAL FIELDS ---
	// Branch A: We know the exact tile and stack position (OldPos < 10)
	FromPos      types.Position
	FromStackPos int8 // -1 if not set

	// Branch B: We only know the Creature ID (OldPos >= 10 or off-screen)
	CreatureID uint32 // 0 if not set

	// Helper to let logic know which fields to use
	KnownSourcePosition bool // True if FromPos/StackPos is valid. False if CreatureID is valid.
}

type PingMsg struct{}

type MagicEffect struct {
	Pos  types.Position
	Type uint8
}

type RemoveTileCreatureMsg struct {
	CreatureID uint32
}

type RemoveTileThingMsg struct {
	Pos      types.Position
	StackPos uint8
}

type AddTileThingMsg struct {
	Pos  types.Position
	Item types.Item
}

type CreatureLightMsg struct {
	CreatureID uint32
	LightLevel uint8
	Color      uint8
}

type WorldLightMsg struct {
	LightLevel uint8
	Color      uint8
}

type CreatureHealthMsg struct {
	CreatureID uint32
	Hppc       uint8
}

type PlayerIconsMsg struct {
	Icons uint8
}

type ContainerMsg struct {
	ContainerID uint8
}

type AddInventoryItemMsg struct {
	Slot uint8
	Item types.Item
}

type RemoveInventoryItemMsg struct {
	Slot uint8
}

type SayMsg struct {
	Slot uint8
}

func ParseLoginResultMessage(pr *protocol.PacketReader) (*LoginResponse, error) {
	lr := &LoginResponse{}

	lr.PlayerId = pr.ReadUint32()
	lr.BeatDuration = pr.ReadUint16()
	lr.CanReportBugs = pr.ReadBool()

	return lr, pr.Err()
}

func ParsePlayerStats(pr *protocol.PacketReader) (*MapDescription, error) {
	return nil, ErrUnknownOpcode
}

func ParseMoveCreature(pr *protocol.PacketReader) (*MoveCreatureMsg, error) {
	msg := &MoveCreatureMsg{
		FromStackPos: -1,
	}

	// 1. Peek at the next 2 bytes.
	// In C++: msg.add<uint16_t>(0xFFFF)
	peekVal, err := pr.PeekUint16()
	if err != nil {
		return nil, err
	}

	// 2. Decide which branch to parse
	if peekVal == 0xFFFF {
		// --- BRANCH B: Unknown Position (CreatureID) ---
		_ = pr.ReadUint16() // Consume the 0xFFFF marker

		msg.CreatureID = pr.ReadUint32()
		msg.KnownSourcePosition = false

		// In this branch, C++ does NOT send FromPos or StackPos
	} else {
		// --- BRANCH A: KnownSourcePosition  ---

		// Read FromPos (X, Y, Z)
		msg.FromPos.X = pr.ReadUint16()
		msg.FromPos.Y = pr.ReadUint16()
		msg.FromPos.Z = pr.ReadByte()

		// Read StackPos
		// In C++, this is msg.addByte(oldStackPos)
		msg.FromStackPos = int8(pr.ReadByte())
		msg.KnownSourcePosition = true
	}

	// 3. Read Destination (Always present)
	msg.ToPos.X = pr.ReadUint16()
	msg.ToPos.Y = pr.ReadUint16()
	msg.ToPos.Z = pr.ReadByte()

	return msg, nil
}

func ParseMagicEffect(pr *protocol.PacketReader) (*MagicEffect, error) {
	me := &MagicEffect{}
	me.Pos = readPosition(pr)
	me.Type = pr.ReadByte()
	return me, nil
}

func ParseRemoveTileThing(pr *protocol.PacketReader) (S2CPacket, error) {
	// 1. Peek
	peekVal, err := pr.PeekUint16()
	if err != nil {
		return nil, err
	}

	// 2. Branch
	if peekVal == 0xFFFF {
		// --- Return Struct A ---
		_ = pr.ReadUint16() // Consume marker

		return &RemoveTileCreatureMsg{
			CreatureID: pr.ReadUint32(),
		}, nil
	}

	// --- Return Struct B ---
	msg := &RemoveTileThingMsg{}
	msg.Pos.X = pr.ReadUint16()
	msg.Pos.Y = pr.ReadUint16()
	msg.Pos.Z = pr.ReadByte()
	msg.StackPos = pr.ReadByte()

	return msg, nil
}

func ParseCreatureLight(pr *protocol.PacketReader) (*CreatureLightMsg, error) {
	cl := &CreatureLightMsg{}
	cl.CreatureID = pr.ReadUint32()
	cl.LightLevel = pr.ReadByte()
	cl.Color = pr.ReadByte()

	return cl, nil
}

func ParseWorldLight(pr *protocol.PacketReader) (*WorldLightMsg, error) {
	cl := &WorldLightMsg{}
	cl.LightLevel = pr.ReadByte()
	cl.Color = pr.ReadByte()

	return cl, nil
}

func ParseCreatureHealth(pr *protocol.PacketReader) (*CreatureHealthMsg, error) {
	cl := &CreatureHealthMsg{}
	cl.CreatureID = pr.ReadUint32()
	cl.Hppc = pr.ReadByte()

	return cl, nil
}

func (lr *WorldLightMsg) Encode(pw *protocol.PacketWriter) {
	pw.WriteByte(S2CWorldLight)
	pw.WriteByte(lr.LightLevel)
	pw.WriteByte(lr.Color)
}

func (cr *CreatureLightMsg) Encode(pw *protocol.PacketWriter) {
	pw.WriteByte(S2CCreatureLight)
	pw.WriteUint32(cr.CreatureID)
	pw.WriteByte(cr.LightLevel)
	pw.WriteByte(cr.Color)
}

func ParsePlayerIcons(pr *protocol.PacketReader) (*PlayerIconsMsg, error) {
	pi := &PlayerIconsMsg{}
	pi.Icons = pr.ReadByte()

	return pi, nil
}

type ServerClosedMsg struct {
	Reason string
}

func ParseServerClosedMsg(pr *protocol.PacketReader) (*ServerClosedMsg, error) {
	scm := &ServerClosedMsg{}
	scm.Reason = pr.ReadString()
	return scm, nil
}

func ParseAddTileThingMsg(pr *protocol.PacketReader) (*AddTileThingMsg, error) {
	ati := &AddTileThingMsg{}
	ati.Pos.X = pr.ReadUint16()
	ati.Pos.Y = pr.ReadUint16()
	ati.Pos.Z = pr.ReadByte()

	itemId, _ := pr.PeekUint16()
	if itemId == 97 || itemId == 98 {
		// TODO do something with it
		readCreatureInMap(pr)
	} else {
		ati.Item = readItem(pr)
	}
	return ati, nil
}

func ParseAddInventoryItemMsg(pr *protocol.PacketReader) (*AddInventoryItemMsg, error) {
	aii := &AddInventoryItemMsg{}
	aii.Slot = pr.ReadByte()
	aii.Item = readItem(pr)
	return aii, nil
}

func ParseRemoveInventoryItemMsg(pr *protocol.PacketReader) (*RemoveInventoryItemMsg, error) {
	rii := &RemoveInventoryItemMsg{}
	rii.Slot = pr.ReadByte()
	return rii, nil
}
