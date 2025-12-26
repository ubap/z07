package packets

import (
	"fmt"
	"goTibia/internal/game/domain"
	"goTibia/internal/protocol"
)

const (
	ClientViewportX = 8
	ClientViewportY = 6

	TileDataCreatureKnown   = 0x62 // 97
	TileDataCreatureUnknown = 0x61 // 98
	TileDataTurnCreature    = 0x63 // 99
)

type MapDescriptionMsg struct {
	PlayerPos domain.Position
	Tiles     map[domain.Position]domain.Tile
}

func ParseMove(pr *protocol.PacketReader, ctx ParsingContext, direction domain.Direction) (*MapDescriptionMsg, error) {
	msg := &MapDescriptionMsg{
		PlayerPos: ctx.PlayerPosition,
	}

	var x int
	var y int
	var z = int(msg.PlayerPos.Z)
	var width = ClientViewportX*2 + 2
	var height = ClientViewportY*2 + 2

	switch direction {
	case domain.North:
		msg.PlayerPos.Y = ctx.PlayerPosition.Y - 1
		x = int(msg.PlayerPos.X) - ClientViewportX
		y = int(msg.PlayerPos.Y) - ClientViewportY
		height = 1
	case domain.South:
		msg.PlayerPos.Y = ctx.PlayerPosition.Y + 1
		x = int(msg.PlayerPos.X) - ClientViewportX
		y = int(msg.PlayerPos.Y) + ClientViewportY + 1
		height = 1
	case domain.West:
		msg.PlayerPos.X = ctx.PlayerPosition.X - 1
		x = int(msg.PlayerPos.X) - ClientViewportX
		y = int(msg.PlayerPos.Y) - ClientViewportY
		width = 1
	case domain.East:
		msg.PlayerPos.X = ctx.PlayerPosition.X + 1
		x = int(msg.PlayerPos.X) + ClientViewportX + 1
		y = int(msg.PlayerPos.Y) - ClientViewportY
		width = 1
	}

	tiles, err := parseMapDescription(pr, x, y, z, width, height)
	if err != nil {
		return nil, err
	}
	msg.Tiles = *tiles
	return msg, nil
}

func ParseMapDescriptionMsg(pr *protocol.PacketReader) (*MapDescriptionMsg, error) {
	msg := &MapDescriptionMsg{
		PlayerPos: readPosition(pr),
	}

	var x = int(msg.PlayerPos.X) - ClientViewportX
	var y = int(msg.PlayerPos.Y) - ClientViewportY
	var z = int(msg.PlayerPos.Z)

	tiles, err := parseMapDescription(pr, x, y, z, ClientViewportX*2+2, ClientViewportY*2+2)
	if err != nil {
		return nil, err
	}
	msg.Tiles = *tiles
	return msg, err
}

func parseMapDescription(pr *protocol.PacketReader, x, y, z, width, height int) (*map[domain.Position]domain.Tile, error) {
	tiles := make(map[domain.Position]domain.Tile)

	// 2. Determine Z-Range
	// If on surface (z<=7), draw from 7 down to 0.
	// If underground (z>7), draw from z-2 to z+2.
	var startZ, endZ, zStep int
	if z > 7 {
		startZ = max(z-2, 0)
		endZ = z + 2 // TODO: Clamp to Max?
		zStep = 1
	} else {
		startZ = 7
		endZ = 0
		zStep = -1
	}

	currentZ := startZ
	tilesProcessed := 0
	tilesPerFloor := width * height

	// 3. Loop until all floors are read
	for {
		// Peek the Token
		// >= 0xFF00 means SKIP
		// <  0xFF00 means TILE (and this is the Ground ID)
		val, err := pr.PeekUint16()
		if err != nil {
			return nil, fmt.Errorf("EOF peeking token at Floor Z=%d, TileIndex=%d", currentZ, tilesProcessed)
		}

		if val >= 0xFF00 {
			// --- SKIP (Run-Length Encoding) ---
			_ = pr.ReadUint16()            // Consume the peeked value
			skipCount := int(val&0xFF) + 1 // Lower byte is count, skip is 0 based counter
			tilesProcessed += skipCount
		} else {
			// --- REAL TILE ---
			// 'val' is the Ground ID

			// Calculate perspective offset for this floor
			// Tibia shifts the view when looking at lower floors
			offsetZ := z - currentZ

			// Calculate actual X,Y based on linear index
			// Tibia loop order: for(x) { for(y) }
			nx := tilesProcessed / height
			ny := tilesProcessed % height

			tilePos := domain.Position{
				X: uint16(x + nx + offsetZ),
				Y: uint16(y + ny + offsetZ),
				Z: uint8(currentZ),
			}

			tile := parseTile(pr)
			tiles[tilePos] = tile
		}

		// skip tiles
		for tilesProcessed >= tilesPerFloor {
			// 1. Check if we are done with the entire volume
			if currentZ == endZ {
				return &tiles, nil
			}

			// 2. Move to next floor
			tilesProcessed -= tilesPerFloor
			currentZ += zStep
			// fmt.Printf("--- MOVING TO FLOOR Z=%d ---\n", currentZ)
		}
	}
}

func parseTile(pr *protocol.PacketReader) domain.Tile {
	// 1. Setup the Tile struct
	t := domain.Tile{
		Items: make([]domain.Item, 0, 4), // Pre-allocate small cap for performance
	}

	groundItem := readItem(pr)
	t.Items = append(t.Items, groundItem)

	// 3. Loop: Read Items on top of the ground
	// We read until we hit a "Skip" marker (>= 0xFF00) which belongs to the NEXT tile.
	for {
		// A. Peek at the next 2 bytes
		nextVal, err := pr.PeekUint16()

		// B. Stop conditions:
		// - Error/EOF
		// - Value is >= 0xFF00 (This is a Skip/RLE marker for the map loop)
		if err != nil || nextVal >= 0xFF00 {
			// fmt.Println("End of Tile")
			break
		}

		if nextVal == TileDataCreatureKnown || nextVal == TileDataCreatureUnknown {
			// It is a CREATURE, not an ITEM.
			err := readCreatureInMap(pr)
			if err != nil {
				// fmt.Printf("Error reading creature in map at tile %v: %v\n", pos, err)
				return domain.Tile{}
			}
			continue
		}

		item := readItem(pr)
		t.Items = append(t.Items, item)
	}

	return t
}

func readCreatureInMap(pr *protocol.PacketReader) error {
	// 1. Read Marker (We already peeked it, but we must consume it)
	marker := pr.ReadUint16()

	// 2. Handle ID / Name logic
	switch marker {
	case TileDataCreatureKnown:
		_ = pr.ReadUint32() // ID
	case TileDataCreatureUnknown:
		_ = pr.ReadUint32() // The id to remove from knowns, it is there to free some slot from known creatures list.
		_ = pr.ReadUint32() // ID
		_ = pr.ReadString() // Name
	default:
		return fmt.Errorf("unknown creature marker: 0x%X", marker)
	}

	// Health
	_ = pr.ReadUint8()

	// Direction
	_ = pr.ReadUint8()

	// Outfit
	if err := readOutfit(pr); err != nil {
		return err
	}

	// Light
	_ = pr.ReadUint8() // Light Level
	_ = pr.ReadUint8() // Light Color

	// Speed
	_ = pr.ReadUint16()

	// Skull & Party
	_ = pr.ReadUint8() // Skull
	_ = pr.ReadUint8() // Party Shield

	return nil
}

func readOutfit(pr *protocol.PacketReader) error {
	lookType := pr.ReadUint16()

	if lookType != 0 {
		_ = pr.ReadUint8() // Head
		_ = pr.ReadUint8() // Body
		_ = pr.ReadUint8() // Legs
		_ = pr.ReadUint8() // Feet
	} else {
		// Item Outfit (Chameleon Rune, etc.)
		_ = pr.ReadUint16() // Look Item ID
	}
	return nil
}
