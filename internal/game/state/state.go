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
		Player:    gs.player,
		Equipment: gs.equipment,
	}

	return snap
}

func (gs *GameState) SetPlayerId(pId uint32) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.player.ID = pId
}

func (gs *GameState) SetPlayerPos(pos domain.Coordinate) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.player.Pos = pos
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
