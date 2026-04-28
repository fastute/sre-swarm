// Package agents defines the SRE Swarm — a multi-agent system built on Google ADK
// that triages payment gateway incidents using Gemini 2.5 Flash.
package agents

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
)

// ──────────────────────────────────────────────
//  Exported Agent Handles (used by the server)
// ──────────────────────────────────────────────

var (
	Model   model.LLM
	Triage  agent.Agent // root — reads runbooks, decides who handles it
	Fixer   agent.Agent // restarts crashed K8s pods
	Alerter agent.Agent // sends Slack alerts to humans
)

// ──────────────────────────────────────────────
//  Bootstrap
// ──────────────────────────────────────────────

// Init wires up the Gemini model and all agents. Called once at startup.
func Init() {
	Model = initModel()
	Fixer = newFixer(Model)
	Alerter = newAlerter(Model)
	Triage = newTriage(Model, Fixer, Alerter)
	log.Println("[ADK] All agents initialized ✓")
}

// ──────────────────────────────────────────────
//  Model
// ──────────────────────────────────────────────

func initModel() model.LLM {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	location := os.Getenv("GOOGLE_CLOUD_LOCATION")
	if project == "" || location == "" {
		log.Fatal("Set GOOGLE_CLOUD_PROJECT and GOOGLE_CLOUD_LOCATION in .env")
	}

	m, err := gemini.NewModel(context.Background(), "gemini-2.5-flash", &genai.ClientConfig{
		Project:  project,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		log.Fatalf("Gemini init failed: %v", err)
	}
	log.Println("[ADK] Gemini 2.5 Flash via Vertex AI ✓")
	return m
}

// ──────────────────────────────────────────────
//  Agent Definitions
// ──────────────────────────────────────────────

func newFixer(m model.LLM) agent.Agent {
	return must(llmagent.New(llmagent.Config{
		Name:        "fixer",
		Model:       m,
		Description: "Restarts crashed services automatically.",
		Instruction: "You are the Fixer. Call restart_service with the service name. Report the result.",
		Tools:       []tool.Tool{restartServiceTool()},
	}))
}

func newAlerter(m model.LLM) agent.Agent {
	return must(llmagent.New(llmagent.Config{
		Name:        "alerter",
		Model:       m,
		Description: "Sends Slack alerts to human teams.",
		Instruction: "You are the Alerter. Call send_alert with the message. Report the result.",
		Tools:       []tool.Tool{sendAlertTool()},
	}))
}

func newTriage(m model.LLM, fixer, alerter agent.Agent) agent.Agent {
	return must(llmagent.New(llmagent.Config{
		Name:        "triage",
		Model:       m,
		Description: "Reads runbooks and decides who handles the incident.",
		Instruction: `You are the Triage agent for a payment gateway SRE team.

On every alert:
1. Call read_runbook with the EXACT signal (OOM_KILL, P99_SPIKE, AML_FLAG).
2. Follow the runbook, then delegate:
   • OOM_KILL / P99_SPIKE → transfer to fixer
   • AML_FLAG → report Human in the Loop required
   • Weekend → report Human in the Loop (defer to Monday)

Briefly explain your reasoning before acting.`,
		Tools:     []tool.Tool{readRunbookTool()},
		// If SubAgents doesn't exist on Config or expects different types, 
		// they can be handled separately, but let's try injecting them if llmagent supports it.
	}))
}

// ──────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────

func must(a agent.Agent, err error) agent.Agent {
	if err != nil {
		log.Fatalf("Agent creation failed: %v", err)
	}
	return a
}
