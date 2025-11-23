package game

import (
	"errors"
	"goTibia/protocol"
	"goTibia/protocol/crypto"
)

// LoginRequest is a special packet. It's the first packet sent by the client to the server
type LoginRequest struct {
	Protocol      uint8
	ClientOS      uint16
	ClientVersion uint16

	XTEAKey       [4]uint32
	Gamemaster    bool
	AccountNumber uint32
	CharacterName string
	Password      string
}

func ParseLoginRequest(packetReader *protocol.PacketReader) (*LoginRequest, error) {
	packet := &LoginRequest{}

	packet.Protocol = packetReader.ReadByte()
	packet.ClientOS = packetReader.ReadUint16()
	packet.ClientVersion = packetReader.ReadUint16()

	encryptedBlock := packetReader.ReadAll()
	if packetReader.Err() != nil {
		return nil, packetReader.Err()
	}

	decryptedBlock := crypto.DecryptRSA(encryptedBlock)
	decryptedBlockReader := protocol.NewPacketReader(decryptedBlock)
	checkByte := decryptedBlockReader.ReadByte()
	if checkByte != 0x00 {
		return nil, errors.New("invalid checkByte")
	}

	packet.XTEAKey[0] = decryptedBlockReader.ReadUint32()
	packet.XTEAKey[1] = decryptedBlockReader.ReadUint32()
	packet.XTEAKey[2] = decryptedBlockReader.ReadUint32()
	packet.XTEAKey[3] = decryptedBlockReader.ReadUint32()
	packet.Gamemaster = decryptedBlockReader.ReadBool()
	packet.AccountNumber = decryptedBlockReader.ReadUint32()
	packet.CharacterName = decryptedBlockReader.ReadString()
	packet.Password = decryptedBlockReader.ReadString()

	return packet, packetReader.Err()
}
