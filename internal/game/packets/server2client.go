package packets

import (
	"goTibia/internal/game/domain"
	"goTibia/internal/protocol"
)

type LoginResponse struct {
	ClientDisconnected bool
	PlayerId           uint32
	BeatDuration       uint16
	CanReportBugs      bool
}

type MoveCreatureMsg struct {
	// The destination is always present
	ToPos domain.Position

	// --- CONDITIONAL FIELDS ---
	// Branch A: We know the exact tile and stack position (OldPos < 10)
	FromPos      domain.Position
	FromStackPos int8 // -1 if not set

	// Branch B: We only know the Creature ID (OldPos >= 10 or off-screen)
	CreatureID uint32 // 0 if not set

	// Helper to let logic know which fields to use
	KnownSourcePosition bool // True if FromPos/StackPos is valid. False if CreatureID is valid.
}

type PingMsg struct{}

type MagicEffect struct {
	Pos  domain.Position
	Type uint8
}

type RemoveTileCreatureMsg struct {
	CreatureID uint32
}

type RemoveTileThingMsg struct {
	Pos      domain.Position
	StackPos uint8
}

type AddTileThingMsg struct {
	Pos  domain.Position
	Item domain.Item
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

type AddInventoryItemMsg struct {
	Slot domain.EquipmentSlot
	Item domain.Item
}

type RemoveInventoryItemMsg struct {
	Slot domain.EquipmentSlot
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
	aii.Slot = domain.EquipmentSlot(pr.ReadByte())
	aii.Item = readItem(pr)
	return aii, nil
}

func ParseRemoveInventoryItemMsg(pr *protocol.PacketReader) (*RemoveInventoryItemMsg, error) {
	rii := &RemoveInventoryItemMsg{}
	rii.Slot = domain.EquipmentSlot(pr.ReadByte())
	return rii, nil
}

type OpenContainerMsg struct {
	ContainerID   uint8
	ContainerItem domain.Item
	ContainerName string
	Capacity      uint8
	HasParent     bool
	Items         []domain.Item
}

func ParseOpenContainerMsg(pr *protocol.PacketReader) (*OpenContainerMsg, error) {

	cm := &OpenContainerMsg{}

	cm.ContainerID = pr.ReadByte()
	cm.ContainerItem = readItem(pr)
	cm.ContainerName = pr.ReadString()
	cm.Capacity = pr.ReadByte()
	cm.HasParent = pr.ReadBool()

	itemCount := pr.ReadByte()

	// Pre-allocate the slice to avoid resizing overhead
	cm.Items = make([]domain.Item, 0, itemCount)

	for i := 0; i < int(itemCount); i++ {
		item := readItem(pr)
		cm.Items = append(cm.Items, item)
	}

	return cm, nil
}

type RemoveContainerItemMsg struct {
	ContainerID uint8
	Slot        uint8
}

func ParseRemoveContainerItemMsg(pr *protocol.PacketReader) (*RemoveContainerItemMsg, error) {
	rcim := &RemoveContainerItemMsg{}

	rcim.ContainerID = pr.ReadByte()
	rcim.Slot = pr.ReadByte()

	return rcim, nil
}

type AddContainerItemMsg struct {
	ContainerID uint8
	Item        domain.Item
}

func ParseAddContainerItemMsg(pr *protocol.PacketReader) (*AddContainerItemMsg, error) {
	acim := &AddContainerItemMsg{}
	acim.ContainerID = pr.ReadByte()
	acim.Item = readItem(pr)

	return acim, nil
}

type UpdateContainerItemMsg struct {
	ContainerID uint8
	Slot        uint8
	Item        domain.Item
}

func ParseUpdateContainerItemMsg(pr *protocol.PacketReader) (*UpdateContainerItemMsg, error) {
	ucim := &UpdateContainerItemMsg{}

	ucim.ContainerID = pr.ReadByte()
	ucim.Slot = pr.ReadByte()
	ucim.Item = readItem(pr)

	return ucim, nil
}

type CloseContainerMsg struct {
	ContainerID uint8
}

func ParseCloseContainerMsg(pr *protocol.PacketReader) (*CloseContainerMsg, error) {
	ccm := &CloseContainerMsg{}

	ccm.ContainerID = pr.ReadByte()

	return ccm, nil
}

type UpdateTileItemMsg struct {
	Position domain.Position
	Stackpos uint8
	Item     domain.Item
}

func ParseUpdateTileItemMsg(pr *protocol.PacketReader) (*UpdateTileItemMsg, error) {

	utim := &UpdateTileItemMsg{}

	utim.Position = readPosition(pr)
	utim.Stackpos = pr.ReadByte()
	utim.Item = readItem(pr)

	// TODO - Item can be creature !!

	return utim, nil
}

type PlayerSkillsMsg struct {
	Skills [domain.SkillLast + 1]domain.Skill
}

func ParsePlayerSkillMsg(pr *protocol.PacketReader) (*PlayerSkillsMsg, error) {
	psm := &PlayerSkillsMsg{}

	for i := domain.SkillFirst; i <= domain.SkillLast; i++ {
		psm.Skills[i] = domain.Skill{
			Level:   pr.ReadByte(),
			Percent: pr.ReadByte(),
		}
	}
	return psm, nil
}

type PlayerStatsMsg struct {
	Health            uint16
	MaxHealth         uint16
	FreeCapacity      uint16
	Experience        uint32
	Level             uint16
	LevelPercent      uint8
	Mana              uint16
	MaxMana           uint16
	MagicLevel        uint8
	MagicLevelPercent uint8
	Soul              uint8
}

func ParsePlayerStatsMsg(pr *protocol.PacketReader) (*PlayerStatsMsg, error) {
	psm := &PlayerStatsMsg{}

	psm.Health = pr.ReadUint16()
	psm.MaxHealth = pr.ReadUint16()
	psm.FreeCapacity = pr.ReadUint16()
	psm.Experience = pr.ReadUint32()
	psm.Level = pr.ReadUint16()
	psm.LevelPercent = pr.ReadByte()
	psm.Mana = pr.ReadUint16()
	psm.MaxMana = pr.ReadUint16()
	psm.MagicLevel = pr.ReadByte()
	psm.MagicLevelPercent = pr.ReadByte()
	psm.Soul = pr.ReadByte()

	return psm, nil
}

type LoginQueueMsg struct {
	Message          string
	RetryTimeSeconds uint8
}

func (lqm *LoginQueueMsg) Encode(pw *protocol.PacketWriter) {
	pw.WriteByte(S2CSLoginQueue)
	pw.WriteString(lqm.Message)
	pw.WriteByte(lqm.RetryTimeSeconds)
}

func ParseLoginQueueMsg(pr *protocol.PacketReader) (*LoginQueueMsg, error) {
	lqm := &LoginQueueMsg{}
	lqm.Message = pr.ReadString()
	lqm.RetryTimeSeconds = pr.ReadByte()
	return lqm, pr.Err()
}
