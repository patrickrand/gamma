package agent

import (
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
}

// Serve wraps an http.Server and serves the latest content of
// the given Agent.
func (s *Server) Serve() error {
	http.HandleFunc(s.Entrypoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.BindAddr, s.Port), nil)
}
