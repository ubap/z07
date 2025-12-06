package packets

import (
	"errors"
	protocol2 "goTibia/internal/protocol"
	"goTibia/internal/protocol/crypto"
)

// ClientCredentialPacket is a special packet. It's the first packet sent by the client to the server
type ClientCredentialPacket struct {
	Protocol      uint8
	ClientOS      uint16
	ClientVersion uint16
	DatSignature  uint32
	SprSignature  uint32
	PicSignature  uint32
	XTEAKey       [4]uint32
	AccountNumber uint32
	Password      string
}

func (lp *ClientCredentialPacket) Encode(pw *protocol2.PacketWriter) {
	pw.WriteByte(lp.Protocol)
	pw.WriteUint16(lp.ClientOS)
	pw.WriteUint16(lp.ClientVersion)
	pw.WriteUint32(lp.DatSignature)
	pw.WriteUint32(lp.SprSignature)
	pw.WriteUint32(lp.PicSignature)

	// RSA Encrypted part starts here
	toEncrypt := protocol2.NewPacketWriter()

	toEncrypt.WriteByte(0x00) // Write the check check byte
	toEncrypt.WriteUint32(lp.XTEAKey[0])
	toEncrypt.WriteUint32(lp.XTEAKey[1])
	toEncrypt.WriteUint32(lp.XTEAKey[2])
	toEncrypt.WriteUint32(lp.XTEAKey[3])
	toEncrypt.WriteUint32(lp.AccountNumber)
	toEncrypt.WriteString(lp.Password)

	// Encrypt the data block with the target server's public key.
	unencodedBytes, err := toEncrypt.GetBytes()
	pw.SetError(err)

	encryptedBlock, err := crypto.EncryptRSA(crypto.RSA.GameServerPublicKey, unencodedBytes)
	pw.SetError(err)

	pw.WriteBytes(encryptedBlock)
}

func ParseCredentialsPacket(packetReader *protocol2.PacketReader) (*ClientCredentialPacket, error) {
	packet := &ClientCredentialPacket{}

	packet.Protocol = packetReader.ReadByte()
	packet.ClientOS = packetReader.ReadUint16()
	packet.ClientVersion = packetReader.ReadUint16()
	packet.DatSignature = packetReader.ReadUint32()
	packet.SprSignature = packetReader.ReadUint32()
	packet.PicSignature = packetReader.ReadUint32()

	encryptedBlock := packetReader.ReadAll()
	if packetReader.Err() != nil {
		return nil, packetReader.Err()
	}

	decryptedBlock := crypto.DecryptRSA(encryptedBlock)
	decryptedBlockReader := protocol2.NewPacketReader(decryptedBlock)
	checkByte := decryptedBlockReader.ReadByte()
	if checkByte != 0x00 {
		return nil, errors.New("invalid checkByte")
	}

	packet.XTEAKey[0] = decryptedBlockReader.ReadUint32()
	packet.XTEAKey[1] = decryptedBlockReader.ReadUint32()
	packet.XTEAKey[2] = decryptedBlockReader.ReadUint32()
	packet.XTEAKey[3] = decryptedBlockReader.ReadUint32()
	packet.AccountNumber = decryptedBlockReader.ReadUint32()
	packet.Password = decryptedBlockReader.ReadString()

	return packet, packetReader.Err()
}

func (p *ClientCredentialPacket) GetXTEAKey() [4]uint32 {
	return p.XTEAKey
}
