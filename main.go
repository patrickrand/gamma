package main

import (
	"github.com/patrickrand/gamma/agent"
	"log"
	"os"
	"time"
)

var (
	config = "templates/agent.json"

	server = new(agent.Server)
	checks = make(map[string]agent.Check)
)

func main() {
	log.Print("[main] starting gamma...")

	cfg, err := os.Open(config)
	if err != nil {
		log.Fatalf("[error] main failed to read config file: %v", err)
	}
	log.Printf("[main] read configuration settings from %s", config)

	if err := agent.Load(cfg, checks, server); err != nil {
		log.Fatalf("[error] failed to load agent modules: %v", err)
	}

	if server.IsActive {
		go func() { log.Fatal(server.ServeHTTP()) }()
	}

	results := make(chan *agent.Result)
	for _, c := range checks {
		go func(check agent.Check) {
			for range time.Tick(check.Interval * time.Second) {
				results <- check.Exec()
			}
		}(c)
	}

	for r := range results {
		go func(r *agent.Result) {
			server.Cache[r.ID] = *r
			if err := agent.Write(r); err != nil {
				log.Printf("[error] failed to write result: %v", err)
			}
		}(r)
	}
}
