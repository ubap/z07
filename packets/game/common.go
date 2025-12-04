package game

import (
	"goTibia/dat"
	"goTibia/protocol"
	"goTibia/types"
)

func writePosition(pw *protocol.PacketWriter, position types.Position) {
	pw.WriteUint16(position.X)
	pw.WriteUint16(position.Y)
	pw.WriteByte(position.Z)
}

func readPosition(pr *protocol.PacketReader) types.Position {
	return types.Position{
		X: pr.ReadUint16(),
		Y: pr.ReadUint16(),
		Z: pr.ReadByte(),
	}
}

func readItem(pr *protocol.PacketReader) types.Item {
	id := pr.ReadUint16()
	item := types.Item{ID: id}
	thing := dat.Get(id)

	if thing.IsStackable || thing.IsFluid {
		item.Count = pr.ReadByte()
		item.HasCount = true
	}

	return item
}
