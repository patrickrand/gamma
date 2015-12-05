// Package agent contains methods and functions for running an Agent
// on a host.
package agent

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// An Agent represents the checks and handlers used to monitor a
// particular host. An Agent is typically loaded from the agent.json
// configuration file.
type Agent struct {
	// AgentID is used to uniquely identify an Agent across a network
	// of separate hosts, each with their own Agent.
	AgentID string `json:"agent_id"`

	// HostID is used to uniquely identify the host this Agent is
	// running on within a network. For example, this could be an IP address,
	// FQDN, or CNAME.
	HostID string `json:"host_id"`

	// Checks represents the list of Checks to be executed by the Agent
	// on its host.
	Checks []Check `json:"checks"`

	// Handlers is the map of configured Handlers available to the Agent
	// for pushing a Check's results.
	Handlers map[string]Handler `json:"handlers"`
}

// LoadFromFile reads an agent.json file and decodes it into an Agent.
func LoadFromFile(file string) (*Agent, error) {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(absPath)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	a := new(Agent)
	if err = json.NewDecoder(f).Decode(a); err != nil {
		return nil, err
	}

	log.Printf("Loaded new Agent from file: %s", absPath)
	return a, nil
}
