package main

import (
	"encoding/binary"
	"io"
)

type LoginPacket struct {
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
