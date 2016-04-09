package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/patrickrand/gamma"
)

type Agent struct {
	Checks         []gamma.Check     `json:"checks"`
	Results        chan gamma.Result `json:"-"`
	gamma.Executer `json:"executer"`
	*Server        `json:"server"`
	Errors         chan error `json:"-"`
}

// ToDo: make slice of "handler" channels

func NewAgent(checks []gamma.Check, results chan gamma.Result, executer gamma.Executer, server *Server, errors chan error) *Agent {
	return &Agent{
		Checks:   checks,
		Results:  results,
		Executer: executer,
		Server:   server,
		Errors:   errors,
	}
}

func (agent *Agent) Run() {
	agent.runChecks()
	agent.handleResults()
}

func (agent *Agent) runChecks() {
	for _, check := range agent.Checks {
		go func(agent *Agent, check gamma.Check) {
			for range time.Tick(check.Interval * time.Second) {
				agent.Results <- check.Run(agent)
			}
		}(agent, check)
	}
}

func (agent *Agent) handleResults() {
	go func(agent *Agent) {
		for result := range agent.Results {
			agent.Server.Cache.Save(result)
			if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
				agent.Errors <- fmt.Errorf("agent.handleResults failed to encode result: %v", err)
			}
		}
	}(agent)
}
