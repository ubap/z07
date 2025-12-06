package game

import (
	"goTibia/internal/game/domain"
	"log"
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
	Pos  domain.Position
}

func (s *GameState) SetPlayerName(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Player.Name = name
}

func (s *GameState) SetPlayerId(ID uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Player.ID = ID
}

func (s *GameState) SetPosition(pos domain.Position) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Player.Pos = pos

	log.Printf("Player position updated to: X=%d, Y=%d, Z=%d", pos.X, pos.Y, pos.Z)
}

func (s *GameState) SetInventoryItem(slot domain.InventorySlot, item domain.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Inventory[slot] = item
}

func (s *GameState) RemoveInventoryItem(slot domain.InventorySlot) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Inventory, slot)
}

func (s *GameState) GetName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Player.Name
}
