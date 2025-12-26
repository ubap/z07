package bot

import (
	"goTibia/internal/game/domain"
	"goTibia/internal/game/packets"
	"goTibia/internal/game/state"
	"log"
	"time"
)

func (b *Bot) loopFishing() {
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	log.Println("[Bot] Auto fishing started")

	for {
		select {
		case <-b.stopChan:
			return

		case <-ticker.C:
			frame := b.state.CaptureFrame()
			pos := frame.Player.Pos

			fishPos := b.findFishPos(frame, pos)
			if fishPos == nil {
				continue
			}
			tileWithFish := frame.WorldMap[*fishPos]

			pkt := packets.UseItemWithCrosshairRequest{
				FromPos:      domain.NewInventoryPosition(domain.SlotAmmo),
				FromItemId:   3483,
				FromStackPos: 0,

				ToPos:      *fishPos,
				ToItemId:   tileWithFish.Ground.ID,
				ToStackPos: 0,
			}

			// map is off, invesgiate why
			b.serverConn.SendPacket(&pkt)
		}
	}
}

func (b *Bot) findFishPos(frame state.WorldSnapshot, pos domain.Position) *domain.Position {
	for x := pos.X - 5; x <= pos.X+5; x++ {
		for y := pos.Y - 5; y <= pos.Y+5; y++ {
			currentPos := domain.Position{X: x, Y: y, Z: pos.Z}
			tile := frame.WorldMap[currentPos]
			if tile.Ground.ID == 4598 {
				log.Printf("[Bot] Found water with tile at (%d, %d, %d)", x, y, pos.Z)
				return &currentPos
			}
		}

	}
	return nil
}
