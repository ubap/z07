package bot

import (
	"log"
	"sync"
	"time"
	"z07/internal/game/packets"
	"z07/internal/game/state"
	"z07/internal/protocol"
)

type Bot struct {
	state *state.GameState

	clientConn protocol.Connection
	serverConn protocol.Connection
	stopChan   chan struct{}  // The broadcast channel
	wg         sync.WaitGroup // To wait for modules to finish
	stopOnce   sync.Once      // To ensure we close the channel only once

	// Module states
	fishingEnabled   bool
	lighthackEnabled bool
	lighthackLevel   uint8
	lighthackColor   uint8

	lastLookedAt uint16
}

func NewBot(state *state.GameState, clientConn protocol.Connection, serverConn protocol.Connection) *Bot {
	return &Bot{
		state: state,

		clientConn: clientConn,
		serverConn: serverConn,
		stopChan:   make(chan struct{}),

		lighthackLevel: 0xFF,
		lighthackColor: 0xD7,
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
			if !b.lighthackEnabled {
				continue
			}
			pId := b.state.CaptureFrame().Player.ID

			if pId == 0 {
				continue
			}

			pkt := &packets.CreatureLightMsg{
				CreatureID: pId,
				LightLevel: b.lighthackLevel,
				Color:      b.lighthackColor,
			}
			err := b.clientConn.SendPacket(pkt)
			if err != nil {
				log.Printf("[Bot][LightHack] Failed to send light packet: %v", err)
				return
			}
		}
	}
}

// InterceptS2CPacket has to return immediately.
func (b *Bot) InterceptS2CPacket(data []byte) ([]byte, error) {
	opcode := packets.S2COpcode(data[0])
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
	opcode := packets.C2SOpcode(data[0])

	// LOG FOR TESTING
	// This only prints in terminal so you can copy it
	//fmt.Println("\n--- COPY ME FOR TEST ---")
	//fmt.Println(FormatForTest(fmt.Sprintf("Opcode: %d", opcode), data))
	//fmt.Println("------------------------")

	switch opcode {
	case packets.C2SLookRequest:
		pr := protocol.NewPacketReader(data)
		pr.ReadUint8() // skip opcode
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
	b.lastLookedAt = p.ItemId
}
