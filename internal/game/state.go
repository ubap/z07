package game

import (
	"goTibia/internal/game/domain"
	"sync"
)

// GameState is the container for all tracking data.
// It must be thread-safe because the Network Loop writes to it
// and the Logic Loop reads from it simultaneously.
type GameState struct {
	mu        sync.RWMutex
	Player    PlayerModel
	Inventory map[domain.InventorySlot]domain.Item
}

func New() *GameState {
	return &GameState{
		Inventory: make(map[domain.InventorySlot]domain.Item),
	}
}

type PlayerModel struct {
	ID   uint32
	Name string
	Pos  domain.Coordinate
}

// Lock
// Use this when UPDATING data (e.g., in GameHandler/Network)
// No other thread can Read OR Write while this is held.
func (s *GameState) Lock() {
	s.mu.Lock()
}

func (s *GameState) Unlock() {
	s.mu.Unlock()
}

// RLock
// Use this when READING data (e.g., in Bot Logic)
// Multiple threads can hold this at the same time.
// It only blocks if someone is currently holding the Write Lock.
func (s *GameState) RLock() {
	s.mu.RLock()
}

func (s *GameState) RUnlock() {
	s.mu.RUnlock()
}
