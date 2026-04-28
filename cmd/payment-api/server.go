package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// ──────────────────────────────────────────────
//  Server Utilities
// ──────────────────────────────────────────────

func respond(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ──────────────────────────────────────────────
//  Entrypoint
// ──────────────────────────────────────────────

func run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/transaction", handleTransaction)
	mux.HandleFunc("/api/inject-failure", handleInjectFailure)

	log.Println("[Payment API] :8080 ✓")
	if err := http.ListenAndServe(":8080", cors(mux)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
