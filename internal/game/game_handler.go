package game

import (
	"fmt"
	"goTibia/internal/bot"
	"goTibia/internal/game/packets"
	"goTibia/internal/game/state"
	"goTibia/internal/protocol"
	"goTibia/internal/proxy"
	"log"
)

type GameHandler struct {
	TargetAddr string
	State      *state.GameState
}

func (h *GameHandler) Handle(client *protocol.Connection) {
	h.State = state.New()

	log.Printf("[Game] New Connection: %s", client.RemoteAddr())

	_, protoServerConn, err := proxy.InitSession(
		"Game",
		client,
		h.TargetAddr,
		packets.ParseLoginRequest,
	)
	if err != nil {
		log.Printf("Game: Failed to initialize session for %s: %v", client.RemoteAddr(), err)
		return
	}
	defer protoServerConn.Close()

	// 4. Start the Bidirectional Pipe
	// We use a channel to detect when either side disconnects.
	// If one side dies, we unblock and the function exits (triggering defers).
	errChan := make(chan error, 2)

	// Loop A: Server -> Client
	go h.pipe(protoServerConn, client, "S2C", errChan)

	// Loop B: Client -> Server
	go h.pipe(client, protoServerConn, "C2S", errChan)

	// 1. Create the Adapter (The Bridge)
	// We inject the specific connections for THIS session
	adapter := &BotAdapter{
		State:      h.State,
		ServerConn: protoServerConn,
		ClientConn: client,
	}

	// 2. Create the Bot (The Brain)
	// The Bot only sees the Interfaces, not the concrete structs
	myBot := bot.NewBot(adapter, adapter, adapter)

	// 3. Start the Bot
	log.Println("[Game] Starting Bot for this session...")
	myBot.Start()

	// 4. Ensure Bot stops when connection dies
	defer myBot.Stop()

	// Wait for the first error (disconnect)
	disconnectErr := <-errChan
	log.Printf("[Game] Connection closed: %v", disconnectErr)
}

// pipe moves data from src to dst indefinitely.
func (h *GameHandler) pipe(src, dst *protocol.Connection, tag string, errChan chan<- error) {
	for {
		// TODO: The proxy could forward raw, unecrypted message right away. This will reduce latency.
		// Right now I can't think of a scenario where we want to edit game packets on the fly.
		rawMsg, packetReader, err := src.ReadMessage()
		if err != nil {
			errChan <- fmt.Errorf("%s Read Error: %w", tag, err)
			return
		}

		if tag == "S2C" {
			h.processPacketsFromServer(packetReader)
		}

		// 2. Forward to Destination
		if err := dst.WriteMessage(rawMsg); err != nil {
			errChan <- fmt.Errorf("%s Write Error: %w", tag, err)
			return
		}
	}
}

func (h *GameHandler) processPacketsFromServer(packetReader *protocol.PacketReader) {
	for packetReader.Remaining() > 0 {

		ctx := packets.ParsingContext{
			PlayerPosition: h.State.CaptureFrame().Player.Pos,
		}

		opcode := packetReader.ReadByte()
		packet, err := packets.ParseS2CPacket(opcode, packetReader, ctx)
		if err != nil {
			log.Printf("[Game] Failed to parse game packet (opcode: 0x%X): %v", opcode, err)
			break
		}
		h.processPacketFromServer(packet)
	}
}

func (h *GameHandler) processPacketFromServer(packet packets.S2CPacket) {
	switch p := packet.(type) {
	case *packets.LoginResponse:
		h.State.SetPlayerId(p.PlayerId)
	case *packets.PingMsg: // Ignore
	case *packets.MapDescriptionMsg:
		h.State.SetPlayerPos(p.PlayerPos)
	case *packets.MoveCreatureMsg:
		// log.Printf("[Game] MoveCreatureMsg %v", p)
	case *packets.MagicEffect:
		log.Printf("[Game] MagicEffect %v", p)
	case *packets.RemoveTileThingMsg:
		// log.Printf("[Game] RemoveTileThingMsg %v", p)
	case *packets.RemoveTileCreatureMsg:
		// log.Printf("[Game] RemoveTileCreatureMsg %v", p)
	case *packets.CreatureLightMsg:
		// log.Printf("[Game] CreatureLightMsg %v", p)
	case *packets.CreatureHealthMsg:
		// log.Printf("[Game] CreatureHealthMsg %v", p)
	case *packets.PlayerIconsMsg:
		log.Printf("[Game] PlayerIconsMsg %v", p)
	case *packets.ServerClosedMsg:
		log.Printf("[Game] ServerClosedMsg %v", p)
	case *packets.AddTileThingMsg:
		log.Printf("[Game] AddTileThingMsg %v", p)
	case *packets.AddInventoryItemMsg:
		h.State.SetEquipment(p.Slot, p.Item)
	case *packets.RemoveInventoryItemMsg:
		h.State.ClearEquipmentSlot(p.Slot)
	case *packets.OpenContainerMsg:
		log.Printf("[Game] OpenContainerMsg %v", p)
	case *packets.RemoveContainerItemMsg:
		log.Printf("[Game] RemoveContainerItemMsg %v", p)
	case *packets.AddContainerItemMsg:
		log.Printf("[Game] AddContainerItemMsg %v", p)

	default:
		log.Printf("[Game] Unhandled game packet type: %T", p)

	}
}
