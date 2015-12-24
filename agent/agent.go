// Package agent contains methods and functions for running an Agent
// on a host.
package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
}

type AgentHost struct {

	// ID is used to uniquely identify the host this Agent is
	// running on within a network. For example, this could be an IP address,
	// FQDN, or CNAME.
	ID string `json:"id"`
}

var Host = new(AgentHost)

// LoadFromFile reads an agent.json file and decodes it into various gamma modules.
func LoadFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to load file: %v", err)
	}

	var modules = struct {
		*AgentHost `json:"host"`
		Checks     map[string]Check   `json:"checks"`
		Handlers   map[string]Handler `json:"handlers"`
		*Server    `json:"server"`
	}{
		AgentHost: Host,
		Checks:    Checks,
		Handlers:  Handlers,
		Server:    HttpServer,
	}

	if err := json.Unmarshal(data, &modules); err != nil {
		return err
	}

	Results = make(map[string]*Result, len(Checks))
	for id, c := range Checks {
		c.ID = id
		Checks[id] = c
		Results[id] = NewResult(&c)
	}

	if HttpServer.IsActive {
		go func() { log.Fatal(ServeHTTP()) }()
	}

	log.Printf("loaded new agent from file: %s", filename)
	return nil
}

// ServeHTTP wraps an http.HttpServer and serves the latest content of
// the given Agent.
func ServeHTTP() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.MarshalIndent(Results, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", string(data))
	})

	fmt.Printf("serving results API at %s:%d\n", HttpServer.BindAddr, HttpServer.Port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", HttpServer.BindAddr, HttpServer.Port), nil)
}
