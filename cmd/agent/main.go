package main

import "calculator_go/internal/grpc/agent"

// Agent - the calculating server, which parses expression ->
// calculates the answer -> returns result to Orchestrator

func main() {
	agent.RunAgentServer()
}