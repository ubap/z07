package domain

import "fmt"

type Coordinate struct {
	X, Y uint16
	Z    uint8
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
	Position Coordinate
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
