package bot

import (
	"fmt"
	"log"
	"net/http"
)

func (b *Bot) handleDashboard(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func (b *Bot) handleFishingView(w http.ResponseWriter, r *http.Request) {
	data := struct {
		FishingData interface{}
	}{
		FishingData: struct {
			Endpoint string
			Enabled  bool
		}{
			Endpoint: "/api/toggle-fishing",
			Enabled:  b.fishingEnabled,
		},
	}
	templates.ExecuteTemplate(w, "fishing-view", data)
}

func (b *Bot) handleToggleFishing(w http.ResponseWriter, r *http.Request) {
	// 1. Toggle the logic
	b.fishingEnabled = !b.fishingEnabled

	// 2. Prepare the FULL data the generic toggle needs
	// If you miss 'Endpoint', the next click won't work!
	data := struct {
		Name     string
		Endpoint string
		Enabled  bool
	}{
		Name:     "Auto Fishing",
		Endpoint: "/api/toggle-fishing",
		Enabled:  b.fishingEnabled,
	}

	// 3. Use the correct template name "toggle"
	err := templates.ExecuteTemplate(w, "toggle", data)
	if err != nil {
		// This log is crucial. If the template fails,
		// you will see why in your terminal.
		log.Printf("[UI Error] Failed to render toggle: %v", err)
	}
}

func (b *Bot) handleGetStats(w http.ResponseWriter, r *http.Request) {
	frame := b.state.CaptureFrame()

	data := struct {
		Name string
		Pos  string
	}{
		Name: frame.Player.Name,
		Pos:  fmt.Sprintf("%d, %d, %d", frame.Player.Pos.X, frame.Player.Pos.Y, frame.Player.Pos.Z),
	}

	templates.ExecuteTemplate(w, "stats", data)
}
