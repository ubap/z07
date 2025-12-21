package game

import (
	"goTibia/internal/bot"
	"goTibia/internal/game/state"
	"goTibia/internal/protocol"
)

type GameSession struct {
	ID               string
	State            *state.GameState
	Bot              *bot.Bot
	ClientConn       protocol.Connection
	ServerConn       protocol.Connection
	c2sRawPacketChan chan []byte // Used to listen to C2S packets
	ErrChan          chan error
}

func newGameSession(client protocol.Connection, server protocol.Connection, gameState *state.GameState) *GameSession {
	return &GameSession{
		ID:               client.RemoteAddr().String(),
		State:            gameState,
		ClientConn:       client,
		ServerConn:       server,
		c2sRawPacketChan: make(chan []byte, 1024),
		ErrChan:          make(chan error, 100),
		Bot:              bot.NewBot(gameState, client, server),
	}
}
