This is the transition from a "Passive Proxy" (just forwarding bytes) to an "Active Proxy/Bot" (understanding and reacting).

To achieve this, you need to introduce a **Game State** layer and a **Logic Engine** layer.

### 1. The Architecture

We will add two new packages:
1.  **`internal/game/state`**: The "Database" in memory. It holds the current HP, position, inventory, etc. It must be **Thread-Safe**.
2.  **`internal/game/logic`**: The "Brain". It runs in the background, checks the state, and decides to send new packets.

```text
/internal
  /game
    /state              <-- Data Models
      player.go         (HP, Mana, Pos)
      world.go          (Map data, Creatures)
      state.go          (The container)
      
    /logic              <-- The Automation
      engine.go         (The loop)
      healer.go         (Example module)
```

---

### 2. The Game State (`internal/game/state`)

This package simply stores data. It needs a `sync.RWMutex` because the **Proxy Loop** will *write* to it (updating HP from packets) while the **Logic Loop** will *read* from it (checking if HP is low).

```go
package state

import (
	"sync"
)

// GameState holds the world view.
type GameState struct {
	mu     sync.RWMutex
	Player PlayerState
	// Map    MapState
}

type PlayerState struct {
	ID        uint32
	Name      string
	Health    uint16
	MaxHealth uint16
	Mana      uint16
	Position  Position
}

type Position struct {
	X, Y, Z uint16
}

func New() *GameState {
	return &GameState{}
}

// UpdateHealth is called by the Proxy when a stats packet arrives.
func (s *GameState) UpdateHealth(hp, maxHp, mana uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Player.Health = hp
	s.Player.MaxHealth = maxHp
	s.Player.Mana = mana
}

// GetHealth is called by the Logic engine.
func (s *GameState) GetHealth() (uint16, uint16) {
	s.mu.RLock() // Read Lock allows multiple logic modules to read at once
	defer s.mu.RUnlock()
	return s.Player.Health, s.Player.MaxHealth
}
```

---

### 3. The Logic Engine (`internal/game/logic`)

This is where your "Game Logic" lives. It needs access to the **State** (to know what's happening) and the **Server Connection** (to perform actions).

```go
package logic

import (
	"goTibia/internal/game/state"
	"goTibia/internal/packets/game"
	"goTibia/internal/protocol"
	"log"
	"time"
)

type Engine struct {
	State      *state.GameState
	ServerConn *protocol.Connection // To send actions (Attack, Move)
	ClientConn *protocol.Connection // To send info (Text, Effects)
	quit       chan struct{}
}

func NewEngine(s *state.GameState, server, client *protocol.Connection) *Engine {
	return &Engine{
		State:      s,
		ServerConn: server,
		ClientConn: client,
		quit:       make(chan struct{}),
	}
}

// Start begins the background logic loop.
func (e *Engine) Start() {
	ticker := time.NewTicker(200 * time.Millisecond) // Tick 5 times a second
	go func() {
		for {
			select {
			case <-ticker.C:
				e.tick()
			case <-e.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (e *Engine) Stop() {
	close(e.quit)
}

// tick runs all your logic modules
func (e *Engine) tick() {
	e.runAutoHealer()
	// e.runAutoLooter()
	// e.runCavebot()
}

// Example Logic Module
func (e *Engine) runAutoHealer() {
	hp, maxHp := e.State.GetHealth()
	
	if maxHp > 0 && hp < (maxHp/2) {
		log.Println("[Logic] Emergency Heal!")
		
		// Create the packet
		healPacket := &game.CastSpellRequest{
			Spell: "exura vita",
		}
		
		// Send it to the server (Injecting the packet)
		// Note: Protocol.Connection.Send MUST be thread-safe (it usually is if writing to net.Conn)
		if err := e.ServerConn.Send(healPacket); err != nil {
			log.Printf("Failed to cast spell: %v", err)
		}
	}
}
```

---

### 4. Integration: The Handler (`internal/handlers/game/handler.go`)

Now we update your existing handler to wire these two new components together. The Handler becomes the **Interceptor**.

```go
package game

import (
	"goTibia/internal/game/logic"
	"goTibia/internal/game/state"
	game_pkt "goTibia/internal/packets/game"
	"goTibia/internal/protocol"
	"goTibia/internal/proxy"
)

type GameHandler struct {
	TargetAddr string
}

func (h *GameHandler) Handle(client *protocol.Connection) {
	// 1. Connect Backend
	server, err := proxy.ConnectToBackend(h.TargetAddr)
	if err != nil {
		return
	}
	defer server.Close()

	// 2. Initialize State & Logic
	gameState := state.New()
	gameLogic := logic.NewEngine(gameState, server, client)
	
	// Start the background logic (Auto Healer, etc.)
	gameLogic.Start()
	defer gameLogic.Stop()

	// 3. Start the Proxy Loops
	// We need to intercept packets to update gameState
	
	// Channel to wait for disconnect
	done := make(chan struct{})

	// Loop A: Server -> Client (Updating State based on what server says)
	go func() {
		defer close(done)
		for {
			// Read Raw
			msg, err := server.ReadMessage()
			if err != nil { return }

			// Parse just the Opcode to decide if we need to Decode
			reader := protocol.NewPacketReader(msg)
			opcode := reader.ReadByte()

			// --- INTERCEPTION LAYER ---
			switch opcode {
			case game_pkt.OpcodePlayerStats: // e.g. 0xA0
				// Reset reader to parse full packet
				reader.Reset() 
				stats, _ := game_pkt.ParsePlayerStatsMsg(reader)
				
				// UPDATE STATE
				gameState.UpdateHealth(stats.Health, stats.MaxHealth, stats.Mana)
			}
			
			// Forward raw bytes to client (Fastest)
			// OR re-encode if you modified it
			client.WriteMessage(msg)
		}
	}()

	// Loop B: Client -> Server (Tracking user actions)
	go func() {
		for {
			msg, err := client.ReadMessage()
			if err != nil { return }
			
			// You can track client actions here too (e.g. if player moved manually)
			
			server.WriteMessage(msg)
		}
	}()

	<-done
}
```

### Summary of the Flow

1.  **Server sends `StatsPacket (HP: 40)`**:
    *   `Proxy Loop A` receives it.
    *   It parses the packet.
    *   It calls `gameState.UpdateHealth(40)`.
    *   It forwards the packet to the Client (so the user sees 40 HP).
2.  **Logic Loop (`runAutoHealer`) Ticks**:
    *   It calls `gameState.GetHealth()`.
    *   It sees `40`. Max is `100`. 40 < 50.
    *   It constructs a `CastSpellRequest("exura")`.
    *   It calls `e.ServerConn.Send(packet)`.
3.  **Real Server receives `CastSpellRequest`**:
    *   It thinks the client sent it.
    *   It heals the player.
    *   It sends back `StatsPacket (HP: 80)`.
4.  **Cycle repeats**.

### Key Technical Details

1.  **Concurrency Safety:** `GameState` uses `sync.RWMutex`. This is mandatory.
2.  **Packet ID Usage:** In the proxy loop, use `reader.Reset()` (if your reader supports it) or create a new Reader if you need to peek at the Opcode and then parse the body.
3.  **Sending from Logic:** The `protocol.Connection.Send` method creates a *new* `PacketWriter` for every call. This makes it thread-safe regarding buffer usage. The underlying `net.Conn.Write` call is thread-safe in Go.
4.  **Efficiency:** Only parse packets you care about. If the opcode is `0xAB` (some animation) and you don't track it, just forward the raw `[]byte` without calling `ParseAnimationMsg`.