package agent

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
