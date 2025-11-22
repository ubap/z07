package packets

import (
	"errors"
	"goTibia/protocol"
	"log"
	"strings"
)

// region LoginResultMessage

type LoginResultMessage struct {
	ClientDisconnected       bool
	ClientDisconnectedReason string
	Motd                     *Motd
	CharacterList            *CharacterList
}

func (lp *LoginResultMessage) Encode(pw *protocol.PacketWriter) error {
	if lp.ClientDisconnected {
		pw.WriteByte(S2COpcodeDisconnectClient)
		pw.WriteString(lp.ClientDisconnectedReason)
	}

	if lp.CharacterList != nil {
		pw.WriteByte(S2COpcodeCharacterList)
		WriteCharacterList(pw, lp.CharacterList)
	}

	return pw.Err()
}

// endregion LoginResultMessage

// region CharacterList
type CharacterList struct {
	Characters  []*CharacterEntry
	PremiumDays uint16
}

type CharacterEntry struct {
	Name      string
	WorldName string
	WorldIp   uint32
	WorldPort uint16
}

func ParseCharacterList(packetReader *protocol.PacketReader) (*CharacterList, error) {
	entryCount := packetReader.ReadByte()

	characterEntries := make([]*CharacterEntry, entryCount)
	for i := 0; i < int(entryCount); i++ {
		name := packetReader.ReadString()
		worldName := packetReader.ReadString()
		worldIp := packetReader.ReadUint32()
		worldPort := packetReader.ReadUint16()

		characterEntries[i] = &CharacterEntry{
			Name:      name,
			WorldName: worldName,
			WorldIp:   worldIp,
			WorldPort: worldPort,
		}
	}

	premiumDays := packetReader.ReadUint16()

	return &CharacterList{Characters: characterEntries, PremiumDays: premiumDays}, packetReader.Err()
}

func WriteCharacterEntry(pw *protocol.PacketWriter, entry *CharacterEntry) {
	pw.WriteString(entry.Name)
	pw.WriteString(entry.WorldName)
	pw.WriteUint32(entry.WorldIp)
	pw.WriteUint16(entry.WorldPort)
}

func WriteCharacterList(pw *protocol.PacketWriter, charList *CharacterList) {
	pw.WriteByte(uint8(len(charList.Characters)))
	for _, entry := range charList.Characters {
		WriteCharacterEntry(pw, entry)
	}
	pw.WriteUint16(charList.PremiumDays)
}

// endregion CharacterList

// region MOTD

type Motd struct {
	MotdId  string
	Message string
}

func ParseMotd(packetReader *protocol.PacketReader) (*Motd, error) {
	data := packetReader.ReadString()
	if packetReader.Err() != nil {
		return nil, packetReader.Err()
	}

	parts := strings.SplitN(data, "\n", 2)

	if len(parts) != 2 {
		return nil, errors.New("invalid format")
	}

	message := parts[1]
	motd := &Motd{MotdId: parts[0], Message: parts[1]}

	log.Printf("motd: ID: %s, MOTD: %s", motd.MotdId, message)
	return motd, nil
}

// endregion MOTD
