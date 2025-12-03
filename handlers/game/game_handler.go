package game

import (
	"errors"
	"fmt"
	"goTibia/bot"
	"goTibia/game/state"
	"goTibia/packets/game"
	"goTibia/protocol"
	"goTibia/proxy"
	"log"
)

type GameHandler struct {
	TargetAddr string
	State      *state.GameState
	// You could add "DB *sql.DB" here later!
	Bot bot.Bot
}

func (h *GameHandler) Handle(client *protocol.Connection) {
	h.State = state.New()

	log.Printf("[Game] New Connection: %s", client.RemoteAddr())

	_, protoServerConn, err := proxy.InitSession(
		"Game",
		client,
		h.TargetAddr,
		game.ParseLoginRequest,
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

	h.Bot = *bot.New(h.State, protoServerConn, client)
	h.Bot.Start()

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
		opcode := packetReader.ReadByte()
		packet, err := game.ParseS2CPacket(opcode, packetReader)
		if errors.Is(err, game.ErrNotFullyImplemented) {
			h.processPacketFromServer(packet)
		}
		if err != nil {
			log.Printf("[Game] Failed to parse game packet (opcode: 0x%X): %v", opcode, err)
			break
		}
		h.processPacketFromServer(packet)
	}
}

func (h *GameHandler) processPacketFromServer(packet game.S2CPacket) {
	switch p := packet.(type) {
	case *game.LoginResponse:
		log.Printf("[Game] Received game message: %d", p.PlayerId)
		h.State.SetPlayerId(p.PlayerId)
	case *game.MapDescription:
		h.State.SetPosition(p.Pos)
	case *game.MoveCreatureMsg:
		// log.Printf("[Game] MoveCreatureMsg %v", p)
	case *game.MagicEffect:
		log.Printf("[Game] MagicEffect %v", p)
	case *game.RemoveTileThingMsg:
		log.Printf("[Game] RemoveTileThingMsg %v", p)
	case *game.RemoveTileCreatureMsg:
		log.Printf("[Game] RemoveTileCreatureMsg %v", p)
	case *game.CreatureLightMsg:
		log.Printf("[Game] CreatureLightMsg %v", p)
	case *game.CreatureHealthMsg:
		log.Printf("[Game] CreatureHealthMsg %v", p)
	case *game.PlayerIconsMsg:
		log.Printf("[Game] PlayerIconsMsg %v", p)
	case *game.ServerClosedMsg:
		log.Printf("[Game] ServerClosedMsg %v", p)
	case *game.AddTileThingMsg:
		log.Printf("[Game] AddTileThingMsg %v", p)
	case *game.AddInventoryItemMsg:
		log.Printf("[Game] AddInventoryItemMsg %v", p)
	case *game.RemoveInventoryItemMsg:
		log.Printf("[Game] RemoveInventoryItemMsg %v", p)
	case *game.PingMsg:
		// Ignore
	default:
		log.Printf("[Game] Unhandled game packet type: %T", p)

	}
}
