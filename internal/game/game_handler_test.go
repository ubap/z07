package game

import (
	"goTibia/internal/game/domain"
	"goTibia/internal/game/packets"
	"goTibia/internal/game/state"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessPacketFromServer(t *testing.T) {
	gameState := state.New()
	session := &GameSession{
		State: gameState,
	}

	t.Run("Handle LoginResponse", func(t *testing.T) {
		pkt := &packets.LoginResponse{
			PlayerId: 12345,
		}

		session.processPacketFromServer(pkt)

		require.Equal(t, uint32(12345), gameState.CaptureFrame().Player.ID)
	})

	t.Run("Handle OpenContainerMsg", func(t *testing.T) {
		pkt := &packets.OpenContainerMsg{
			ContainerID:   1,
			ContainerItem: domain.Item{ID: 2853},
			ContainerName: "Backpack",
			Capacity:      20,
			Items: []domain.Item{
				{ID: 3031, Count: 100},
			},
		}

		session.processPacketFromServer(pkt)

		containers := gameState.CaptureFrame().Containers
		c := containers[1]
		require.NotNil(t, c, "container with ID 1 should be present in state")
		require.Equal(t, "Backpack", c.Name)
		require.Equal(t, uint8(20), c.Capacity)
		require.Equal(t, uint16(2853), c.ItemID)

		require.Len(t, c.Items, 1, "container should have 1 item")
		require.Equal(t, uint16(3031), c.Items[0].ID)
		require.Equal(t, uint8(100), c.Items[0].Count)
	})

	t.Run("Handle AddInventoryItemMsg", func(t *testing.T) {
		pkt := &packets.AddInventoryItemMsg{
			Slot: 1,
			Item: domain.Item{ID: 3350},
		}

		session.processPacketFromServer(pkt)

		equip := gameState.CaptureFrame().Equipment
		item := equip[1]
		require.Equal(t, uint16(3350), item.ID)
	})

	t.Run("Handle SetPlayerPos", func(t *testing.T) {
		targetPos := domain.Position{X: 32368, Y: 32234, Z: 7}

		pkt := &packets.MapDescriptionMsg{
			PlayerPos: targetPos,
		}

		session.processPacketFromServer(pkt)

		currentPos := gameState.CaptureFrame().Player.Pos
		require.Equal(t, targetPos, currentPos, "Player position in state should match the packet position")
	})
}
