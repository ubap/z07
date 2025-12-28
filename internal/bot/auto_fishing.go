package bot

import (
	"goTibia/internal/game/domain"
	"goTibia/internal/game/packets"
	"goTibia/internal/game/state"
	"log"
	"time"
)

const fishingRodItemId = 3483

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

			fishingRod := frame.FindItemInEqAndOpenWindows(fishingRodItemId)
			if fishingRod == nil {
				log.Println("[Bot] No fishing rod found in equipment or containers.")
				continue
			}

			tileWithFish := b.findFishPos(frame)
			if tileWithFish == nil {
				continue
			}

			pkt := packets.UseItemWithCrosshairRequest{
				FromPos:      domain.NewInventoryPosition(domain.SlotAmmo),
				FromItemId:   3483,
				FromStackPos: 0,

				ToPos:      tileWithFish.Position,
				ToItemId:   tileWithFish.Items[0].ID,
				ToStackPos: 0,
			}

			b.serverConn.SendPacket(&pkt)
		}
	}
}

func (b *Bot) findFishPos(frame state.WorldSnapshot) *domain.Tile {
	pos := frame.Player.Pos
	for x := pos.X - 7; x <= pos.X+7; x++ {
		for y := pos.Y - 5; y <= pos.Y+5; y++ {
			currentPos := domain.Position{X: x, Y: y, Z: pos.Z}
			tile, ok := frame.WorldMap[currentPos]
			if ok && tile.Items[0].ID == 4598 {
				log.Printf("[Bot] Found water with tile at (%d, %d, %d)", x, y, pos.Z)
				return tile
			}
		}

	}
	return nil
}
