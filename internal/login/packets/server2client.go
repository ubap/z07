package packets

import (
	"errors"
	"fmt"
	protocol2 "goTibia/internal/protocol"
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

func (lp *LoginResultMessage) Encode(pw *protocol2.PacketWriter) {
	if lp.ClientDisconnected {
		pw.WriteByte(S2COpcodeDisconnectClient)
		pw.WriteString(lp.ClientDisconnectedReason)
	}

	if lp.Motd != nil {
		pw.WriteByte(S2COpcodeMOTD)
		motdData := fmt.Sprintf("%s\n%s", lp.Motd.MotdId, lp.Motd.Message)
		pw.WriteString(motdData)
	}

	if lp.CharacterList != nil {
		pw.WriteByte(S2COpcodeCharacterList)
		writeCharacterList(pw, lp.CharacterList)
	}
}

func ParseLoginResultMessage(pr *protocol2.PacketReader) (*LoginResultMessage, error) {
	message := LoginResultMessage{}
	for pr.Remaining() > 0 {
		opcode := pr.ReadByte()
		if err := pr.Err(); err != nil {
			return nil, fmt.Errorf("failed to read opcode: %w", err)
		}

		switch opcode {
		case S2COpcodeDisconnectClient:
			disconnectedReason := pr.ReadString()
			if err := pr.Err(); err != nil {
				return nil, fmt.Errorf("failed to read disconnect reason: %w", err)
			}
			message.ClientDisconnected = true
			message.ClientDisconnectedReason = disconnectedReason

		case S2COpcodeMOTD:
			motd, err := parseMotd(pr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse MOTD: %w", err)
			}
			message.Motd = motd

		case S2COpcodeCharacterList:
			charList, err := parseCharacterList(pr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse CharList: %w", err)
			}
			message.CharacterList = charList

		default:
			return nil, fmt.Errorf("unknown login opcode: %#x", opcode)
		}

		if err := pr.Err(); err != nil {
			return nil, fmt.Errorf("error parsing packet content for opcode %#x: %w", opcode, err)
		}

	}
	return &message, nil
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

func parseCharacterList(packetReader *protocol2.PacketReader) (*CharacterList, error) {
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

func writeCharacterEntry(pw *protocol2.PacketWriter, entry *CharacterEntry) {
	pw.WriteString(entry.Name)
	pw.WriteString(entry.WorldName)
	pw.WriteUint32(entry.WorldIp)
	pw.WriteUint16(entry.WorldPort)
}

func writeCharacterList(pw *protocol2.PacketWriter, charList *CharacterList) {
	pw.WriteByte(uint8(len(charList.Characters)))
	for _, entry := range charList.Characters {
		writeCharacterEntry(pw, entry)
	}
	pw.WriteUint16(charList.PremiumDays)
}

// endregion CharacterList

// region MOTD

type Motd struct {
	MotdId  string
	Message string
}

func parseMotd(packetReader *protocol2.PacketReader) (*Motd, error) {
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
