package game

import (
	"goTibia/protocol"
)

type LoginResponse struct {
	ClientDisconnected bool
	PlayerId           uint32
	BeatDuration       uint16
	CanReportBugs      bool
}

func ParseLoginResultMessage(pr *protocol.PacketReader) (*LoginResponse, error) {
	lr := &LoginResponse{}

	// first byte is 0x0A for success, errors not handled yet.
	success := pr.ReadByte()
	if success != 0x0A { // Regular opcode opcode
		lr.ClientDisconnected = true
		return lr, nil
	}

	lr.PlayerId = pr.ReadUint32()
	lr.BeatDuration = pr.ReadUint16()
	lr.CanReportBugs = pr.ReadBool()

	return lr, pr.Err()
}
