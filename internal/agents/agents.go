package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sre-swarm/pkg/adklite"
)

var (
	Orchestrator            *adklite.ADKLite
	IncidentCommanderAgent  *adklite.Agent
	AutoRemediatorAgent     *adklite.Agent
	CommsLeadAgent          *adklite.Agent
	ResilienceEngineerAgent *adklite.Agent
)

// Init configures the ADK-style swarm orchestrator and its specialized agents.
// It maps the topology for the decentralized physical machines.
func Init() {
	// 1. Setup Swarm Orchestrator
	Orchestrator = adklite.New()

	// We will map all agents to the local machine for testing to ensure the demo works seamlessly.
	modelHost := "http://localhost:11434" 
	modelName := "gemma4:26b"

	// 2. Define Tools (Capabilities)
	readKBTool := adklite.Tool{
		Name:        "ReadKnowledgeBase",
		Description: "Reads the contents of a markdown runbook. You MUST provide the exact Telemetry Signal as the argument (e.g., OOM_KILL, DEADLOCK).",
		Execute: func(args string) string {
			// Karpathy's Librarian Approach: Load markdown into context dynamically
			// Update the relative path since we are now inside internal/agents
			baseDir := "../../knowledge"
			files, _ := os.ReadDir(baseDir)

			for _, f := range files {
				content, err := os.ReadFile(filepath.Join(baseDir, f.Name()))
				if err == nil && strings.Contains(string(content), args) {
					return string(content)
				}
			}
			return "No runbook found for that error."
		},
	}

	restartK8sTool := adklite.Tool{
		Name:        "RestartKubernetesPod",
		Description: "Restarts a service in Kubernetes. Provide the service name as an argument.",
		Execute: func(args string) string {
			return fmt.Sprintf("Successfully restarted pod for %s", args)
		},
	}

	sendSlackAlertTool := adklite.Tool{
		Name:        "SendSlackAlert",
		Description: "Sends an alert to a Slack channel. Provide the message as an argument.",
		Execute: func(args string) string {
			return fmt.Sprintf("Message sent to #sre-alerts: %s", args)
		},
	}

	executeExponentialBackoffTool := adklite.Tool{
		Name:        "ExecuteExponentialBackoff",
		Description: "Executes an exponential backoff sequence for a deadlocked transaction.",
		Execute: func(args string) string {
			// Simulate backoff sequence
			time.Sleep(1 * time.Second)
			return "Attempt 1: Failed (Lock acquired by another process). Attempt 2: Failed (Lock timeout). Attempt 3: Succeeded after 2.5s backoff."
		},
	}

	// 3. Define Agents (Personas)

	// Incident Commander Agent (Runs on MacBook)
	// Responsible for bridging the Information Vacuum using the Second Brain.
	IncidentCommanderAgent = &adklite.Agent{
		Role: "Incident_Commander",
		Instructions: `You are the Incident Commander. You receive raw transaction telemetry (Amount, Date, Signal).
STEP 1 (Discovery): Based on the transaction amount, date, and signal, brainstorm a plausible failure scenario.
STEP 2 (Investigation): Use the ReadKnowledgeBase tool to find the relevant rulebook. You MUST pass the EXACT Telemetry Signal (e.g., OOM_KILL, DEADLOCK) to the tool. Do NOT guess error codes.
STEP 3 (Resolution): Decide the path strictly based on the rulebook findings.
You MUST output your final decision in exactly one of these formats:
1. "HANDOFF: Auto_Remediator | <reason>"
2. "HANDOFF: Human_in_Loop | <reason>"
3. "HANDOFF: Comms_Lead | <reason>"
4. "HANDOFF: Resilience_Engineer | <reason>"`,
		Tools:     []adklite.Tool{readKBTool},
		ModelHost: modelHost,
		ModelName: modelName,
	}

	// Auto Remediator Agent
	// Responsible for asynchronous infrastructure remediation.
	AutoRemediatorAgent = &adklite.Agent{
		Role: "Auto_Remediator",
		Instructions: `You are the Auto Remediator agent. You resolve infrastructure issues autonomously without human intervention.
Use the RestartKubernetesPod tool on the service name provided to you. Return a success message when done.`,
		Tools:     []adklite.Tool{restartK8sTool},
		ModelHost: modelHost,
		ModelName: modelName,
	}

	// Comms Lead Agent
	// Responsible for outward communications.
	CommsLeadAgent = &adklite.Agent{
		Role: "Comms_Lead",
		Instructions: `You are the Comms Lead. You receive handoffs from the Incident Commander when a team must be notified.
You must use the SendSlackAlert tool to notify the relevant team.`,
		Tools:     []adklite.Tool{sendSlackAlertTool},
		ModelHost: modelHost,
		ModelName: modelName,
	}

	// Resilience Engineer Agent
	// Responsible for orchestrating retries
	ResilienceEngineerAgent = &adklite.Agent{
		Role: "Resilience_Engineer",
		Instructions: `You are the Resilience Engineer. You handle transactional deadlocks by executing exponential backoffs.
You MUST use the ExecuteExponentialBackoff tool to simulate resolving the transaction block, and return the result.`,
		Tools:     []adklite.Tool{executeExponentialBackoffTool},
		ModelHost: modelHost,
		ModelName: modelName,
	}
}
