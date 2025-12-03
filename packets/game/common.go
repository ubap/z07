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

// readItem reads a full Item (ID + Optional Count/Subtype) from the stream.
// This is used for Inventory, Containers, and Tile Stacks.
func readItem(pr *protocol.PacketReader) types.Item {
	// 1. Read ID
	id := pr.ReadUint16()

	// 2. Setup Struct
	item := types.Item{ID: id}

	thing := dat.Get(id)

	if thing.IsStackable || thing.IsFluid {
		item.Count = pr.ReadByte()
		item.HasCount = true
	}

	return item
}
