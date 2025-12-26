package state

import "goTibia/internal/game/domain"

type WorldSnapshot struct {
	Player     domain.Player
	Equipment  [11]domain.Item
	Containers [16]*domain.Container
	WorldMap   map[domain.Position]*domain.Tile
}
