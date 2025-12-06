package packets

import (
	"goTibia/internal/assets"
	"goTibia/internal/game/types"
	protocol2 "goTibia/internal/protocol"
)

func writePosition(pw *protocol2.PacketWriter, position types.Position) {
	pw.WriteUint16(position.X)
	pw.WriteUint16(position.Y)
	pw.WriteByte(position.Z)
}

func readPosition(pr *protocol2.PacketReader) types.Position {
	return types.Position{
		X: pr.ReadUint16(),
		Y: pr.ReadUint16(),
		Z: pr.ReadByte(),
	}
}

func readItem(pr *protocol2.PacketReader) types.Item {
	id := pr.ReadUint16()
	item := types.Item{ID: id}
	thing := assets.Get(id)

	if thing.IsStackable || thing.IsFluid {
		item.Count = pr.ReadByte()
		item.HasCount = true
	}

	return item
}
