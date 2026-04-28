package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"

	"github.com/sre-swarm/internal/agents"
)

// ──────────────────────────────────────────────
//  Constants
// ──────────────────────────────────────────────

const (
	appName = "sre-swarm"
	userID  = "dashboard"
)

// ──────────────────────────────────────────────
//  Handlers
// ──────────────────────────────────────────────

// POST /api/simulate — trigger the swarm on an alert
func handleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
		Date   string  `json:"date"`
		Signal string  `json:"signal"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("\n══ SWARM TRIGGERED ══ £%.2f | %s | %s", req.Amount, req.Date, req.Signal)
	broadcast("Incident Commander is analyzing...")

	// Build prompt
	ctx := context.Background()
	now := time.Now().Format("Monday, 15:04:05")
	prompt := fmt.Sprintf(
		"Transaction £%.2f GBP on %s (System Time: %s). Signal: %s. Investigate and resolve per SRE policy.",
		req.Amount, req.Date, now, req.Signal,
	)

	// Create session service + session
	sessSvc := session.InMemoryService()
	createResp, err := sessSvc.Create(ctx, &session.CreateRequest{
		AppName: appName,
		UserID:  userID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("session create: %v", err), http.StatusInternalServerError)
		return
	}

	// Create runner
	agentRunner, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          agents.Triage,
		SessionService: sessSvc,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("runner init: %v", err), http.StatusInternalServerError)
		return
	}

	// Run the swarm and stream events
	res := struct {
		Status       string `json:"status"`
		Handoff      string `json:"handoff"`
		Target       string `json:"target"`
		IsAutonomous bool   `json:"isAutonomous"`
	}{Status: "Resolution Complete"}

	msg := genai.NewContentFromText(prompt, genai.RoleUser)

	for event, runErr := range agentRunner.Run(
		ctx,
		createResp.Session.UserID(),
		createResp.Session.ID(),
		msg,
		agent.RunConfig{},
	) {
		if runErr != nil {
			emit("Error", runErr.Error())
			continue
		}
		if event == nil || event.Content == nil {
			continue
		}

		author := event.Author
		if author == "" {
			author = "incident_commander"
		}

		for _, part := range event.Content.Parts {
			if part == nil {
				continue
			}
			switch {
			case part.Text != "":
				text := part.Text
				emit(author, text)
				switch author {
				case "fixer":
					res.Handoff, res.Target, res.IsAutonomous = "Fixer", text, true
				case "alerter":
					res.Handoff, res.Target, res.IsAutonomous = "Alerter", text, false
				default:
					if isHumanLoop(text) {
						res.Handoff, res.Target, res.IsAutonomous = "Human_in_Loop", text, false
					}
				}
			case part.FunctionCall != nil:
				emit("Tool", fmt.Sprintf("%s(%v)", part.FunctionCall.Name, part.FunctionCall.Args))
			case part.FunctionResponse != nil:
				emit("Result", fmt.Sprintf("%v", part.FunctionResponse.Response))
			}
		}
	}

	if res.Handoff == "" {
		res.Handoff, res.Target, res.IsAutonomous = "Triage", "Triage completed", true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GET /api/stream-logs — SSE endpoint for the dashboard
func handleStreamLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case msg := <-logStream:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// POST /api/log-external — receives logs from the Payment API
func handleExternalLog(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Msg string `json:"msg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
		broadcast("[PAYMENT] " + req.Msg)
	}
}

// ──────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────

func isHumanLoop(text string) bool {
	lower := strings.ToLower(text)
	for _, kw := range []string{"human in the loop", "compliance review", "defer to monday"} {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}
