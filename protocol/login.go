package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type LoginPacket struct {
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

func (p *LoginPacket) Write(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, p.ClientOS); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, p.ClientVersion); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, p.AccountNumber); err != nil {
		return err
	}
	// For strings, first write the length (as uint16), then the string itself.
	if err := binary.Write(w, binary.LittleEndian, uint16(len(p.Password))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(p.Password)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, p.XTEAKey); err != nil {
		return err
	}
	return nil
}

func ParseLoginPacket(data []byte) (*LoginPacket, error) {
	// Use a bytes.Reader to treat the byte slice like a file or network stream.
	// This makes it easy to read structured data sequentially.
	reader := bytes.NewReader(data)
	packet := &LoginPacket{}

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

	decryptedBlock, err := DecryptRSA(encryptedBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt rsa block: %w", err)
	}

	// --- 3. Parse the Decrypted Block ---
	decryptedReader := bytes.NewReader(decryptedBlock)

	if err := binary.Read(decryptedReader, binary.LittleEndian, &packet.XTEAKey); err != nil {
		return nil, fmt.Errorf("failed to read xtea key from decrypted block: %w", err)
	}
	if err := binary.Read(decryptedReader, binary.LittleEndian, &packet.AccountNumber); err != nil {
		return nil, fmt.Errorf("failed to read account number from decrypted block: %w", err)
	}

	password, err := readString(decryptedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse password: %w", err)
	}
	packet.Password = password

	return packet, nil
}
