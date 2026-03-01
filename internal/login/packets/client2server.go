package packets

import (
	"errors"
	"z07/internal/protocol"
	"z07/internal/protocol/crypto"
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

func NewClientCredentialPacket(accountNumber uint32, password string) *ClientCredentialPacket {
	return &ClientCredentialPacket{
		Protocol:      1,
		ClientOS:      1,
		ClientVersion: 772,
		DatSignature:  1,
		SprSignature:  2,
		PicSignature:  3,
		XTEAKey:       [4]uint32{0x11111111, 0x22222222, 0x33333333, 0x44444444},
		AccountNumber: accountNumber,
		Password:      password,
	}
}

func (lp *ClientCredentialPacket) Encode(pw *protocol.PacketWriter) {
	pw.WriteUint8(lp.Protocol)
	pw.WriteUint16(lp.ClientOS)
	pw.WriteUint16(lp.ClientVersion)
	pw.WriteUint32(lp.DatSignature)
	pw.WriteUint32(lp.SprSignature)
	pw.WriteUint32(lp.PicSignature)

	// RSA Encrypted part starts here
	toEncrypt := protocol.NewPacketWriter()

	toEncrypt.WriteUint8(0x00) // Write the check check byte
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

func ParseCredentialsPacket(packetReader *protocol.PacketReader) (*ClientCredentialPacket, error) {
	packet := &ClientCredentialPacket{}

	packet.Protocol = packetReader.ReadUint8()
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
	decryptedBlockReader := protocol.NewPacketReader(decryptedBlock)
	checkByte := decryptedBlockReader.ReadUint8()
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
