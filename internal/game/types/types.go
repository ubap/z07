package types

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

type Tile struct {
	Position Position
	Ground   *Item
	Items    []Item
}
