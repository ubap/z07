package domain

import "fmt"

type Position struct {
	X uint16
	Y uint16
	Z uint8
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
