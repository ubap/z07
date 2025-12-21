package domain

import "fmt"

type Position struct {
	X, Y uint16
	Z    uint8
}

type Container struct {
	// 1. Identification
	ID     uint8  // The window index (0-15). Crucial for logic.
	ItemID uint16 // The visual ID (e.g., 1988 for Brown Backpack).
	Name   string

	// 2. State
	Capacity  uint8 // Total slots available (e.g. 20).
	HasParent bool  // Useful to know if this is inside another container.

	// 3. Contents
	Items []Item
}

type Item struct {
	ID       uint16
	Count    uint8 // Used for stack count, fluid type, or rune charges
	HasCount bool  // Helper to know if we should write the Count byte
}

func (i Item) String() string {
	if i.HasCount || i.Count > 1 {
		return fmt.Sprintf("ID: %d (x%d)", i.ID, i.Count)
	}

	// 2. Simple items
	return fmt.Sprintf("ID: %d", i.ID)
}

type Tile struct {
	Position Position
	Ground   *Item
	Items    []Item
}

type Direction uint8

const (
	North Direction = 0
	East  Direction = 1
	South Direction = 2
	West  Direction = 3
)

type Player struct {
	ID   uint32
	Name string
	Pos  Position
}

type SkillType uint8

const (
	Fist     SkillType = 0
	Club     SkillType = 1
	Sword    SkillType = 2
	Axe      SkillType = 3
	Distance SkillType = 4
	Shield   SkillType = 5
	Fishing  SkillType = 6

	Maglevel SkillType = 7
	Level    SkillType = 8

	SkillFirst = Fist
	SkillLast  = Fishing
)

type Skill struct {
	Level   uint8
	Percent uint8
}
