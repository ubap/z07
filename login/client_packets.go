package login

import (
	"bytes"
	"encoding/binary"
	"errors"
	"goTibia/protocol"
	"goTibia/protocol/crypto"
)

// Special packet - first received from the client during login.
type ClientCredentialPacket struct {
	Protocol      uint8
	ClientOS      uint16
	ClientVersion uint16
	DatSignature  uint32
	SprSignature  uint32
	PicSignature  uint32
	// RSA Encrypted part starts here
	XTEAKey       [4]uint32
	AccountNumber uint32
	Password      string
}

func (lp *ClientCredentialPacket) Marshal() ([]byte, error) {
	// 1. Prepare the plaintext block that needs to be encrypted.
	rsaPlaintext := new(bytes.Buffer)

	// Write the check byte, XTEA key, account number, and password.
	rsaPlaintext.WriteByte(0x00)
	binary.Write(rsaPlaintext, binary.LittleEndian, lp.XTEAKey)
	binary.Write(rsaPlaintext, binary.LittleEndian, lp.AccountNumber)

	// Write the length-prefixed password string.
	binary.Write(rsaPlaintext, binary.LittleEndian, uint16(len(lp.Password)))
	rsaPlaintext.WriteString(lp.Password)

	// 2. Encrypt the plaintext block with the target server's public key.
	encryptedBlock, err := crypto.EncryptRSA(crypto.RSA.GameServerPublicKey, rsaPlaintext.Bytes())
	if err != nil {
		return nil, err
	}

	// 3. Assemble the full packet by prepending the unencrypted header.
	fullPacket := new(bytes.Buffer)
	binary.Write(fullPacket, binary.LittleEndian, lp.Protocol)
	binary.Write(fullPacket, binary.LittleEndian, lp.ClientOS)
	binary.Write(fullPacket, binary.LittleEndian, lp.ClientVersion)
	binary.Write(fullPacket, binary.LittleEndian, lp.DatSignature)
	binary.Write(fullPacket, binary.LittleEndian, lp.SprSignature)
	binary.Write(fullPacket, binary.LittleEndian, lp.PicSignature)
	fullPacket.Write(encryptedBlock)

	return fullPacket.Bytes(), nil
}

func ParseCredentialsPacket(packetReader *protocol.PacketReader) (*ClientCredentialPacket, error) {
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
	decryptedBlockReader := protocol.NewPacketReader(decryptedBlock)
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
