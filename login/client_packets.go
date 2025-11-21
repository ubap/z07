package login

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"goTibia/protocol"
	"goTibia/protocol/crypto"
	"io"
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

func ParseCredentialsPacket(data []byte) (*ClientCredentialPacket, error) {
	// Use a bytes.Reader to treat the byte slice like a file or network stream.
	// This makes it easy to read structured data sequentially.
	reader := bytes.NewReader(data)
	packet := &ClientCredentialPacket{}

	// --- 1. Read the Unencrypted Header ---
	// We read each field in order, using binary.Read for fixed-size data.
	// The byte order must be LittleEndian, which is standard for this protocol.
	if err := binary.Read(reader, binary.LittleEndian, &packet.Protocol); err != nil {
		return nil, fmt.Errorf("failed to read protocol version: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &packet.ClientOS); err != nil {
		return nil, fmt.Errorf("failed to read client os: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &packet.ClientVersion); err != nil {
		return nil, fmt.Errorf("failed to read client version: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &packet.DatSignature); err != nil {
		return nil, fmt.Errorf("failed to read dat signature: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &packet.SprSignature); err != nil {
		return nil, fmt.Errorf("failed to read spr signature: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &packet.PicSignature); err != nil {
		return nil, fmt.Errorf("failed to read pic signature: %w", err)
	}

	// --- 2. Decrypt the Remainder of the Packet ---
	encryptedBlock, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted block: %w", err)
	}

	decryptedBlock := crypto.DecryptRSA(encryptedBlock)

	messagePayload := bytes.TrimLeft(decryptedBlock, "\x00")

	// Now, create a reader from the clean, left-aligned payload.
	decryptedReader := bytes.NewReader(messagePayload)

	if err := binary.Read(decryptedReader, binary.LittleEndian, &packet.XTEAKey); err != nil {
		return nil, fmt.Errorf("failed to read xtea key from decrypted block: %w", err)
	}
	if err := binary.Read(decryptedReader, binary.LittleEndian, &packet.AccountNumber); err != nil {
		return nil, fmt.Errorf("failed to read account number from decrypted block: %w", err)
	}

	password, err := protocol.ReadString(decryptedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse password: %w", err)
	}
	packet.Password = password

	return packet, nil
}
