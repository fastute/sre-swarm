package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// ──────────────────────────────────────────────
//  Tool Schemas & Implementations
// ──────────────────────────────────────────────

type kbInput struct {
	Signal string `json:"signal"`
}

type kbOutput struct {
	Result string `json:"result"`
}

type podInput struct {
	ServiceName string `json:"service_name"`
}

type podOutput struct {
	Result string `json:"result"`
}

type slackInput struct {
	Message string `json:"message"`
}

type slackOutput struct {
	Result string `json:"result"`
}

// readKB loads the matching runbook from the knowledge/ directory.
func readKB(ctx tool.Context, in kbInput) (kbOutput, error) {
	for _, dir := range []string{"knowledge", "../../knowledge"} {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, f := range files {
			data, err := os.ReadFile(filepath.Join(dir, f.Name()))
			if err == nil && strings.Contains(string(data), in.Signal) {
				return kbOutput{Result: string(data)}, nil
			}
		}
	}
	return kbOutput{Result: "No runbook found for: " + in.Signal}, nil
}

func readRunbookTool() tool.Tool {
	t, err := functiontool.New(functiontool.Config{
		Name:        "read_runbook",
		Description: "Reads the SRE runbook for a signal.",
	}, readKB)
	if err != nil {
		panic(err)
	}
	return t
}

// restartPod simulates a Kubernetes pod restart.
func restartPod(ctx tool.Context, in podInput) (podOutput, error) {
	return podOutput{Result: fmt.Sprintf("✓ Restarted pod: %s", in.ServiceName)}, nil
}

func restartServiceTool() tool.Tool {
	t, err := functiontool.New(functiontool.Config{
		Name:        "restart_service",
		Description: "Restarts a crashed service by name.",
	}, restartPod)
	if err != nil {
		panic(err)
	}
	return t
}

// sendSlack simulates posting to Slack.
func sendSlack(ctx tool.Context, in slackInput) (slackOutput, error) {
	return slackOutput{Result: fmt.Sprintf("✓ Slack #sre-alerts: %s", in.Message)}, nil
}

func sendAlertTool() tool.Tool {
	t, err := functiontool.New(functiontool.Config{
		Name:        "send_alert",
		Description: "Posts to #sre-alerts Slack channel.",
	}, sendSlack)
	if err != nil {
		panic(err)
	}
	return t
}
