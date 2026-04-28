// Package main is the SRE Swarm orchestrator — receives alerts from the
// dashboard, runs them through the ADK agent swarm, and streams results via SSE.
package main

import (
	"fmt"
	"log"
)

// ──────────────────────────────────────────────
//  SSE Log Stream
// ──────────────────────────────────────────────

var logStream = make(chan string, 100)

func broadcast(msg string) {
	select {
	case logStream <- msg:
	default:
	}
}

func emit(tag, text string) {
	msg := fmt.Sprintf("[%s] %s", tag, text)
	log.Println(msg)
	broadcast(msg)
}
