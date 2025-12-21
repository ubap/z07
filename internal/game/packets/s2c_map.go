package packets

import (
	"fmt"
	"goTibia/internal/game/domain"
	"goTibia/internal/protocol"
)

const (
	MapWidth  = 18
	MapHeight = 14

	TileDataCreatureKnown   = 0x62 // 97
	TileDataCreatureUnknown = 0x61 // 98
	TileDataTurnCreature    = 0x63 // 99
)

type MapDescriptionMsg struct {
	PlayerPos domain.Position
	Tiles     []domain.Tile
}

func ParseMove(pr *protocol.PacketReader, ctx ParsingContext, direction domain.Direction) (*MapDescriptionMsg, error) {
	msg := &MapDescriptionMsg{
		PlayerPos: ctx.PlayerPosition,
	}

	width := MapWidth
	height := MapHeight

	switch direction {
	case domain.North:
		msg.PlayerPos.Y = ctx.PlayerPosition.Y - 1
		height = 1
	case domain.South:
		msg.PlayerPos.Y = ctx.PlayerPosition.Y + 1
		height = 1
	case domain.West:
		msg.PlayerPos.X = ctx.PlayerPosition.X - 1
		width = 1
	case domain.East:
		msg.PlayerPos.X = ctx.PlayerPosition.X + 1
		width = 1
	}

	tiles, err := parseMapDescription(pr, ctx.PlayerPosition, width, height)
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

	tiles, err := parseMapDescription(pr, msg.PlayerPos, MapWidth, MapHeight)
	if err != nil {
		return nil, err
	}
	msg.Tiles = *tiles
	return msg, err
}

func parseMapDescription(pr *protocol.PacketReader, pos domain.Position, width int, height int) (*[]domain.Tile, error) {
	tiles := make([]domain.Tile, 0, width*height)

	// 2. Determine Z-Range (Protocol 7.72 Logic)
	// If on surface (z<=7), draw from 7 down to 0.
	// If underground (z>7), draw from z-2 to z+2.
	var startZ, endZ, zStep int
	if pos.Z > 7 {
		startZ = int(pos.Z) - 2
		endZ = int(pos.Z) + 2
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
		// Calculate perspective offset for this floor
		// Tibia shifts the view when looking at lower floors
		offsetZ := int(pos.Z) - currentZ

		// fmt.Printf("Floor Z=%d, Processed=%d/%d\n", currentZ, tilesProcessed, tilesPerFloor)

		// Peek the Token
		// >= 0xFF00 means SKIP
		// <  0xFF00 means TILE (and this is the Ground ID)
		val, err := pr.PeekUint16()
		if err != nil {
			return nil, fmt.Errorf("EOF peeking token at Floor Z=%d, TileIndex=%d", currentZ, tilesProcessed)
		}

		if val >= 0xFF00 {
			// --- SKIP (RLE) ---
			_ = pr.ReadUint16()            // Consume the peeked value
			skipCount := int(val&0xFF) + 1 // Lower byte is count, skip is 0 based counter

			// fmt.Printf("  [SKIP] Count: %d (Token: %04X) | Total Processed: %d\n", skipCount, val, tilesProcessed)

			tilesProcessed += skipCount
		} else {
			// --- REAL TILE ---
			// 'val' is the Ground ID

			// Calculate actual X,Y based on linear index
			// Tibia loop order: for(x) { for(y) }
			nx := tilesProcessed / height
			ny := tilesProcessed % height

			tilePos := domain.Position{
				X: uint16(int(pos.X) + nx + offsetZ),
				Y: uint16(int(pos.Y) + ny + offsetZ),
				Z: uint8(currentZ),
			}

			// fmt.Printf("  [TILE] Parsing %v (Ground ID: %d)...\n", tilePos, val)

			// Parse the items on this tile
			tile := parseTile(pr, tilePos)
			tiles = append(tiles, tile)
		}

		// Check for Floor End
		// Use a loop because a massive skip (255) could theoretically
		// skip the remainder of Floor A, all of Floor B (unlikely 18x14=252), and start of Floor C.
		for tilesProcessed >= tilesPerFloor {

			// 1. Check if we are done with the entire volume
			if (zStep > 0 && currentZ == endZ) || (zStep < 0 && currentZ == endZ) {
				// We finished the last floor.
				// Any remaining 'skip' count is irrelevant (padding).
				// fmt.Println("--- END MAP DEBUG (Success) ---")

				return &tiles, nil
			}

			// 2. Move to next floor
			tilesProcessed -= tilesPerFloor
			currentZ += zStep
			// fmt.Printf("--- MOVING TO FLOOR Z=%d ---\n", currentZ)
		}
	}

}

func parseTile(pr *protocol.PacketReader, pos domain.Position) domain.Tile {
	// 1. Setup the Tile struct
	t := domain.Tile{
		Position: pos,
		Items:    make([]domain.Item, 0, 4), // Pre-allocate small cap for performance
	}

	groundItem := readItem(pr)
	t.Ground = &groundItem

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
			// We must consume the bytes to keep the stream aligned.

			// Note: For a pure MapDescription parser, we often discard creature data
			// because creatures are usually tracked via 0x6A/0x6B packets.
			// However, we MUST read it to advance the cursor.

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
	if marker == TileDataCreatureKnown { // 0x62
		// C++: if (known) msg.add<uint32_t>(creature->getID());
		_ = pr.ReadUint32() // ID

	} else if marker == TileDataCreatureUnknown { // 0x61
		_ = pr.ReadUint32() // The id to remove from knowns, it is there to free some slot from known creatures list.
		_ = pr.ReadUint32() // ID
		_ = pr.ReadString() // Name
	} else {
		return fmt.Errorf("unknown creature marker: 0x%X", marker)
	}

	// Health
	_ = pr.ReadByte()

	// Direction
	_ = pr.ReadByte()

	// Outfit
	if err := readOutfit(pr); err != nil {
		return err
	}

	// Light
	_ = pr.ReadByte() // Light Level
	_ = pr.ReadByte() // Light Color

	// Speed
	_ = pr.ReadUint16()

	// Skull & Party
	_ = pr.ReadByte() // Skull
	_ = pr.ReadByte() // Party Shield

	return nil
}

func readOutfit(pr *protocol.PacketReader) error {
	lookType := pr.ReadUint16()

	if lookType != 0 {
		_ = pr.ReadByte() // Head
		_ = pr.ReadByte() // Body
		_ = pr.ReadByte() // Legs
		_ = pr.ReadByte() // Feet
	} else {
		// Item Outfit (Chameleon Rune, etc.)
		_ = pr.ReadUint16() // Look Item ID
	}
	return nil
}
