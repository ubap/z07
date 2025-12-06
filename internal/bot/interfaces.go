package bot

import "goTibia/internal/game/domain"

type WorldStateReader interface {
	GetPlayerPosition() domain.Coordinate
	GetInventoryItem(slot domain.InventorySlot) domain.Item
	GetPlayerID() uint32
}

type ActionDispatcher interface {
	Say(text string)
}

// ClientManipulator allows manipulating the client's visual state.
// This is a separate interface, there might be a different implementation,
// which writes to the client process.
type ClientManipulator interface {
	SetLocalPlayerLight(lightLevel uint8, color uint8)
}
