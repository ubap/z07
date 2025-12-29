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
	fishingEnabled bool

	state *state.GameState

	clientConn protocol.Connection
	serverConn protocol.Connection
	stopChan   chan struct{}  // The broadcast channel
	wg         sync.WaitGroup // To wait for modules to finish
	stopOnce   sync.Once      // To ensure we close the channel only once
}

func NewBot(state *state.GameState, clientConn protocol.Connection, serverConn protocol.Connection) *Bot {
	return &Bot{
		state: state,

		clientConn: clientConn,
		serverConn: serverConn,
		stopChan:   make(chan struct{}),
	}
}

func (b *Bot) Start() {
	log.Println("[Bot] Engine started")

	b.runModule("LightHack", b.loopLightHack)
	b.runModule("Fishing", b.loopFishing)
	b.runModule("UI", b.loopWebUI)
}

func (b *Bot) StartUIOnly() {
	log.Println("[Bot] Engine started in UI-only mode")

	b.runModule("UI", b.loopWebUI)
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
			err := b.clientConn.SendPacket(pkt)
			if err != nil {
				return
			}
		}
	}
}

// InterceptS2CPacket has to return immediately.
func (b *Bot) InterceptS2CPacket(data []byte) ([]byte, error) {
	pr := protocol.NewPacketReader(data)
	opcode := packets.S2COpcode(pr.ReadUint8())
	switch opcode {
	case packets.S2CSLoginQueue:
		pw := protocol.NewPacketWriter()
		msg := packets.LoginQueueMsg{Message: "Queue hack active.", RetryTimeSeconds: 1}
		msg.Encode(pw)
		return pw.GetBytes()
	}
	return data, nil
}

// InterceptC2SPacket has to return immediately.
func (b *Bot) InterceptC2SPacket(data []byte) ([]byte, error) {
	pr := protocol.NewPacketReader(data)
	firstByte, err := pr.PeekUint8()
	if err != nil {
		log.Println("[Bot] InterceptC2SPacket err:", err)
		return data, err
	}
	opcode := packets.C2SOpcode(firstByte)
	switch opcode {
	case packets.C2SLookRequest:
		b.handleLookRequest(pr)
	}
	return data, nil
}

func (b *Bot) handleLookRequest(pr *protocol.PacketReader) {
	p, err := packets.ParseLookRequest(pr)
	if err != nil {
		log.Printf("Failed to parse look request: %v", err)
		return
	}
	log.Printf("User looked at item ID: %d", p.ItemId)
}
