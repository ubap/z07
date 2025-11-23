package game

import (
	"fmt"
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
	if err != nil {
		log.Printf("Game: Failed to initialize session for %s: %v", client.RemoteAddr(), err)
		return
	}
	defer protoServerConn.Close()

	//message, err := protoServerConn.ReadMessage()
	//if err != nil {
	//	log.Printf("Game: Failed to read server response for %s: %v", client.RemoteAddr(), err)
	//	return
	//}

	//loginResult, err := game.ParseLoginResultMessage(message)
	//if err != nil {
	//	log.Printf("Game: Failed to receive login result message for %s: %v", client.RemoteAddr(), err)
	//	return
	//}
	//
	//log.Printf("Game: PlayerId: %d", loginResult.PlayerId)

	// 4. Start the Bidirectional Pipe
	// We use a channel to detect when either side disconnects.
	// If one side dies, we unblock and the function exits (triggering defers).
	errChan := make(chan error, 2)

	// Loop A: Server -> Client
	go h.pipe(protoServerConn, client, "S2C", errChan)

	// Loop B: Client -> Server
	go h.pipe(client, protoServerConn, "C2S", errChan)

	// Wait for the first error (disconnect)
	disconnectErr := <-errChan
	log.Printf("[Game] Connection closed: %v", disconnectErr)
}

// pipe moves data from src to dst indefinitely.
func (h *GameHandler) pipe(src, dst *protocol.Connection, tag string, errChan chan<- error) {
	for {
		// 1. Read Raw Encrypted Message
		rawMsg, packetReader, err := src.ReadMessage()
		if err != nil {
			errChan <- fmt.Errorf("%s Read Error: %w", tag, err)
			return
		}

		if tag == "C2S" {
			for packetReader.Remaining() > 0 {
				opcode := packetReader.ReadByte()
				packet, err := game.ParseS2CPacket(opcode, packetReader)
				if err != nil {
					log.Printf("[Game] Failed to parse game packet (opcode: 0x%X): %v", opcode, err)
					break
				}
				log.Printf("[Game] Received game message: %v", packet)
			}
		}

		// 2. Forward to Destination
		if err := dst.WriteMessage(rawMsg); err != nil {
			errChan <- fmt.Errorf("%s Write Error: %w", tag, err)
			return
		}
	}
}
