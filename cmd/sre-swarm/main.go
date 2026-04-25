package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sre-swarm/internal/agents"
)

var logStream = make(chan string, 100)

func broadcastLog(msg string) {
	select {
	case logStream <- msg:
	default:
		// Channel full, drop message
	}
}

func streamLogsHandler(w http.ResponseWriter, r *http.Request) {
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
			// Replace newlines with HTML breaks to avoid breaking the SSE format
			cleanMsg := strings.ReplaceAll(msg, "\n", "<br>")
			fmt.Fprintf(w, "data: %s\n\n", cleanMsg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// simulateIncidentHandler handles requests from the Web UI to trigger the swarm
func simulateIncidentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	log.Printf("\n==================================================")
	log.Printf("AUTONOMOUS DISCOVERY INITIATED")
	log.Printf("Transaction: £%.2f GBP | Date: %s | Signal: %s", req.Amount, req.Date, req.Signal)
	log.Printf("==================================================")

	// 1. Triage Agent Discovery & Investigation
	// We prompt the agent to hallucinate a plausible failure based on the metadata, the signal hint, and its runbooks.
	currentTime := time.Now().Format("Monday, 15:04:05")
	triagePrompt := fmt.Sprintf(`A transaction of £%.2f GBP was processed on date %s (Current System Time: %s).
Primary Telemetry Signal: %s

1. If the Signal is not AUTO, prioritize brainstorming a failure scenario related to that signal.
2. Otherwise, brainstorm any plausible SRE failure (Settlement, Compliance, or Latency).
3. Use your tools to investigate and resolve it based on company policy.`, req.Amount, req.Date, currentTime, req.Signal)
	
	msgStart := "Strategic_Triage Agent is analyzing the telemetry..."
	broadcastLog(msgStart)
	
	triageResult := agents.Orchestrator.Run(agents.IncidentCommanderAgent, triagePrompt, broadcastLog)

	// Prepare a structured response for the UI
	var response struct {
		Status      string `json:"status"`
		Handoff     string `json:"handoff"`
		Target      string `json:"target"`
		IsAutonomous bool   `json:"isAutonomous"`
	}
	response.Status = "Resolution Complete"

	// Step 2: Handoff Parsing
	if strings.Contains(triageResult, "HANDOFF: Auto_Remediator") {
		parts := strings.Split(triageResult, "|")
		serviceName := "Unknown"
		if len(parts) > 1 {
			serviceName = strings.TrimSpace(parts[1])
		}
		
		msgHealer := fmt.Sprintf("[Handoff] Incident_Commander -> Auto_Remediator: Resolving %s", serviceName)
		broadcastLog(msgHealer)

		healerPrompt := fmt.Sprintf("Restart the following service: %s", serviceName)
		agents.Orchestrator.Run(agents.AutoRemediatorAgent, healerPrompt, broadcastLog)

		response.Handoff = "Auto_Remediator"
		response.Target = serviceName
		response.IsAutonomous = true

	} else if strings.Contains(triageResult, "HANDOFF: Human_in_Loop") {
		parts := strings.Split(triageResult, "|")
		reason := "Architectural Risk detected."
		if len(parts) > 1 {
			reason = strings.TrimSpace(parts[1])
		}
		
		msgHIL := fmt.Sprintf("[Handoff] Incident_Commander -> Human_in_Loop: %s", reason)
		broadcastLog(msgHIL)
		
		response.Handoff = "Human_in_Loop"
		response.Target = reason
		response.IsAutonomous = false

	} else if strings.Contains(triageResult, "HANDOFF: Comms_Lead") {
		parts := strings.Split(triageResult, "|")
		reason := "Notification required."
		if len(parts) > 1 {
			reason = strings.TrimSpace(parts[1])
		}

		msgNotif := fmt.Sprintf("[Handoff] Incident_Commander -> Comms_Lead: %s", reason)
		broadcastLog(msgNotif)
		
		notificationPrompt := fmt.Sprintf("Please notify the relevant team via Slack regarding: %s", reason)
		agents.Orchestrator.Run(agents.CommsLeadAgent, notificationPrompt, broadcastLog)

		response.Handoff = "Comms_Lead"
		response.Target = reason
		response.IsAutonomous = false
	} else if strings.Contains(triageResult, "HANDOFF: Resilience_Engineer") {
		parts := strings.Split(triageResult, "|")
		reason := "Initiating Exponential Backoff."
		if len(parts) > 1 {
			reason = strings.TrimSpace(parts[1])
		}

		msgBackoff := fmt.Sprintf("[Handoff] Incident_Commander -> Resilience_Engineer: %s", reason)
		broadcastLog(msgBackoff)
		
		backoffPrompt := fmt.Sprintf("Execute the backoff tool to resolve this deadlock: %s", reason)
		agents.Orchestrator.Run(agents.ResilienceEngineerAgent, backoffPrompt, broadcastLog)

		response.Handoff = "Resilience_Engineer"
		response.Target = "Database Transaction Retried"
		response.IsAutonomous = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	log.Println("Initializing ADK Swarm Agents...")
	agents.Init()

	mux := http.NewServeMux()
	
	// Serve Web UI
	fs := http.FileServer(http.Dir("../../web"))
	mux.Handle("/", fs)
	
	mux.HandleFunc("/api/simulate", simulateIncidentHandler)
	mux.HandleFunc("/api/stream-logs", streamLogsHandler)
	mux.HandleFunc("/api/log-external", externalLogHandler)

	log.Println("[SRE Swarm Orchestrator] UI & API listening on :9090")
	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func externalLogHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Msg string `json:"msg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
		broadcastLog("[PAYMENT] " + req.Msg)
	}
}
