package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Configuration
const (
	InputPath       = "data/772/Tibia.dat"
	OutputPath      = "data/772/items.json"
	DatSignature772 = 0x439D5A33
)

// ItemAttributes represents the JSON structure for your bot.
type ItemAttributes struct {
	ID   uint16 `json:"id"`
	Name string `json:"name,omitempty"` // Placeholder for manual editing

	// --- Logic Flags ---
	IsGround       bool   `json:"is_ground,omitempty"`
	Speed          uint16 `json:"speed,omitempty"`
	IsBlocking     bool   `json:"is_blocking,omitempty"`      // Solids (Walls)
	IsMissileBlock bool   `json:"is_missile_block,omitempty"` // Blocks Projectiles
	IsPathBlock    bool   `json:"is_path_block,omitempty"`    // Unpassable (Magic Walls)

	IsContainer  bool `json:"is_container,omitempty"`
	IsStackable  bool `json:"is_stackable,omitempty"`
	IsFluid      bool `json:"is_fluid,omitempty"`
	IsMultiUse   bool `json:"is_multi_use,omitempty"` // Runes, fluids
	IsPickupable bool `json:"is_pickupable,omitempty"`
	IsRotatable  bool `json:"is_rotatable,omitempty"`

	// --- Visuals (Useful for bot context) ---
	LightLevel   uint16 `json:"light_level,omitempty"`
	LightColor   uint16 `json:"light_color,omitempty"`
	Elevation    uint16 `json:"elevation,omitempty"`
	MinimapColor uint16 `json:"minimap_color,omitempty"`
}

func main() {
	fmt.Printf("Reading %s...\n", InputPath)

	data, err := os.ReadFile(InputPath)
	if err != nil {
		log.Fatalf("Error reading .dat file: %v\nEnsure Tibia.dat (7.72) is in the data/ folder.", err)
	}
	r := bytes.NewReader(data)

	// 1. Verify Signature
	var signature uint32
	binary.Read(r, binary.LittleEndian, &signature)
	if signature != DatSignature772 {
		log.Fatalf("Invalid .dat signature: 0x%X (Expected 7.72: 0x%X)", signature, DatSignature772)
	}

	// 2. Read Counts
	var items, creatures, effects, missiles uint16
	binary.Read(r, binary.LittleEndian, &items)
	binary.Read(r, binary.LittleEndian, &creatures)
	binary.Read(r, binary.LittleEndian, &effects)
	binary.Read(r, binary.LittleEndian, &missiles)

	totalIDs := int(items) + int(creatures) + int(effects) + int(missiles)
	fmt.Printf("Found %d total IDs. Parsing...\n", totalIDs)

	// 3. Parse All
	var exportList []ItemAttributes

	// Tibia IDs start at 100
	currentID := uint16(100)

	for i := 0; i < totalIDs; i++ {
		item := parseItem(r, currentID)

		// Only save if it has logic properties we care about
		// (Optional optimization: uncomment to reduce file size)
		// if item.IsGround || item.IsBlocking || item.IsStackable || item.IsContainer {
		exportList = append(exportList, item)
		// }

		currentID++
	}

	// 4. Write JSON
	fmt.Printf("Writing JSON to %s...\n", OutputPath)
	outFile, err := os.Create(OutputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportList); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success!")
}

func parseItem(r *bytes.Reader, id uint16) ItemAttributes {
	item := ItemAttributes{ID: id}

	// --- 1. PARSE FLAGS (Specific to 7.72) ---
	for {
		flag, err := r.ReadByte()
		if err != nil {
			break
		}
		if flag == 0xFF {
			break
		} // End Flag

		switch flag {
		case 0x00: // Ground
			item.IsGround = true
			item.Speed = readUint16(r)
		case 0x01: // Clip
			// Rendering Flag
		case 0x02: // Bottom
			// Rendering Flag
		case 0x03: // Top
			// Decorative
		case 0x04: // Container
			item.IsContainer = true
		case 0x05: // Stackable
			item.IsStackable = true
		case 0x06: // Corpse
			// No data
		case 0x07: // Useable
			item.IsMultiUse = true
		case 0x08: // Writable
			readUint16(r) // Max Length
		case 0x09: // Writable Once (
			readUint16(r) // Max Length
		case 0x0A: // Fluid
			item.IsFluid = true
		case 0x0B: // Splash
			item.IsFluid = true
		case 0x0C: // Blocking
			item.IsBlocking = true
		case 0x0D: // Immovable
			item.IsPickupable = true
		case 0x0E: // Block Missiles
			item.IsMissileBlock = true
		case 0x0F: // Block Path
			// for example Magic Walls or Fire Fields
			item.IsPathBlock = true
		case 0x10: // Pickupable
			item.IsPickupable = true
		case 0x11: // Hangable
			// No data
		case 0x12: // Vertical
			// No data
		case 0x13: // Horizontal
			// No data
		case 0x14: // Rotatable
			item.IsRotatable = true
		case 0x15: // Light Info
			// Note: In 7.72, Light Level/Color are uint16
			item.LightLevel = readUint16(r)
			item.LightColor = readUint16(r)
		case 0x16: // DontHide
			// No data
		case 0x17: // Floor Change
			// No data
		case 0x18: // Shift
			// In 7.72, shift x/y are uint16
			readUint16(r) // x
			readUint16(r) // y
		case 0x19: // Elevation
			item.Elevation = readUint16(r)
		case 0x1A: // Lying Corpse
			// No data
		case 0x1B: // Animate Always
			// No data
		case 0x1C: // Minimap Color
			item.MinimapColor = readUint16(r)
		case 0x1D: // Lens Help
			readUint16(r) // Value
		case 0x1E: // Full Ground
			// No data
		default:
			log.Fatalf("Unknown 7.72 Flag: 0x%X at ID %d. File offset broken.", flag, id)
		}
	}

	// --- 2. SKIP SPRITE DATA ---
	// We must simulate reading sprite headers to advance the cursor correctly.

	width := int(readByte(r))
	height := int(readByte(r))

	if width > 1 || height > 1 {
		readByte(r) // Exact Size
	}

	blendFrames := int(readByte(r))
	xDiv := int(readByte(r))
	yDiv := int(readByte(r))
	zDiv := int(readByte(r))
	animCount := int(readByte(r))

	// Formula for 7.72 sprite count
	totalSprites := width * height * blendFrames * xDiv * yDiv * zDiv * animCount

	// Each sprite ID is 2 bytes
	r.Seek(int64(totalSprites*2), io.SeekCurrent)

	return item
}

// Helpers
func readByte(r *bytes.Reader) uint8 {
	b, _ := r.ReadByte()
	return b
}

func readUint16(r *bytes.Reader) uint16 {
	var v uint16
	binary.Read(r, binary.LittleEndian, &v)
	return v
}
