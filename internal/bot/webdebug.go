package bot

import (
	"context"
	"fmt"
	"goTibia/internal/game/domain"
	"log"
	"net/http"
	"time"
)

func (b *Bot) loopWebDebug() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", b.handleRenderMap)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Run server in a goroutine so we can listen for the stop signal
	go func() {
		log.Println("[WebDebug] Map debugger available at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[WebDebug] Error: %v", err)
		}
	}()

	// Wait for stop signal
	<-b.stopChan

	// Shutdown the server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func (b *Bot) handleRenderMap(w http.ResponseWriter, r *http.Request) {
	// 1. Capture a consistent snapshot of the game state
	// We call this once to ensure player pos and map data are synchronized
	frame := b.state.CaptureFrame()

	pPos := frame.Player.Pos
	currentZ := pPos.Z
	worldTiles := frame.WorldMap // map[domain.Position]*domain.Tile

	// 2. Configuration
	const radius = 8 // How many tiles to show in each direction

	// 3. Start HTML Header
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<html><head>
    <title>goTibia Map Debugger</title>
    <meta http-equiv="refresh" content="1">
	<style>
		body { background: #121212; color: #eee; font-family: 'Courier New', monospace; display: flex; flex-direction: column; align-items: center; padding-top: 30px; }
		.stats { margin-bottom: 15px; background: #222; padding: 10px 20px; border-radius: 5px; border: 1px solid #444; }
        .map-container { 
            background: #000; 
            padding: 15px; 
            border: 3px solid #333; 
            line-height: 14px; 
            letter-spacing: 3px; 
            font-size: 16px;
            box-shadow: 0 0 20px rgba(0,0,0,0.5);
        }
		.tile { display: inline-block; width: 16px; height: 16px; vertical-align: middle; text-align: center; }
		.player { color: #00ff00; font-weight: bold; text-shadow: 0 0 8px #00ff00; }
		.item { color: #f1c40f; }
		.ground { color: #333; }
		.creature { color: #e74c3c; font-weight: bold; }
		.unknown { color: #1a1a1a; }
	</style></head><body>`)

	// Display metadata
	fmt.Fprintf(w, `<div class="stats">
        <b>Player ID:</b> %d | 
        <b>Position:</b> (%d, %d, %d)
    </div>`, frame.Player.ID, pPos.X, pPos.Y, currentZ)

	fmt.Fprint(w, "<div class='map-container'>")

	// 4. Render the grid centered on the player
	for y := pPos.Y - radius; y <= pPos.Y+radius; y++ {
		for x := pPos.X - radius; x <= pPos.X+radius; x++ {

			currPos := domain.Position{X: x, Y: y, Z: currentZ}

			// A. Check if this is the player
			if currPos == pPos {
				fmt.Fprint(w, "<span class='tile player' title='YOU'>@</span>")
				continue
			}

			// B. Lookup tile in the Snapshot
			tile, ok := worldTiles[currPos]
			if !ok {
				// Not in the proxy's current map cache
				fmt.Fprint(w, "<span class='tile unknown'>.</span>")
				continue
			}

			// C. Determine what to draw (Creature > Item > Ground)
			char := "."
			class := "ground"
			title := fmt.Sprintf("Pos: %d, %d", x, y)

			if len(tile.Items) > 0 {
				char = "i"
				class = "item"
				title += fmt.Sprintf(" | %d items", len(tile.Items))
			}

			fmt.Fprintf(w, "<span class='tile %s' title='%s'>%s</span>", class, title, char)
		}
		// End of row
		fmt.Fprint(w, "<br>")
	}

	fmt.Fprint(w, "</div>")
	fmt.Fprint(w, "<p style='color:#666; font-size:12px;'>Centering logic: Player is always middle. '.' is ground, 'i' is item, 'C' is creature.</p>")
	fmt.Fprint(w, "</body></html>")
}
