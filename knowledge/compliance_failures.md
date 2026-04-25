# Knowledge Base: Compliance Rulebook
(KB-SRE-002)

## Resolution Matrix

### Scenario 1: AML Engine Sanctions Match (Fraud Alert)
- **Symptom:** Telemetry Signal is `AML_FLAG`.
- **Action Required:** Compliance Escalation.
- **Protocol:**
  1. The transaction must be immediately quarantined.
  2. If the amount is > £10,000 GBP, it MUST be reviewed by a human in the loop.
     - **Instruction:** `HANDOFF: Human_in_Loop | High Value Compliance Review (> £10k)`
  3. Otherwise (<= £10,000 GBP), alert the FinOps team to investigate and file a Jira ticket.
     - **Instruction:** `HANDOFF: Comms_Lead | Fraud Alert Sent to FinOps team & Jira Ticket Created`
