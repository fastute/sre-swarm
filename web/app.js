

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
    const div = document.createElement('div');
    div.className = `log-entry ${type}`;
    div.innerHTML = `[${new Date().toLocaleTimeString()}] ${message}`;
    container.appendChild(div);
    container.scrollTop = container.scrollHeight;
}

function appendLog(message, type) {
    appendLogTo(message, type, 'log-container');
}

// Set default date to today and disallow past dates
document.addEventListener('DOMContentLoaded', () => {
    const dateInput = document.getElementById('tx-date');
    if (dateInput) {
        const today = new Date().toISOString().split('T')[0];
        dateInput.value = today;
        dateInput.min = today;
    }
});

let isSwarmActive = false;

async function initiateAutonomousSwarm() {
    if (isSwarmActive) {
        appendLog(`⚠️ <b>Swarm is currently active. Please wait for the current cycle to complete.</b>`, 'warning-log');
        return;
    }

    const amount = document.getElementById('tx-amount').value;
    const date = document.getElementById('tx-date').value;
    const signal = document.getElementById('tx-signal').value;

    isSwarmActive = true;
    const logContainer = document.getElementById('log-container');
    logContainer.innerHTML = '';
    
    appendLog(`Transaction Initiated: <strong>£${amount} GBP</strong> on <strong>${date}</strong>`, 'info');
    appendLog(`Narrative Mode: <strong>${signal}</strong>`, 'info');
    appendLog(`Triggering Autonomous SRE Swarm...`, 'info');

    try {
        appendLog(`Waking up Strategic_Triage Agent for Discovery...`, 'agent-triage');
        
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
            appendLog(`🤖 <b>Autonomous Behavior Detected</b>`, 'agent-healer');
            appendLog(`Target Execution: ${swarmData.target}`, 'tool-exec');
            appendLog(`Status: ${swarmData.status}`, 'info');
        } else if (swarmData.handoff === "Human_in_Loop") {
            appendLog(`👤 <b>Human in the Loop Required!</b>`, 'warning-log');
            appendLog(`Reason: ${swarmData.target}`, 'warning-log');
        } else if (swarmData.handoff === "Notification (Code Fix)") {
            appendLog(`📝 <b>Code Fix Required</b>`, 'info');
            appendLog(`Target: PR Drafted for team review.`, 'tool-exec');
        } else {
            appendLog(`Swarm Result: ${swarmData.status}`, 'info');
        }

    } catch (error) {
        appendLog(`Error communicating with backend: ${error.message}`, 'danger-log');
    } finally {
        isSwarmActive = false;
    }
}
