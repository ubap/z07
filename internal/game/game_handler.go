package game

import (
	"fmt"
	"goTibia/internal/bot"
	"goTibia/internal/game/domain"
	"goTibia/internal/game/packets"
	"goTibia/internal/game/state"
	"goTibia/internal/protocol"
	"goTibia/internal/proxy"
	"log"
)

type GameHandler struct {
	TargetAddr string
}

type GameSession struct {
	ID         string
	State      *state.GameState
	Bot        *bot.Bot
	ClientConn *protocol.Connection
	ServerConn *protocol.Connection
	ErrChan    chan error
}

func (h *GameHandler) Handle(client *protocol.Connection) {
	log.Printf("[Game] New Connection: %s", client.RemoteAddr())

	loginPkt, protoServerConn, err := proxy.InitSession(
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

	gameState := state.New()
	gameState.SetPlayerName(loginPkt.CharacterName)

	session := &GameSession{
		ID:         client.RemoteAddr().String(),
		State:      gameState,
		ClientConn: client,
		ServerConn: protoServerConn,
		ErrChan:    make(chan error, 100),
	}
	session.Bot = bot.NewBot(gameState, client, protoServerConn)

	go session.loopS2C()
	go session.loopC2S()
	go session.Bot.Start()

	disconnectErr := <-session.ErrChan
	log.Printf("[Game] Connection closed: %v", disconnectErr)
	session.Bot.Stop()
}

func (g *GameSession) loopS2C() {
	for {
		// 1. Read Raw
		rawMsg, packetReader, err := g.ServerConn.ReadMessage()
		if err != nil {
			g.ErrChan <- fmt.Errorf("S2C Read: %w", err)
			return
		}
		if err := g.ClientConn.WriteMessage(rawMsg); err != nil {
			g.ErrChan <- fmt.Errorf("S2C Write: %w", err)
			return
		}

		go g.processPacketsFromServer(packetReader)
	}
}

func (g *GameSession) loopC2S() {
	for {
		rawMsg, _, err := g.ClientConn.ReadMessage()
		if err != nil {
			g.ErrChan <- fmt.Errorf("C2S Read: %w", err)
			return
		}
		if err := g.ServerConn.WriteMessage(rawMsg); err != nil {
			g.ErrChan <- fmt.Errorf("C2S Write: %w", err)
			return
		}

		packetReader := protocol.NewPacketReader(rawMsg)
		opcode := packetReader.ReadByte()
		packet, err := packets.ParseC2SPacket(opcode, packetReader)
		if err == nil {
			g.Bot.UserActions <- packet
		}

	}
}

func (g *GameSession) processPacketsFromServer(packetReader *protocol.PacketReader) {
	for packetReader.Remaining() > 0 {

		ctx := packets.ParsingContext{
			PlayerPosition: g.State.CaptureFrame().Player.Pos,
		}

		opcode := packetReader.ReadByte()
		packet, err := packets.ParseS2CPacket(opcode, packetReader, ctx)
		if err != nil {
			log.Printf("[Game] Failed to parse game packet (opcode: 0x%X): %v", opcode, err)
			break
		}
		g.processPacketFromServer(packet)
	}
}

func (g *GameSession) processPacketFromServer(packet packets.S2CPacket) {
	switch p := packet.(type) {
	case *packets.LoginResponse:
		g.State.SetPlayerId(p.PlayerId)
	case *packets.PingMsg: // Ignore
	case *packets.MapDescriptionMsg:
		g.State.SetPlayerPos(p.PlayerPos)
	case *packets.MoveCreatureMsg:
		// log.Printf("[Game] MoveCreatureMsg %v", p)
	case *packets.MagicEffect:
		// log.Printf("[Game] MagicEffect %v", p)
	case *packets.RemoveTileThingMsg:
		// log.Printf("[Game] RemoveTileThingMsg %v", p)
	case *packets.RemoveTileCreatureMsg:
		// log.Printf("[Game] RemoveTileCreatureMsg %v", p)
	case *packets.WorldLightMsg:
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
		g.State.SetEquipment(p.Slot, p.Item)
	case *packets.RemoveInventoryItemMsg:
		g.State.ClearEquipmentSlot(p.Slot)
	case *packets.OpenContainerMsg:
		g.handleContainerOpen(p)
	case *packets.CloseContainerMsg:
		g.State.CloseContainer(p.ContainerID)
	case *packets.RemoveContainerItemMsg:
		g.State.RemoveContainerItem(p.ContainerID, p.Slot)
	case *packets.AddContainerItemMsg:
		g.State.AddContainerItem(p.ContainerID, p.Item)
	case *packets.UpdateContainerItemMsg:
		g.State.UpdateContainerItem(p.ContainerID, p.Slot, p.Item)
	case *packets.UpdateTileItemMsg:
		log.Printf("[Game] UpdateTileItem %v", p)
	case *packets.PlayerSkillsMsg:
		log.Printf("[Game] PlayerSkillsMsg %v", p)
	case *packets.PlayerStatsMsg:
		log.Printf("[Game] PlayerStatsMsg %v", p)

	default:
		log.Printf("[Game] Unhandled game packet type: %T", p)

	}
}

func (h *GameSession) handleContainerOpen(p *packets.OpenContainerMsg) {
	// 1. Translate Packet -> Domain
	container := domain.Container{
		ID:       p.ContainerID,
		ItemID:   p.ContainerItem.ID,
		Name:     p.ContainerName,
		Capacity: p.Capacity,
		Items:    p.Items,
	}

	h.State.OpenContainer(container)
}
