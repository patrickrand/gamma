// Package agent contains methods and functions for running an Agent
// on a host. The agent is a "singleton class" that represents the checks
// and handlers used to monitor a particular host. The agent is typically
// loaded from the agent.json configuration file.
package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
    "sync"
)

var mu *sync.Mutex

// Load decodes the given data into various gamma modules.
func Load(r io.Reader, checks map[string]Check, server *Server) error {
	var modules = struct {
		Checks  map[string]Check `json:"checks"`
		*Server `json:"server"`
	}{
		Checks: checks,
		Server: server,
	}

	if err := json.NewDecoder(r).Decode(&modules); err != nil {
		return fmt.Errorf("agent.Load failed decoding modules: %v", err)
	}

	for id, c := range checks {
		c.ID = id
		checks[id] = c
	}

	server.Cache = make(map[string]Result, len(checks))
	return nil
}

func Write(r *Result) error {
    mu.Lock()
    defer mu.Unlock()

	if err := json.NewEncoder(os.Stdout).Encode(r); err != nil {
		return fmt.Errorf("agent.Write failed to encode result to stdout: %v", err)
	}
	return nil
}
