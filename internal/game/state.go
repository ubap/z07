package game

import (
	"goTibia/internal/game/types"
	"log"
	"sync"
)

// GameState is the container for all tracking data.
// It must be thread-safe because the Network Loop writes to it
// and the Logic Loop reads from it simultaneously.
type GameState struct {
	mu     sync.RWMutex
	Player PlayerModel
}

func New() *GameState {
	return &GameState{}
}

type PlayerModel struct {
	ID   uint32
	Name string
	Pos  types.Position
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

func (s *GameState) SetPosition(pos types.Position) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Player.Pos = pos

	log.Printf("Player position updated to: X=%d, Y=%d, Z=%d", pos.X, pos.Y, pos.Z)
}

func (s *GameState) GetName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Player.Name
}
