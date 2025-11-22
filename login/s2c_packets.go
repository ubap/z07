package login

import (
	"encoding/binary"
	"errors"
	"goTibia/protocol"
	"io"
	"log"
	"strconv"
	"strings"
)

// region LoginResultMessage

type LoginResultMessage struct {
	ClientDisconnected       bool
	ClientDisconnectedReason string
	Motd                     *Motd
	CharacterList            *CharacterList
}

func (lp *LoginResultMessage) Marshal() ([]byte, error) {
	pw := protocol.NewPacketWriter()

	if lp.ClientDisconnected {
		pw.WriteByte(S2COpcodeDisconnectClient)
		pw.WriteString(lp.ClientDisconnectedReason)
	}

	if lp.CharacterList != nil {
		pw.WriteByte(S2COpcodeCharacterList)
		WriteCharacterList(pw, lp.CharacterList)
	}

	return pw.GetBytes()
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

func ReadCharacterList(r io.Reader) (*CharacterList, error) {
	entryCount, err := protocol.ReadByte(r)
	if err != nil {
		return nil, err
	}

	var characterEntries []*CharacterEntry
	for i := 0; i < int(entryCount); i++ {
		name, err := protocol.ReadString(r)
		if err != nil {
			return nil, err
		}

		worldName, err := protocol.ReadString(r)
		if err != nil {
			return nil, err
		}

		var worldIp uint32
		if err := binary.Read(r, binary.LittleEndian, &worldIp); err != nil {
			return nil, err
		}

		var worldPort uint16
		if err := binary.Read(r, binary.LittleEndian, &worldPort); err != nil {
			return nil, err
		}

		characterEntries = append(characterEntries, &CharacterEntry{Name: name, WorldName: worldName, WorldIp: worldIp, WorldPort: worldPort})
	}

	var premiumDays uint16
	if err := binary.Read(r, binary.LittleEndian, &premiumDays); err != nil {
		return nil, err
	}

	return &CharacterList{Characters: characterEntries, PremiumDays: premiumDays}, nil
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
	MotdId  int
	Message string
}

func ReadMotd(r io.Reader) (*Motd, error) {
	data, err := protocol.ReadString(r)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(data, "\n", 2)

	if len(parts) != 2 {
		return nil, errors.New("invalid format")
	}

	motdId, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errors.New("failed to parse MOTD ID")
	}

	message := parts[1]

	log.Printf("MOTDID :%d, MOTD: %s", motdId, message)
	return &Motd{MotdId: motdId, Message: message}, nil
}

// endregion MOTD
