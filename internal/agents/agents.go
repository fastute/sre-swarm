package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
		Description: "Reads markdown runbook. Arg = exact Telemetry Signal (OOM_KILL, DEADLOCK, P99_SPIKE, AML_FLAG).",
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
		Instructions: `Role: Incident Commander. Input: Amount, Date, Signal.
1. Call ReadKnowledgeBase with EXACT Signal (OOM_KILL, DEADLOCK, P99_SPIKE, AML_FLAG). No guessing.
2. Follow rulebook protocol exactly.
3. Output ONE line, exact format:
"HANDOFF: Auto_Remediator | reason"
"HANDOFF: Human_in_Loop | reason"
"HANDOFF: Comms_Lead | reason"
"HANDOFF: Resilience_Engineer | reason"`,
		Tools:     []adklite.Tool{readKBTool},
		ModelHost: modelHost,
		ModelName: modelName,
	}

	// Auto Remediator — restarts crashed services
	AutoRemediatorAgent = &adklite.Agent{
		Role: "Auto_Remediator",
		Instructions: `Role: Auto Remediator. Call RestartKubernetesPod with given service name. Return result.`,
		Tools:     []adklite.Tool{restartK8sTool},
		ModelHost: modelHost,
		ModelName: modelName,
	}

	// Comms Lead — alerts humans
	CommsLeadAgent = &adklite.Agent{
		Role: "Comms_Lead",
		Instructions: `Role: Comms Lead. Call SendSlackAlert with given message. Return result.`,
		Tools:     []adklite.Tool{sendSlackAlertTool},
		ModelHost: modelHost,
		ModelName: modelName,
	}

	// Resilience Engineer — handles deadlock retries
	ResilienceEngineerAgent = &adklite.Agent{
		Role: "Resilience_Engineer",
		Instructions: `Role: Resilience Engineer. Call ExecuteExponentialBackoff. Return result.`,
		Tools:     []adklite.Tool{executeExponentialBackoffTool},
		ModelHost: modelHost,
		ModelName: modelName,
	}
}
