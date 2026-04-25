# 🚨 SRE Swarm — AI Agents That Fix Your Infrastructure

> _What if your on-call engineer was an AI that never sleeps, reads every runbook instantly, and coordinates a team of specialists in seconds?_

**SRE Swarm** is a live demo of exactly that. It's a team of AI agents — built in Go, powered by local LLMs — that detect payment system failures, look up the right fix from your company knowledge base, and either fix it automatically or escalate it to the right human. All in real-time, streamed to a slick dashboard.

---

## 🧠 Meet the Agents

| Agent | What it does |
| :--- | :--- |
| **Incident Commander** | The brain. Reads telemetry, consults the rulebooks, and decides what happens next. |
| **Auto Remediator** | The hands. Restarts crashed services autonomously. |
| **Resilience Engineer** | The strategist. Handles tricky problems like database deadlocks with smart retry logic. |
| **Comms Lead** | The messenger. Alerts humans via Slack/Jira when the AI can't (or shouldn't) act alone. |

---

## ⚡ Quick Start

You need **Go**, **Ollama**, and the `gemma4:26b` model.

```bash
# 1. Pull the model (one-time setup)
ollama pull gemma4:26b

# 2. Start the mock payment backend
cd cmd/payment-api && go run main.go

# 3. In a new terminal — start the swarm
cd cmd/sre-swarm && go run main.go

# 4. Open the dashboard
open web/index.html
```

That's it. Pick a failure scenario from the dropdown, hit **Initiate Autonomous Swarm**, and watch the agents reason, delegate, and resolve.

---

## 🔧 How It Works

```
You trigger a failure  →  Incident Commander reads the rulebook
                          →  Decides: auto-fix? escalate? retry?
                              →  Hands off to the right agent
                                  →  You watch it happen live
```

The secret sauce is the **knowledge base** — a set of plain Markdown files in `knowledge/`. The AI reads these files at runtime to decide what to do. Change a rule in the Markdown, and the AI's behavior changes instantly. No retraining, no redeploying.

---

## 📁 Project Layout

```
├── cmd/
│   ├── payment-api/     # Mock payment backend
│   └── sre-swarm/       # Main orchestrator server
├── internal/agents/     # Agent definitions & tools
├── knowledge/           # Markdown rulebooks (the AI's "brain")
├── pkg/adklite/         # Lightweight agent framework
├── web/                 # Dashboard UI
└── docs/                # Architecture diagrams & user guide
```

---

## 🎯 Why This Exists

Most AI agent demos are slow Python scripts talking to cloud APIs. This one is different:

- **100% local** — runs on your machine with Ollama. No API keys, no cloud bills.
- **Blazingly fast** — written in Go, not Python. Direct HTTP to Ollama's `/api/chat`.
- **Actually useful** — models real SRE workflows (OOM kills, deadlocks, compliance flags).
- **Easy to customize** — edit a Markdown file to change the AI's decision logic.

---

## 📖 Further Reading

See [`docs/user_guide.md`](docs/user_guide.md) for detailed architecture diagrams and a full technical walkthrough.
