# PaymentGuard — Product Guide

> A plain-English explanation of what PaymentGuard does, how it works, and why it matters. No technical background required.

---

## What is PaymentGuard?

Imagine you run a company that processes thousands of payments every day. Sometimes things go wrong — a server crashes, a database gets stuck, or a suspicious transaction trips a compliance flag. When that happens, someone needs to figure out what went wrong, look up the correct procedure, and either fix it or call the right person.

Today, that "someone" is a human engineer who gets paged at 2 AM, spends 20 minutes reading documentation, and then runs a fix. **PaymentGuard replaces that process with a team of AI agents that do it in seconds.**

---

## What is an "AI Agent"?

Think of an AI agent like a very capable virtual employee. It can:

- **Read** — Look up information from your company's internal documentation
- **Think** — Analyze a situation and decide what to do
- **Act** — Execute a fix (like restarting a crashed server) or send a message to a human

A single agent is useful. But PaymentGuard uses **four agents working together as a team** — each one specializing in a different job. We call this a "swarm."

---

## Meet the Team

### 🧠 Incident Commander
**Job:** The team leader. When something goes wrong, this agent is the first to respond.

**What it does:**
1. Receives the alert (e.g., "Server crashed" or "Suspicious transaction detected")
2. Opens the company's internal rulebook to understand the correct procedure
3. Decides which specialist to call in

**Real-world analogy:** A hospital triage nurse who assesses a patient and decides whether they need surgery, medication, or just observation.

---

### 🔧 Auto Remediator
**Job:** The fixer. Handles problems that have a known, safe, automated solution.

**What it does:**
- Restarts crashed services
- Clears stuck processes

**When it's called in:** The Incident Commander looks up the rulebook and finds that the problem (e.g., a memory crash) has a simple fix: restart the service. No human needed.

**Real-world analogy:** An IT helpdesk that reboots a frozen computer — it's a well-known fix that doesn't need a senior engineer.

---

### 🔄 Resilience Engineer
**Job:** The strategist. Handles problems that are too risky to fix with a simple restart.

**What it does:**
- When a database gets "stuck" (a deadlock), this agent doesn't just restart everything — that could make things worse
- Instead, it carefully rolls back the failed operation and retries it with increasing wait times (called "exponential backoff")

**When it's called in:** The Incident Commander recognizes that a restart would be dangerous and delegates to this specialist instead.

**Real-world analogy:** A surgeon who can't just "reboot" the patient — they need a careful, step-by-step approach.

---

### 📢 Comms Lead
**Job:** The messenger. Alerts the right humans when AI alone isn't enough.

**What it does:**
- Sends Slack messages to the on-call team
- Creates Jira tickets for compliance reviews

**When it's called in:** The problem involves money (high-value transactions), legal requirements (compliance flags), or happens at a time when automated fixes are restricted (weekends).

**Real-world analogy:** A receptionist who pages the right doctor when a case is beyond the nurse's authority.

---

## How Does It Decide What To Do?

This is the most important part. The AI doesn't just guess — it follows your company's rules.

### The Knowledge Base (the "Rulebook")

Inside the project, there's a folder called `knowledge/` containing simple text files. These files are the AI's rulebook. They contain instructions like:

> *"If the alert is OOM_KILL, restart the settlement-worker service."*
>
> *"If the payment is over £10,000 and there's a compliance flag, a human must review it."*
>
> *"If the payment date falls on a weekend, defer processing to Monday."*

**The key insight:** You can change these text files at any time, and the AI's behavior changes instantly. No programming required. No retraining. Just edit the text.

### The Decision Flow

Here's a simplified view of how every alert is processed:

```
Alert arrives
    → Incident Commander reads the rulebook
        → Is it a simple crash?        → Auto Remediator fixes it
        → Is it a database deadlock?   → Resilience Engineer retries it
        → Is it a compliance issue?    → Comms Lead alerts a human
        → Is it a weekend or high-value? → Human is called directly
```

---

## What Happens on the Dashboard?

When you open PaymentGuard in your browser, you see a live dashboard with three main areas:

### 1. Transaction Injector (left side)
This is where you simulate a problem. You enter:
- **Amount** — How much the payment is worth (e.g., £12,500)
- **Date** — When the payment is being processed
- **Alert type** — What went wrong (e.g., server crash, deadlock, compliance flag)

Then you press **"Initiate Autonomous Swarm"** and watch the AI respond.

### 2. Swarm Orchestrator Logs (below the injector)
This is a live feed of everything the AI is doing in real time:
- Reading the rulebook
- Thinking through its options
- Handing off to a specialist agent
- Executing a fix

There's a **"Show Thinking"** checkbox. When checked, you can see the AI's full reasoning process. When unchecked, you only see the key actions — useful for a cleaner demo.

### 3. Decision Flow Diagram (right side)
A visual map showing all the possible paths the AI can take. As the AI makes its decision, you can trace which path it followed.

### 4. Decision Banner
When the AI finishes, a clear, color-coded banner appears showing the final outcome:
- 🟢 **Green** — Problem fixed automatically
- 🔴 **Red** — Human escalation required
- 🟣 **Purple** — Team notified via Slack/Jira

---

## What is a "Model"?

The agents are powered by a **language model** — the same kind of technology behind ChatGPT, but running entirely on your own computer. We use a model called **Gemma 4** made by Google.

**Why does this matter?**
- **Privacy:** Your data never leaves your machine. No cloud. No API keys.
- **Speed:** Because it runs locally, there's no network delay.
- **Cost:** It's completely free to run.

The model is managed by a tool called **Ollama**, which makes it easy to download and run language models locally — think of it as an "app store for AI models."

---

## Key Concepts Explained

### What is "Caveman Protocol"?
When we talk to the AI, we use as few words as possible. Instead of saying *"You are the Auto Remediator agent. You resolve infrastructure issues autonomously without human intervention. Use the RestartKubernetesPod tool..."*, we say *"Role: Auto Remediator. Call RestartKubernetesPod. Return result."*

Fewer words = the AI processes the request faster. We call this the "Caveman Protocol."

### What is the "Second Brain" / "Librarian Pattern"?
Instead of teaching the AI everything upfront (which is slow and expensive), we give it access to a library of documents. When a problem occurs, the AI looks up the relevant document, reads it, and follows the instructions it finds.

This idea comes from [Andrej Karpathy's concept of an AI-driven "Second Brain"](https://techstrong.ai/features/karpathys-instructions-for-building-an-ai-driven-second-brain/).

### What is a "Handoff"?
When one agent finishes its job and passes the work to another agent, that's called a handoff. For example:
- The Incident Commander hands off to the Auto Remediator
- The Incident Commander hands off to a Human in the Loop

Each handoff is logged in the dashboard so you can see exactly what happened and why.

---

## Frequently Asked Questions

**Q: Can I change what the AI does without coding?**
Yes! Edit the Markdown files in the `knowledge/` folder. The AI reads them every time it processes a new alert.

**Q: Does this need an internet connection?**
No. Everything — the AI model, the agents, the dashboard — runs on your local machine.

**Q: What happens if the AI makes a wrong decision?**
The system is designed with safety rails. High-value transactions and compliance issues always get escalated to a human. The AI can only auto-fix problems that your rulebook explicitly marks as safe to fix automatically.

**Q: Can I add more agents?**
Yes. Adding a new agent involves defining its role, giving it a tool, and adding a routing rule in the orchestrator. See the `internal/agents/agents.go` file.

**Q: How fast is it?**
From alert to resolution, the typical cycle is **5–15 seconds** depending on the complexity of the scenario and the speed of your hardware.
