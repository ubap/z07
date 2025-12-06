package game

import (
	"goTibia/internal/bot"
	"goTibia/internal/game/domain"
	"goTibia/internal/game/packets"
	"goTibia/internal/protocol"
	"log"
)

type BotAdapter struct {
	State      *GameState
	ServerConn *protocol.Connection
	ClientConn *protocol.Connection
}

// region Implementing WorldStateReader

func (ba *BotAdapter) GetPlayerPosition() domain.Coordinate {
	// Thread-safe access to state
	ba.State.Lock()
	defer ba.State.Unlock()
	return ba.State.Player.Pos
}

func (ba *BotAdapter) GetInventoryItem(slot domain.InventorySlot) domain.Item {
	ba.State.Lock()
	defer ba.State.Unlock()
	return ba.State.Inventory[slot]
}

func (ba *BotAdapter) GetPlayerID() uint32 {
	ba.State.Lock()
	defer ba.State.Unlock()
	return ba.State.Player.ID
}

// endregion Implementing WorldStateReader

// region Implementing ActionDispatcher

func (ba *BotAdapter) Say(text string) {
	log.Printf("[Game] Say: %s", text)
}

// endregion Implementing ActionDispatcher

// region Implementing ClientManipulator

var _ bot.ClientManipulator = (*BotAdapter)(nil)

func (ba *BotAdapter) SetLocalPlayerLight(lightLevel uint8, color uint8) {
	ba.State.Lock()
	id := ba.State.Player.ID
	ba.State.Unlock()

	// Safety check
	if id == 0 {
		return
	}

	// 2. Construct the specific packet
	pkt := &packets.CreatureLightMsg{
		CreatureID: id,
		LightLevel: lightLevel,
		Color:      color,
	}

	// 3. Send it to the CLIENT (Cheating visually)
	if err := ba.ClientConn.SendPacket(pkt); err != nil {
		// You might log errors here or handle disconnects
	}
}

// endregion Implementing ClientManipulator
