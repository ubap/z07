package game

import (
	"goTibia/packets/game"
	"goTibia/protocol"
	"goTibia/proxy"
	"log"
)

func HandleGameConnection(client *protocol.Connection, targetAddr string) {
	log.Printf("[Game] New Connection: %s", client.RemoteAddr())

	packetReader, err := client.ReadMessage()
	if err != nil {
		log.Printf("Game: error reading message from %s: %v", client.RemoteAddr(), err)
		return
	}

	_, err = game.ParseLoginRequest(packetReader)
	if err != nil {
		log.Printf("Login: Failed to parse login packet: %v", err)
		return
	}

	protoServerConn, err := proxy.ConnectToBackend(targetAddr)
	if err != nil {
		log.Printf("Login: Failed to connect to %s: %v", client.RemoteAddr(), err)
		return
	}
	defer protoServerConn.Close()

}
