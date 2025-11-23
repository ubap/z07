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

	lr.PlayerId = pr.ReadUint32()
	lr.BeatDuration = pr.ReadUint16()
	lr.CanReportBugs = pr.ReadBool()

	return lr, pr.Err()
}
