package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Transaction struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	ErrorCode string    `json:"errorCode,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

var (
	transactions = make(map[string]*Transaction)
	mu           sync.Mutex
)

func sendToSwarm(msg string) {
	body, _ := json.Marshal(map[string]string{"msg": msg})
	http.Post("http://localhost:9090/api/log-external", "application/json", bytes.NewBuffer(body))
}

func transactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	txID := fmt.Sprintf("TX-%d", rand.Int63())
	msg := fmt.Sprintf("Processing Transaction %s...", txID)
	log.Println(msg)
	sendToSwarm(msg)

	time.Sleep(500 * time.Millisecond)

	// In a real demo, this is where the failure would happen
	msgDone := fmt.Sprintf("Transaction %s completed successfully", txID)
	log.Println(msgDone)
	sendToSwarm(msgDone)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":     txID,
		"status": "success",
	})
}

// InjectFailure handler forces a transaction to fail to simulate an incident
func injectFailureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	msg := fmt.Sprintf("CRITICAL FAILURE: %s in Transaction %s", req.ErrorCode, txID)
	log.Println(msg)
	sendToSwarm(msg)

	mu.Lock()
	tx := &Transaction{
		ID:        txID,
		Amount:    1000.50,
		Status:    "FAILED",
		ErrorCode: req.ErrorCode,
		CreatedAt: time.Now(),
	}
	transactions[tx.ID] = tx
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":     txID,
		"status": "error",
		"error":  req.ErrorCode,
	})
}

func getStatusHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var txList []*Transaction
	for _, tx := range transactions {
		txList = append(txList, tx)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txList)
}

// enableCORS is a simple middleware to allow cross-origin requests from the UI
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/inject-failure", injectFailureHandler)
	mux.HandleFunc("/api/status", getStatusHandler)

	log.Println("[Mock Payment API] Listening on :8080 (CORS enabled)")
	if err := http.ListenAndServe(":8080", enableCORS(mux)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
