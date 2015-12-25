package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A Server is the local HTTP server running on the Agent's host.
// It serves the latest, real-time results of the Agent's checks,
// regardless of their Status/AlertOn values.
type Server struct {
	// IsActive indicates whether or not this Server should be run.
	IsActive bool `json:"active"`

	//Entrypoint is the subpath to that content will be served on.
	Entrypoint string `json:"entrypoint"`

	// BindAddr is the local IP address that this Server should
	// bind incoming HTTP requests to.
	BindAddr string `json:"bind_addr"`

	// Port is the local port that this Server is listening on.
	Port int `json:"port"`

	Cache map[string]Result `json:"-"`
}

// ServeHTTP wraps an http.HttpServer and serves the latest content of
// the given Agent.
func (server *Server) ServeHTTP() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.MarshalIndent(server.Cache, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", string(data))
	})

	entrypoint := fmt.Sprintf("%s:%d", server.BindAddr, server.Port)
	fmt.Printf("serving results API at %s\n", entrypoint)
	return http.ListenAndServe(entrypoint, nil)
}
