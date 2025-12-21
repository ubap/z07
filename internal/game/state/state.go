package state

import (
	"goTibia/internal/game/domain"
	"sync"
)

// GameState is the container for all tracking data.
// It must be thread-safe because the Network Loop writes to it
// and the Logic Loop reads from it simultaneously.
type GameState struct {
	player     domain.Player
	equipment  [11]domain.Item
	containers [16]*domain.Container // nil means closed

	mu sync.RWMutex
}

func New() *GameState {
	return &GameState{}
}

func (gs *GameState) CaptureFrame() WorldSnapshot {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	snap := WorldSnapshot{
		Player:     gs.player,
		Equipment:  gs.equipment,
		Containers: gs.containers,
	}

	return snap
}

func (gs *GameState) SetPlayerId(pId uint32) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.player.ID = pId
}

func (gs *GameState) SetPlayerPos(pos domain.Position) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.player.Pos = pos
}

func (gs *GameState) SetPlayerName(Name string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.player.Name = Name
}

func (gs *GameState) SetEquipment(slot domain.EquipmentSlot, item domain.Item) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if slot > 0 && int(slot) < len(gs.equipment) {
		gs.equipment[slot] = item
	}
}

func (gs *GameState) ClearEquipmentSlot(slot domain.EquipmentSlot) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if slot > 0 && int(slot) < len(gs.equipment) {
		gs.equipment[slot] = domain.Item{}
	}
}

// OpenContainer overwrites the container slot with the full state provided.
func (gs *GameState) OpenContainer(c domain.Container) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if int(c.ID) >= len(gs.containers) {
		return
	}

	// 2. Store the Container
	// Since 'c' is passed by value, we have a copy of the struct headers.
	// However, c.Items is a slice (pointer to array).
	// Because the Handler creates this slice fresh from the packet and then discards it,
	// it is generally safe to take ownership of it here without a deep copy loop.
	gs.containers[c.ID] = &c
}

func (gs *GameState) CloseContainer(cId uint8) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if int(cId) >= len(gs.containers) {
		return
	}
	gs.containers[cId] = nil
}

func (gs *GameState) RemoveContainerItem(cId uint8, slot uint8) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if int(cId) >= len(gs.containers) {
		return
	}

	container := gs.containers[cId]
	if container == nil {
		return
	}

	if int(slot) >= len(container.Items) {
		return
	}

	// 4. Perform the Removal (Shift Logic)
	// We take everything BEFORE the slot, and append everything AFTER the slot.
	// Go automatically handles the shifting of memory.
	container.Items = append(container.Items[:slot], container.Items[slot+1:]...)
}

func (gs *GameState) AddContainerItem(cId uint8, item domain.Item) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if int(cId) >= len(gs.containers) {
		return
	}

	container := gs.containers[cId]
	if container == nil {
		return
	}

	container.Items = append([]domain.Item{item}, container.Items...)
}

func (gs *GameState) UpdateContainerItem(cId uint8, slot uint8, item domain.Item) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if int(cId) >= len(gs.containers) {
		return
	}

	container := gs.containers[cId]
	if container == nil {
		return
	}

	if int(slot) >= len(container.Items) {
		return
	}

	container.Items[slot] = item
}
