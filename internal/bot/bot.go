package bot

import (
	"goTibia/internal/game/packets"
	"goTibia/internal/game/state"
	"goTibia/internal/protocol"
	"log"
	"sync"
	"time"
)

type Bot struct {
	state *state.GameState

	UserActions chan packets.C2SPacket // packets sent by client to server

	clientConn *protocol.Connection
	stopChan   chan struct{}  // The broadcast channel
	wg         sync.WaitGroup // To wait for modules to finish
	stopOnce   sync.Once      // To ensure we close the channel only once
}

func NewBot(state *state.GameState, clientConn *protocol.Connection, serverConn *protocol.Connection) *Bot {
	return &Bot{
		state: state,

		UserActions: make(chan packets.C2SPacket, 100),

		clientConn: clientConn,
		stopChan:   make(chan struct{}),
	}
}

func (b *Bot) Start() {
	log.Println("[Bot] Engine started")

	// 1. The Light Hack (For testing S2C injection)
	b.runModule("LightHack", b.loopLightHack)
	b.runModule("HandleUserAction", b.loopHandleUserAction)
}

func (b *Bot) Stop() {
	b.stopOnce.Do(func() {
		log.Println("[Bot] Stopping engine...")
		close(b.stopChan) // This broadcasts the signal to ALL loops instantly
	})

	b.wg.Wait()
	log.Println("[Bot] Engine stopped cleanly.")
}

func (b *Bot) runModule(name string, logic func()) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		log.Printf("[Bot] Module %s running", name)
		logic()
		log.Printf("[Bot] Module %s stopped", name)
	}()
}

func (b *Bot) loopLightHack() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	log.Println("[Bot] LightHack started")

	for {
		select {
		case <-b.stopChan:
			return

		case <-ticker.C:
			pId := b.state.CaptureFrame().Player.ID

			if pId == 0 {
				continue
			}

			pkt := &packets.WorldLightMsg{
				LightLevel: 0xFF,
				Color:      0xD7,
			}
			b.clientConn.SendPacket(pkt)
		}
	}
}

func (b *Bot) loopHandleUserAction() {
	for {
		select {
		case <-b.stopChan:
			return
		case packet := <-b.UserActions:
			b.handleUserAction(packet)
		}
	}
}

func (b *Bot) handleUserAction(packet packets.C2SPacket) {
	switch p := packet.(type) {
	case *packets.LookRequest:
		log.Printf("User looked at item ID: %d", p.ItemId)
	}
}
