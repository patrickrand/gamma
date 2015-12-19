// Package agent contains methods and functions for running an Agent
// on a host.
package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// An Agent represents the checks and handlers used to monitor a
// particular host. An Agent is typically loaded from the agent.json
// configuration file.
type Agent struct {
	// ID is used to uniquely identify an Agent across a network
	// of separate hosts, each with their own Agent.
	ID string `json:"agent_id"`

	// HostID is used to uniquely identify the host this Agent is
	// running on within a network. For example, this could be an IP address,
	// FQDN, or CNAME.
	HostID string `json:"host_id"`

	// Server is the HTTP server that is (optionally) run on the Agent's host.
	// It will serve the latest results of each Check.
	Server `json:"server"`

	// Checks represents the set of Checks to be executed by the Agent
	// on its host.
	Checks map[string]Check `json:"checks"`

	// Handlers is the map of configured Handlers available to the Agent
	// for pushing Check results.
	Handlers map[string]Handler `json:"handlers"`

	// Results is the in-memory cache of the most recent results of each Check.
	// The JSON representation of Results is the response body of requests to server.
	Results map[string]*Result `json:"-"`
}

// LoadFromFile reads an agent.json file and decodes it into an Agent.
func LoadFromFile(filename string) (*Agent, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	a := new(Agent)
	if err = json.NewDecoder(f).Decode(a); err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(*a, "", "    ")
	fmt.Println(string(data))
	log.Printf("loaded new agent \"%s\" from file: %s", a.ID, absPath)
	return a, nil
}

// Initialize initializes certain dynamic aspects of the Agent.
// In particular, the pre-loading of the Results cache with the IDs
// of all the checks, and running the HTTP server if specified.
func (a *Agent) Initialize() {
	a.Results = make(map[string]*Result, len(a.Checks))
	for id, c := range a.Checks {
		c.ID = id
		a.Checks[id] = c
		a.Results[id] = NewResult(&c)
	}

	if a.Server.IsActive {
		go func() { log.Fatal(a.ServeHTTP()) }()
	}
	log.Printf("initialized agent \"%s\"", a.ID)
}

// ServeHTTP wraps an http.Server and serves the latest content of
// the given Agent.
func (a *Agent) ServeHTTP() error {
	http.HandleFunc(a.Server.Entrypoint, func(w http.ResponseWriter, r *http.Request) {
		data, err := json.MarshalIndent(a.Results, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", string(data))
	})
	return http.ListenAndServe(fmt.Sprintf("%s:%d", a.Server.BindAddr, a.Server.Port), nil)
}
