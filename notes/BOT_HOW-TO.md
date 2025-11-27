The best place for this is a dedicated package at **`internal/bot`**.

This package sits at the top of your "Game Logic" hierarchy. It aggregates the **State** (data), the **Connections** (transport), and the **Packets** (definitions) into a unified controller.

Here is the structure and implementation.

### 1. The Directory Structure

```text
/internal
  /game
    /state              (Your GameState struct)
    /types              (Position, etc.)
  /packets
    /game               (Packet definitions)
  /bot                  <-- NEW PACKAGE
    controller.go       (The "Class" definition)
    actions.go          (High-level commands like Walk, Say)
    modules.go          (Loops like AutoHealer)
```

### 2. The Bot Controller (`internal/bot/controller.go`)

This struct acts as the bridge. It holds the connections and the state.

**Key Concept:**
*   **Acting as Client:** You write **C2S** packets (Move, Say) to the **Server Connection**.
*   **Acting as Server:** You write **S2C** packets (Text, Effects) to the **Client Connection**.

```go
package bot

import (
	"goTibia/internal/game/state"
	"goTibia/internal/packets/game"
	"goTibia/internal/protocol"
	"log"
)

type Bot struct {
	// The Eyes: Read-only access to the world state
	State *state.GameState

	// The Hands: Connections to perform actions
	serverConn *protocol.Connection // Send actions here (Walk, Attack)
	clientConn *protocol.Connection // Send info here (Text Message, Effects)
}

func New(s *state.GameState, server, client *protocol.Connection) *Bot {
	return &Bot{
		State:      s,
		serverConn: server,
		clientConn: client,
	}
}

// SendToServer sends a packet to the Real Game Server.
// Use this to perform actions (Walk, Say, Attack).
// T must be a Client Packet (C2S).
func (b *Bot) SendToServer(pkt protocol.Encodable) {
	if err := b.serverConn.SendPacket(pkt); err != nil {
		log.Printf("[Bot] Failed to send to server: %v", err)
	}
}

// SendToClient sends a packet to the Real Game Client.
// Use this to display things to the user (Text, Effects, Windows).
// T must be a Server Packet (S2C).
func (b *Bot) SendToClient(pkt protocol.Encodable) {
	if err := b.clientConn.SendPacket(pkt); err != nil {
		log.Printf("[Bot] Failed to send to client: %v", err)
	}
}
```

### 3. High-Level Actions (`internal/bot/actions.go`)

Don't construct raw packets in your logic loops. Create helper methods here. This makes your bot logic readable.

```go
package bot

import "goTibia/internal/packets/game"

// --- ACTING AS CLIENT (Actions) ---

func (b *Bot) CastSpell(spell string) {
	b.SendToServer(&game.SayRequest{
		Type: 1, // Say/Talk
		Text: spell,
	})
}

func (b *Bot) Walk(dir uint8) {
	b.SendToServer(&game.MoveRequest{
		Direction: dir,
	})
}

// --- ACTING AS SERVER (Feedback) ---

func (b *Bot) SendTextMessage(msg string) {
	b.SendToClient(&game.TextMessageMsg{
		Type:    game.MessageBlue,
		Message: msg,
	})
}

func (b *Bot) SendEffect(x, y uint16, z uint8, effectId uint8) {
	// e.g. MagicEffect packet
}
```

### 4. Logic Modules (`internal/bot/modules.go`)

This is where the automation lives.

```go
package bot

import (
	"log"
	"time"
)

// RunAutoHealer checks HP every 500ms and heals if necessary.
// It returns a function that stops the loop when called.
func (b *Bot) RunAutoHealer() func() {
	ticker := time.NewTicker(500 * time.Millisecond)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				hp, maxHp := b.State.GetHealth()
				
				// Safety check
				if maxHp == 0 { continue }

				percentage := (float64(hp) / float64(maxHp)) * 100

				if percentage < 70 {
					log.Println("[Bot] HP Low - Casting Heal")
					
					// 1. Act as Client: Cast spell
					b.CastSpell("exura")

					// 2. Act as Server: Tell user we healed
					b.SendTextMessage("[Bot] Auto-Healed!")
				}
			}
		}
	}()

	return func() { close(stop); ticker.Stop() }
}
```

### 5. Integration in Handler (`internal/handlers/game/handler.go`)

Initialize the Bot when the connection starts.

```go
import (
	"goTibia/internal/bot"   // Import the new package
	"goTibia/internal/game/state"
	// ...
)

func (h *GameHandler) Handle(client *protocol.Connection) {
	// ... InitSession ...
	
	// 1. Setup State
	gameState := state.New()

	// 2. Setup Bot
	// We pass the connections so the bot can "Drive"
	myBot := bot.New(gameState, protoServerConn, client)

	// 3. Start Modules
	stopHealer := myBot.RunAutoHealer()
	defer stopHealer() // Stop when player disconnects

	// 4. Start Pipe Loop
	// (Your existing loop that updates 'gameState' from S2C packets)
	// ...
}
```

### Why this structure?

1.  **Dependency Graph:** `bot` imports `state`, `packets`, and `protocol`. Nothing imports `bot`. No cycles.
2.  **Abstraction:** Your logic code (`modules.go`) reads like English: `b.CastSpell`, `b.SendTextMessage`. It doesn't know about `PacketWriter` or `Opcode`.
3.  **Duality:** The `Bot` struct explicitly separates `serverConn` (Output) and `clientConn` (Input/Feedback), making it clear which "Role" you are playing when you call a method.