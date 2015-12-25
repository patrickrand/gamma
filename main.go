package main

import (
	"github.com/patrickrand/gamma/agent"
	"log"
)

var (
	config = "templates/agent.json"

	host   = new(agent.Host)
	server = new(agent.Server)

	checks   = make(map[string]agent.Check)
	handlers = make(map[string]agent.Handler)
)

func main() {
	log.Printf("Starting Gamma...")
	if err := agent.LoadFromFile(config, host, checks, handlers, server); err != nil {
		log.Fatalf("failed to load agent from file %s: %v", err)
	}

	if server.IsActive {
		go func() { log.Fatal(server.ServeHTTP()) }()
	}

	results := make(chan *agent.Result)
	for _, c := range checks {
		go func(check agent.Check) {
			check.Run(results)
		}(c)
	}

	for r := range results {
		go func(r *agent.Result) {
			server.Cache[r.CheckID] = *r
			check := checks[r.CheckID]
			for _, id := range check.HandlerIDs {
				if err := handlers[id].Handle(r); err != nil {
					log.Printf("handler error: %v", err)
				}
			}
		}(r)
	}
}
