package bot

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed web/*
var webAssets embed.FS
var templates = template.Must(template.ParseFS(webAssets,
	"web/index.html",
	"web/components/*.html",
	"web/views/*.html",
))

func (b *Bot) loopWebUI() {
	mux := http.NewServeMux()

	// Main Page
	mux.HandleFunc("/", b.handleDashboard)

	// 2. The Detail Views (Master-Detail Pattern)
	// This matches <button hx-get="/views/fishing" ...> in index.html
	mux.HandleFunc("/views/fishing", b.handleFishingView)

	// 3. API Endpoints (Fragments)
	// This matches <button hx-post="/api/toggle-fishing" ...>
	mux.HandleFunc("/api/toggle-fishing", b.handleToggleFishing)

	// API Endpoints
	mux.HandleFunc("/api/stats", b.handleGetStats) // Simple text return

	fmt.Println("[UI] Dashboard live at http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func (b *Bot) handleDashboard(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func (b *Bot) handleFishingView(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Enabled bool
	}{
		Enabled: b.fishingEnabled,
	}
	// Return only the fishing-view partial
	templates.ExecuteTemplate(w, "fishing-view", data)
}

func (b *Bot) handleToggleFishing(w http.ResponseWriter, r *http.Request) {
	b.fishingEnabled = !b.fishingEnabled
	data := struct {
		Enabled bool
	}{
		Enabled: b.fishingEnabled,
	}
	templates.ExecuteTemplate(w, "fishing-toggle", data)
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
