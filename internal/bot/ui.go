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
	"web/components/*.gohtml",
	"web/views/*.gohtml",
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
