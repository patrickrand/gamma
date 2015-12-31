// Package agent contains methods and functions for running an Agent
// on a host. The agent is a "singleton class" that represents the checks
// and handlers used to monitor a particular host. The agent is typically
// loaded from the agent.json configuration file.
package agent

import (
	"encoding/json"
	"fmt"
)

// Load decodes the given data into various gamma modules.
func Load(data []byte, host *Host, checks map[string]Check, handlers map[string]Handler, server *Server) error {

	var modules = struct {
		*Host    `json:"host"`
		Checks   map[string]Check   `json:"checks"`
		Handlers map[string]Handler `json:"handlers"`
		*Server  `json:"server"`
	}{
		Host:     host,
		Checks:   checks,
		Handlers: handlers,
		Server:   server,
	}

	if err := json.Unmarshal(data, &modules); err != nil {
		return fmt.Errorf("failed to load modules: %v", err)
	}

	for id, c := range checks {
		c.ID = id
		checks[id] = c
	}

	server.Cache = make(map[string]Result, len(checks))
	return nil
}
