

// Setup SSE for Live Log Streaming
const eventSource = new EventSource('/api/stream-logs');
eventSource.onmessage = (event) => {
    let logType = 'info';
    const msg = event.data;
    
    // Route to different consoles
    if (msg.startsWith('[PAYMENT]')) {
        const cleanMsg = msg.replace('[PAYMENT] ', '');
        appendLogTo(cleanMsg, 'info', 'payment-logs');
        return;
    }

    if (msg.includes('[Reasoning]')) logType = 'agent-triage';
    if (msg.includes('[Executing Tool]')) logType = 'tool-exec';
    if (msg.includes('[Tool Result]')) logType = 'tool-exec';
    if (msg.includes('[Handoff]')) logType = 'agent-healer';
    
    appendLogTo(msg, logType, 'log-container');
};

function appendLogTo(message, type, containerId) {
    const container = document.getElementById(containerId);
    if (!container) return;
    const div = document.createElement('div');
    div.className = `log-entry ${type}`;

    // If this is a reasoning log and Show Thinking is unchecked, hide it
    if (type === 'agent-triage') {
        const showThinking = document.getElementById('show-thinking');
        if (showThinking && !showThinking.checked) {
            div.classList.add('hidden-reasoning');
        }
    }

    div.innerHTML = `[${new Date().toLocaleTimeString()}] ${message}`;
    container.appendChild(div);
    container.scrollTop = container.scrollHeight;
}

function appendLog(message, type) {
    appendLogTo(message, type, 'log-container');
}

// Show Thinking toggle — retroactively show/hide reasoning logs
document.addEventListener('DOMContentLoaded', () => {
    const showThinking = document.getElementById('show-thinking');
    if (showThinking) {
        showThinking.addEventListener('change', () => {
            const reasoningLogs = document.querySelectorAll('#log-container .log-entry.agent-triage');
            reasoningLogs.forEach(el => {
                if (showThinking.checked) {
                    el.classList.remove('hidden-reasoning');
                } else {
                    el.classList.add('hidden-reasoning');
                }
            });
        });
    }

    updateScenario(); // Initialize the hidden inputs
});

function getNextWeekday() {
    const d = new Date();
    d.setDate(d.getDate() + ((3 + 7 - d.getDay()) % 7 || 7)); // Next Wednesday
    return d.toISOString().split('T')[0];
}

function getNextWeekend() {
    const d = new Date();
    d.setDate(d.getDate() + ((6 + 7 - d.getDay()) % 7 || 7)); // Next Saturday
    return d.toISOString().split('T')[0];
}

function updateScenario() {
    const scenario = document.getElementById('demo-scenario').value;
    const amountInput = document.getElementById('tx-amount');
    const dateInput = document.getElementById('tx-date');
    const signalInput = document.getElementById('tx-signal');
    const simulateBtn = document.getElementById('simulate-btn');
    const cueBox = document.getElementById('presentation-cue');

    if (!scenario) {
        if (simulateBtn) simulateBtn.disabled = true;
        if (cueBox) cueBox.style.display = 'none';
        return;
    }

    if (simulateBtn) simulateBtn.disabled = false;
    if (cueBox) cueBox.style.display = 'block';

    if (scenario === "1") {
        amountInput.value = "250";
        dateInput.value = getNextWeekday();
        signalInput.value = "OOM_KILL";
        if (cueBox) cueBox.innerHTML = "💡 Watch the agent autonomously diagnose and fix a routine OOM Kill without human intervention.";
    } else if (scenario === "2") {
        amountInput.value = "12500";
        dateInput.value = getNextWeekday();
        signalInput.value = "OOM_KILL"; 
        if (cueBox) cueBox.innerHTML = "💡 Now a £12,500 payment fails. The agent reads the High-Value policy and overrides auto-remediation to escalate!";
    } else if (scenario === "3") {
        amountInput.value = "250";
        dateInput.value = getNextWeekend();
        signalInput.value = "OOM_KILL"; 
        if (cueBox) cueBox.innerHTML = "💡 It's the weekend. The agent detects the Weekend Guardrail policy and defers processing until Monday.";
    }
}

let isSwarmActive = false;

function showDecisionBanner(type, icon, title, detail) {
    const banner = document.getElementById('decision-banner');
    banner.className = `decision-banner ${type}`;
    banner.innerHTML = `
        <div class="decision-icon">${icon}</div>
        <div class="decision-title">${title}</div>
        <div class="decision-detail">${detail}</div>
    `;
}

async function initiateAutonomousSwarm() {
    if (isSwarmActive) {
        appendLog(`⚠️ <b>Swarm is currently active. Please wait for the current cycle to complete.</b>`, 'warning-log');
        return;
    }

    const amount = document.getElementById('tx-amount').value;
    const date = document.getElementById('tx-date').value;
    const signal = document.getElementById('tx-signal').value;

    // Reset UI
    isSwarmActive = true;
    const logContainer = document.getElementById('log-container');
    logContainer.innerHTML = '';
    const banner = document.getElementById('decision-banner');
    banner.className = 'decision-banner';
    banner.innerHTML = '';
    
    appendLog(`Transaction Initiated: <strong>£${amount} GBP</strong> on <strong>${date}</strong>`, 'info');
    appendLog(`Signal: <strong>${signal}</strong>`, 'info');
    appendLog(`Triggering Autonomous SRE Swarm...`, 'info');

    try {
        appendLog(`Waking up Incident_Commander for Discovery...`, 'agent-triage');
        
        const swarmRes = await fetch('/api/simulate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                amount: parseFloat(amount), 
                date: date,
                signal: signal 
            })
        });
        
        const swarmData = await swarmRes.json();
        
        if (swarmData.isAutonomous) {
            showDecisionBanner('autonomous', '✅', `Resolved Autonomously → ${swarmData.handoff}`, `Target: ${swarmData.target}`);
            appendLog(`✅ <b>Autonomous Resolution Complete</b>: ${swarmData.target}`, 'agent-healer');
        } else if (swarmData.handoff === "Human_in_Loop") {
            showDecisionBanner('human', '🚨', 'Human in the Loop Required', `Reason: ${swarmData.target}`);
            appendLog(`🚨 <b>Human Escalation</b>: ${swarmData.target}`, 'warning-log');
        } else if (swarmData.handoff === "Comms_Lead") {
            showDecisionBanner('notification', '📢', `Comms Lead Alerted`, `${swarmData.target}`);
            appendLog(`📢 <b>Comms Lead Notified</b>: ${swarmData.target}`, 'info');
        } else if (swarmData.handoff === "Resilience_Engineer") {
            showDecisionBanner('autonomous', '🔄', `Resilience Engineer → Retry Succeeded`, `${swarmData.target}`);
            appendLog(`🔄 <b>Resilience Engineer</b>: ${swarmData.target}`, 'agent-healer');
        } else {
            showDecisionBanner('notification', 'ℹ️', swarmData.status, `Handoff: ${swarmData.handoff || 'None'}`);
            appendLog(`Swarm Result: ${swarmData.status}`, 'info');
        }

    } catch (error) {
        appendLog(`Error communicating with backend: ${error.message}`, 'danger-log');
    } finally {
        isSwarmActive = false;
    }
}
