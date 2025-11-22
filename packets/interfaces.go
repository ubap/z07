package packets

import "goTibia/protocol"

type OutgoingPacket interface {
	Encode(pw *protocol.PacketWriter)
}

type IncomingPacket interface {
	Decode(pr *protocol.PacketReader) error
}
