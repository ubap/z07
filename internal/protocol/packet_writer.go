package protocol

import (
	"bytes"
	"encoding/binary"
)

type PacketWriter struct {
	buff *bytes.Buffer
	err  error
}

func NewPacketWriter() *PacketWriter {
	p := &PacketWriter{buff: new(bytes.Buffer)}
	return p
}

func (pw *PacketWriter) Err() error {
	return pw.err
}

func (pw *PacketWriter) WriteByte(b byte) {
	if pw.err != nil {
		return
	}
	// Use native WriteByte to avoid allocation
	pw.err = pw.buff.WriteByte(b)
}

func (pw *PacketWriter) WriteUint16(val uint16) {
	if pw.err != nil {
		return
	}
	// Optimization: Avoid binary.Write (reflection).
	// Write 2 raw bytes directly.
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], val)
	_, pw.err = pw.buff.Write(b[:])
}

func (pw *PacketWriter) WriteUint32(val uint32) {
	if pw.err != nil {
		return
	}
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], val)
	_, pw.err = pw.buff.Write(b[:])
}

func (pw *PacketWriter) WriteString(s string) {
	if pw.err != nil {
		return
	}
	// Write string length
	pw.WriteUint16(uint16(len(s)))
	if pw.err != nil {
		return
	}
	// Use WriteString to avoid casting string to []byte
	_, pw.err = pw.buff.WriteString(s)
}

func (pw *PacketWriter) WriteBytes(data []byte) {
	if pw.err != nil {
		return
	}
	if len(data) > 0 {
		_, pw.err = pw.buff.Write(data)
	}
}

func (pw *PacketWriter) WriteBool(data bool) {
	if pw.err != nil {
		return
	}
	if data {
		pw.err = pw.buff.WriteByte(1)
	} else {
		pw.err = pw.buff.WriteByte(0)
	}
}

func (pw *PacketWriter) GetBytes() ([]byte, error) {
	if pw.err != nil {
		return nil, pw.err
	}

	finalBytes := pw.buff.Bytes()

	return finalBytes, nil
}

func (pw *PacketWriter) SetError(err error) {
	pw.err = err
}
