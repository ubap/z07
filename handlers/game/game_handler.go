package game

import (
	"goTibia/packets/game"
	"goTibia/protocol"
	"goTibia/proxy"
	"log"
)

type GameHandler struct {
	TargetAddr string
	// You could add "DB *sql.DB" here later!
}

func (h *GameHandler) Handle(client *protocol.Connection) {
	log.Printf("[Game] New Connection: %s", client.RemoteAddr())

	_, protoServerConn, err := proxy.InitSession(
		"Game",
		client,
		h.TargetAddr,
		game.ParseLoginRequest,
	)
	defer protoServerConn.Close()

	message, err := protoServerConn.ReadMessage()
	if err != nil {
		log.Printf("Game: Failed to read server response for %s: %v", client.RemoteAddr(), err)
		return
	}

	loginResult, err := game.ParseLoginResultMessage(message)
	if err != nil {
		log.Printf("Game: Failed to receive login result message for %s: %v", client.RemoteAddr(), err)
		return
	}

	log.Printf("Game: PlayerId: %d", loginResult.PlayerId)
}
