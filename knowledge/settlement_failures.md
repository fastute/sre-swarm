# Knowledge Base: System Failures Rulebook
(KB-SRE-001)

The `settlement-worker` and `api-gateway` are high-criticality systems.

## Resolution Matrix

### Scenario 1: Memory Leak / OOM Kill (Autonomous)
- **Symptom:** Telemetry Signal is `OOM_KILL`.
- **Action Required:** Auto-Remediation.
- **Protocol:**
  1. Restart the deployment using `RestartKubernetesPod | settlement-worker`.
  2. Instruction: `HANDOFF: Auto_Remediator | settlement-worker`

### Scenario 2: Ledger Deadlock (Architectural Risk)
- **Symptom:** Telemetry Signal is `DEADLOCK`.
- **Action Required:** Exponential Backoff & Retry.
- **Protocol:**
  1. Rollback the current transaction immediately to prevent lock contention.
  2. Trigger an exponential backoff approach to attempt a retry.
  3. Instruction: `HANDOFF: Resilience_Engineer | Initiate Exponential Backoff`

### Scenario 3: P99 Latency Spike
- **Symptom:** Telemetry Signal is `P99_SPIKE`.
- **Action Required:** Auto-Remediation.
- **Protocol:**
  1. Restart the `api-gateway` to clear the hung connection pools.
  2. Instruction: `HANDOFF: Auto_Remediator | api-gateway`

### Global Override Policies
- **Weekend Policy:** If the `Payment Date` is Saturday or Sunday, defer processing to the next business day (Monday). Phone the human agent for resolution.
  - **Instruction:** `HANDOFF: Human_in_Loop | Weekend Guardrail (Deferred to Monday)`
- **High Value Policy:** If the payment amount is > £5,000 GBP, the human agent MUST be notified.
  - **Instruction:** `HANDOFF: Comms_Lead | High Value Payment > £5,000 GBP`
