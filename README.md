# 🚨 PaymentGuard — Autonomous AI Agents for Payment Infrastructure

> _What if your on-call engineer was an AI that never sleeps, reads every runbook instantly, and coordinates a team of specialists in seconds?_

**PaymentGuard** is a team of AI agents that monitor a payment system, detect failures, look up the right fix from a knowledge base of plain Markdown files, and either fix the problem automatically or call the right human. Everything runs locally — no cloud APIs, no API keys — and streams live to a dashboard you can watch in real time.

---

## 🧠 Meet the Agents

Think of each agent as a virtual employee with a specific job:

| Agent | Role | Real-World Analogy |
| :--- | :--- | :--- |
| **Incident Commander** | The brain. Reads incoming alerts, looks up the rulebook, and decides what happens next. | A hospital triage nurse |
| **Auto Remediator** | The fixer. Restarts crashed services without waiting for a human. | IT helpdesk rebooting a frozen computer |
| **Resilience Engineer** | The strategist. Handles tricky problems like database deadlocks using smart retry logic. | A surgeon who needs a careful approach |
| **Comms Lead** | The messenger. Alerts humans via Slack or Jira when the AI can't (or shouldn't) act alone. | A receptionist paging the right doctor |

---

## 🗺️ Decision Flow

When something goes wrong, the **Incident Commander** reads the rulebook and delegates to the right specialist:

```mermaid
flowchart TD
    classDef start fill:#1e293b,stroke:#38bdf8,stroke-width:2px,color:#fff;
    classDef ai fill:#0f172a,stroke:#818cf8,stroke-width:2px,color:#bae6fd;
    classDef kb fill:#020617,stroke:#10b981,stroke-width:2px,color:#a7f3d0;
    classDef heal fill:#064e3b,stroke:#34d399,stroke-width:2px,color:#fff;
    classDef human fill:#7f1d1d,stroke:#f87171,stroke-width:2px,color:#fff;

    Start([Incoming Alert]):::start --> IC{Incident Commander}:::ai
    IC -->|Looks up rulebook| KB[(Knowledge Base)]:::kb
    KB -.->|Crash or latency spike| AR[Auto Remediator: Restart Service]:::heal
    KB -.->|Database deadlock| RE[Resilience Engineer: Retry with Backoff]:::ai
    KB -.->|Weekend or high value| HIL[Human in the Loop]:::human
    KB -.->|Compliance flag| CL[Comms Lead: Alert FinOps Team]:::ai
```

---

## ⚡ Quick Start

You need **Go**, **Ollama**, and the **Gemma 4 (26B)** model.

```bash
# 1. Pull the model (one-time setup)
ollama pull gemma4:26b

# 2. Start the mock payment backend
cd cmd/payment-api && go run main.go

# 3. In a new terminal — start the agent orchestrator
cd cmd/sre-swarm && go run main.go

# 4. Open the dashboard
open web/index.html
```

Pick a failure scenario from the dropdown, hit **Initiate Autonomous Swarm**, and watch the agents reason, delegate, and resolve — all streamed live.

---

## 🔧 How It Works

```
Alert arrives
    → Incident Commander reads the rulebook
        → Is it a simple crash?         → Auto Remediator fixes it
        → Is it a database deadlock?    → Resilience Engineer retries it
        → Is it a compliance issue?     → Comms Lead alerts a human
        → Is it a weekend or high-value? → Human is called directly
```

### The Knowledge Base (the "Rulebook")

Inside the project there's a `knowledge/` folder containing simple Markdown files. These are the AI's rulebook. They contain instructions like:

> *"If the alert is OOM_KILL, restart the settlement-worker service."*
>
> *"If the payment is over £10,000 and there's a compliance flag, a human must review it."*
>
> *"If the payment date falls on a weekend, defer processing to Monday."*

**You can change these files at any time and the AI's behavior changes instantly.** No programming. No retraining. Just edit the text.

### What is a "Handoff"?

When one agent finishes its job and passes the work to another, that's a handoff. For example, the Incident Commander hands off to the Auto Remediator to restart a crashed service. Every handoff is logged in the dashboard so you can see exactly what happened and why.

---

## 🖥️ What You See on the Dashboard

| Area | What it shows |
| :--- | :--- |
| **Transaction Injector** (left) | Where you simulate a problem — set the amount, date, and alert type |
| **Decision Banner** | A color-coded card showing the final outcome (🟢 fixed, 🔴 escalated, 🟣 notified) |
| **Swarm Logs** | A live feed of the AI thinking, looking up rulebooks, and executing fixes |
| **Show Thinking toggle** | Check/uncheck to show or hide the AI's internal reasoning |
| **Decision Flow Diagram** (right) | A visual map of all the paths the AI can take |

---

## 📁 Project Layout

```
├── cmd/
│   ├── payment-api/     # Mock payment backend
│   └── sre-swarm/       # Agent orchestrator server
├── internal/agents/     # Agent definitions & tools
├── knowledge/           # Markdown rulebooks (the AI's "brain")
├── pkg/adklite/         # Lightweight agent framework
├── web/                 # Dashboard UI
└── docs/                # Architecture diagrams & walkthrough
```

---

## 🎯 Why This Exists

Most AI agent demos are slow Python scripts calling cloud APIs. This one is different:

- **100% local** — runs on your machine with Ollama. No API keys, no cloud bills.
- **Fast** — written in Go. Talks directly to the AI model over HTTP.
- **Realistic** — models real payment failures (crashes, deadlocks, compliance flags).
- **Easy to customize** — edit a Markdown file to change how the AI makes decisions.
- **Safe by design** — high-value and compliance issues always get escalated to a human.

---

## 📚 Key Concepts

### What is an "AI Agent"?
A virtual employee that can **read** (look up documentation), **think** (analyze a situation), and **act** (execute a fix or send a message). PaymentGuard uses four agents working together as a coordinated team.

### What is a "Model"?
The AI brain powering the agents. We use **Gemma 4** by Google, running entirely on your own computer via a tool called **Ollama**. Your data never leaves your machine.

### What is the "Second Brain" / Librarian Pattern?
Instead of teaching the AI everything upfront, we give it access to a library of documents. When a problem occurs, the AI looks up the relevant document and follows the instructions.  
→ [Karpathy's AI-Driven Second Brain](https://techstrong.ai/features/karpathys-instructions-for-building-an-ai-driven-second-brain/)

### What is "Caveman Protocol"?
We talk to the AI using as few words as possible. Fewer words = faster responses. Instead of *"You are the Auto Remediator agent. You resolve infrastructure issues autonomously..."* we say *"Role: Auto Remediator. Call RestartKubernetesPod. Return result."*  
→ [CaVEMan on GitHub](https://github.com/cancerit/CaVEMan)

---

## ❓ FAQ

**Can I change what the AI does without coding?**  
Yes. Edit the Markdown files in `knowledge/`. The AI reads them fresh every time.

**Does this need the internet?**  
No. Everything runs on your local machine.

**What if the AI makes a wrong decision?**  
High-value transactions and compliance issues always get escalated to a human. The AI can only auto-fix things your rulebook explicitly marks as safe.

**How fast is it?**  
From alert to resolution, typically **5–15 seconds** depending on the scenario and your hardware.

---

## 📖 Further Reading

- [`docs/user_guide.md`](docs/user_guide.md) — Architecture diagrams and technical walkthrough
- [`docs/product_guide.md`](docs/product_guide.md) — Extended product guide with detailed explanations
