package packets

import (
	"errors"
	"goTibia/protocol"
	"goTibia/protocol/crypto"
)

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

func (lr *LoginRequest) Encode(pw *protocol.PacketWriter) {
	pw.WriteByte(lr.Protocol)
	pw.WriteUint16(lr.ClientOS)
	pw.WriteUint16(lr.ClientVersion)

	// RSA Encrypted part starts here
	toEncrypt := protocol.NewPacketWriter()

	toEncrypt.WriteByte(0x00) // Write the check check byte
	toEncrypt.WriteUint32(lr.XTEAKey[0])
	toEncrypt.WriteUint32(lr.XTEAKey[1])
	toEncrypt.WriteUint32(lr.XTEAKey[2])
	toEncrypt.WriteUint32(lr.XTEAKey[3])
	toEncrypt.WriteBool(lr.Gamemaster)
	toEncrypt.WriteUint32(lr.AccountNumber)
	toEncrypt.WriteString(lr.CharacterName)
	toEncrypt.WriteString(lr.Password)

	// Encrypt the data block with the target server's public key.
	unencodedBytes, err := toEncrypt.GetBytes()
	pw.SetError(err)

	encryptedBlock, err := crypto.EncryptRSA(crypto.RSA.GameServerPublicKey, unencodedBytes)
	pw.SetError(err)

	pw.WriteBytes(encryptedBlock)
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

func (lr *LoginRequest) GetXTEAKey() [4]uint32 {
	return lr.XTEAKey
}
