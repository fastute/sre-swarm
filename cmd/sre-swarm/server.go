package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/sre-swarm/internal/agents"
)

// ──────────────────────────────────────────────
//  Server
// ──────────────────────────────────────────────

func main() {
	// Load .env
	if err := godotenv.Load("../../.env"); err != nil {
		godotenv.Load(".env")
	}

	log.Println("Initializing ADK Swarm...")
	agents.Init()

	mux := http.NewServeMux()
	// Serve Web UI (check both cmd/sre-swarm and root context)
	webPath := "../../web"
	if _, err := os.Stat(webPath); os.IsNotExist(err) {
		webPath = "./web"
	}
	mux.Handle("/", http.FileServer(http.Dir(webPath)))
	mux.HandleFunc("/api/simulate", handleSimulate)
	mux.HandleFunc("/api/stream-logs", handleStreamLogs)
	mux.HandleFunc("/api/log-external", handleExternalLog)

	log.Println("[SRE Swarm] :9090 ✓")
	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
