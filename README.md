# Autonomous SRE Swarm Dashboard

An agentic AI demonstration of an SRE (Site Reliability Engineering) Swarm. This project uses multiple specialized AI agents to triage, remediate, and escalate simulated infrastructure and compliance failures using local LLMs (Ollama) and a custom Go orchestrator.

## System Architecture

The swarm uses a decentralized "Intelligence Mesh" pattern, where specialized agents act as a coordinated team:

*   **Incident_Commander:** Analyzes telemetry and consults the Markdown rulebooks using a RAG (Retrieval-Augmented Generation) "Librarian Pattern".
*   **Auto_Remediator:** Autonomously executes infrastructure commands (e.g., restarting Kubernetes pods).
*   **Resilience_Engineer:** Handles architectural risks like database deadlocks by executing exponential backoff retries.
*   **Comms_Lead:** Escalates compliance and severe issues to human teams via Slack and Jira.

## How to Start the Project

### Prerequisites
*   [Go](https://golang.org/) installed.
*   [Ollama](https://ollama.com/) installed and running locally.
*   The `gemma4:26b` model pulled via Ollama (`ollama run gemma4:26b`).

### Startup Instructions

1.  **Start Local LLM (Ollama):**
    Ensure Ollama is running on your machine. If you are distributing agents across multiple machines, you must bind Ollama to the public network interface:
    ```bash
    OLLAMA_HOST=0.0.0.0 ollama serve
    ```
    *(If running entirely on one machine, standard `ollama serve` is fine).*

2.  **Start the Mock Payment API:**
    Open a terminal and run the mock backend:
    ```bash
    cd cmd/payment-api
    go run main.go
    ```

3.  **Start the SRE Swarm Orchestrator:**
    Open a second terminal and start the main swarm orchestrator:
    ```bash
    cd cmd/sre-swarm
    go run main.go
    ```

4.  **Open the Dashboard:**
    Simply open the `web/index.html` file in your preferred web browser.

## Customizing the Swarm Logic

The `Incident_Commander` reads directly from the `knowledge/` directory. You can dynamically alter the swarm's behavior by editing the Markdown rulebooks (`settlement_failures.md` and `compliance_failures.md`). The AI will adapt to your new rules in real-time on the next transaction!

## How it was built so fast
Rather than using massive, bloated node/python frameworks, the core of this system is a custom, lightweight orchestrator (`internal/agents/agents.go`) that implements raw primitives:
- **Agents** hold instructions, host IPs, and a list of Tools.
- **Context/Knowledge** is dynamically loaded by the agents reading your markdown files.
- **Direct LLM Communication**: The `Swarm.Run` method directly formats the agent's instructions and the available tools into a JSON payload and sends it to the `/api/chat` Ollama endpoint. This removes overhead and makes the code blazingly fast.
