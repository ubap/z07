package packets

import (
	"goTibia/internal/assets"
	"goTibia/internal/game/domain"
	"goTibia/internal/protocol"
)

func writePosition(pw *protocol.PacketWriter, position domain.Coordinate) {
	pw.WriteUint16(position.X)
	pw.WriteUint16(position.Y)
	pw.WriteByte(position.Z)
}

func readPosition(pr *protocol.PacketReader) domain.Coordinate {
	return domain.Coordinate{
		X: pr.ReadUint16(),
		Y: pr.ReadUint16(),
		Z: pr.ReadByte(),
	}
}

func readItem(pr *protocol.PacketReader) domain.Item {
	id := pr.ReadUint16()
	item := domain.Item{ID: id}
	thing := assets.Get(id)

	if thing.IsStackable || thing.IsFluid {
		item.Count = pr.ReadByte()
		item.HasCount = true
	}

	return item
}
