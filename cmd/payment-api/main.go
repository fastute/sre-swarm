// Package main runs a mock payment gateway that simulates transactions
// and injects failures for the SRE Swarm to triage.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// ──────────────────────────────────────────────
//  Swarm Integration
// ──────────────────────────────────────────────

// notify forwards a log line to the SRE Swarm's SSE stream.
func notify(msg string) {
	body, _ := json.Marshal(map[string]string{"msg": msg})
	http.Post("http://localhost:9090/api/log-external", "application/json", bytes.NewBuffer(body))
}

// ──────────────────────────────────────────────
//  Handlers
// ──────────────────────────────────────────────

// POST /api/transaction — simulate a successful payment
func handleTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	txID := fmt.Sprintf("TX-%d", rand.Int63())
	notify(fmt.Sprintf("Processing %s...", txID))

	time.Sleep(500 * time.Millisecond) // simulate latency

	notify(fmt.Sprintf("✓ %s completed", txID))
	respond(w, map[string]string{"id": txID, "status": "success"})
}

// POST /api/inject-failure — force a transaction failure to trigger the swarm
func handleInjectFailure(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ErrorCode string `json:"errorCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	txID := fmt.Sprintf("TX-%d", rand.Int63())
	notify(fmt.Sprintf("⚠ FAILURE: %s on %s", req.ErrorCode, txID))

	respond(w, map[string]string{"id": txID, "status": "error", "error": req.ErrorCode})
}

func main() { run() }
