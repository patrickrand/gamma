package agent

import (
	"fmt"
	"log"
	"net/http"
)

// A Server is the local HTTP server running on the Agent's host.
// It serves the latest, real-time results of the Agent's checks,
// regardless of their Status/AlertOn values.
type Server struct {
	ID   string `json:"id"`
	FQDN string `json:"fqdn"`
	IP   string `json:"ip"`
	SSL  bool   `json:"ssl"`

	// IsActive indicates whether or not this Server should be run.
	IsActive bool `json:"active"`

	//Entrypoint is the subpath to that content will be served on.
	Entrypoint string `json:"entrypoint"`

	// BindAddr is the local IP address that this Server should
	// bind incoming HTTP requests to.
	BindAddr string `json:"bind_addr"`

	// Port is the local port that this Server is listening on.
	Port int `json:"port"`

	*Cache `json:"-"`
}

// ServeHTTP wraps an http.HttpServer and serves the latest content of
// the given Agent.
func (server *Server) ServeHTTP() error {
	http.HandleFunc(server.Entrypoint, func(w http.ResponseWriter, r *http.Request) {
		if err := server.Cache.Load(w, jsonFormatterFunc); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	addr := fmt.Sprintf("%s:%d", server.BindAddr, server.Port)
	log.Printf("serving gamma at %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
