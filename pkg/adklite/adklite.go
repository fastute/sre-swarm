package adklite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Agent represents a specialized persona in the swarm
type Agent struct {
	Role         string
	Instructions string
	Tools        []Tool
	ModelHost    string
	ModelName    string
}

// Tool represents a capability the agent can use
type Tool struct {
	Name        string
	Description string
	Execute     func(args string) string
}

// ADKLite orchestrates the agents
type ADKLite struct {
}

func New() *ADKLite {
	return &ADKLite{}
}

// Run executes an agent against a prompt using Ollama, emitting logs to the optional logCb.
// It supports a multi-turn tool-use loop (up to 5 turns).
func (s *ADKLite) Run(agent *Agent, prompt string, logCb func(string)) string {
	// Build the initial system prompt
	systemPrompt := agent.Instructions + "\n\nYou have access to the following tools:\n"
	for _, t := range agent.Tools {
		systemPrompt += fmt.Sprintf("- %s: %s\n", t.Name, t.Description)
	}
	systemPrompt += `
CRITICAL INSTRUCTION:
You MUST format your decisions exactly as shown in these examples. DO NOT use conversational text.
Example 1 (Tool): CALL_TOOL: ReadKnowledgeBase | PAYMENT_SETTLEMENT_FAILED
Example 2 (Handoff): HANDOFF: Auto_Healer | payment-service
`

	// Track conversation history for the multi-turn loop
	messages := []map[string]string{
		{"role": "system", "content": systemPrompt},
		{"role": "user", "content": prompt},
	}

	for turn := 0; turn < 5; turn++ {
		// Call Ollama
		reqBody := map[string]interface{}{
			"model":    agent.ModelName,
			"messages": messages,
			"stream":   false,
		}

		jsonData, _ := json.Marshal(reqBody)
		resp, err := http.Post(fmt.Sprintf("%s/api/chat", agent.ModelHost), "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			errMsg := fmt.Sprintf("Error connecting to Ollama: %v", err)
			log.Println(errMsg)
			if logCb != nil { logCb(errMsg) }
			return errMsg
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf("Ollama API Error [%d]: %s", resp.StatusCode, string(body))
			log.Println(errMsg)
			if logCb != nil { logCb(errMsg) }
			return errMsg
		}

		var ollamaResp struct {
			Message struct {
				Role     string `json:"role"`
				Content  string `json:"content"`
				Thinking string `json:"thinking"`
			} `json:"message"`
		}
		json.Unmarshal(body, &ollamaResp)

		// Combine thinking and content
		content := strings.TrimSpace(ollamaResp.Message.Thinking + "\n" + ollamaResp.Message.Content)
		msgReasoning := fmt.Sprintf("[Reasoning] %s", content)
		log.Println(msgReasoning)
		if logCb != nil { logCb(msgReasoning) }

		// Add assistant's response to history
		messages = append(messages, map[string]string{"role": "assistant", "content": content})

		// Check for tool usage
		if strings.Contains(content, "CALL_TOOL:") {
			parts := strings.SplitN(content, "CALL_TOOL:", 2)
			toolCall := strings.TrimSpace(parts[1])
			toolParts := strings.SplitN(toolCall, "|", 2)
			if len(toolParts) == 2 {
				toolName := strings.TrimSpace(toolParts[0])
				toolArgs := strings.TrimSpace(toolParts[1])

				// Find and execute tool
				found := false
				for _, t := range agent.Tools {
					if t.Name == toolName {
						msgExec := fmt.Sprintf("[Executing Tool] %s(%s)", t.Name, toolArgs)
						log.Println(msgExec)
						if logCb != nil { logCb(msgExec) }

						result := t.Execute(toolArgs)
						msgResult := fmt.Sprintf("[Tool Result] %s", result)
						log.Println(msgResult)
						if logCb != nil { logCb(msgResult) }

						// Append tool result to history and continue loop
						messages = append(messages, map[string]string{"role": "user", "content": fmt.Sprintf("Tool %s returned: %s", toolName, result)})
						found = true
						break
					}
				}
				if found {
					continue
				}
			}
		}

		// If no tool call was found, this is the final answer or handoff
		return content
	}

	return "Max tool execution depth reached"
}
