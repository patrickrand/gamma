package main

import (
	"github.com/patrickrand/gamma/agent"
	"io/ioutil"
	"log"
	"time"
)

var (
	config = "templates/agent.json"

	host   = new(agent.Host)
	server = new(agent.Server)

	checks   = make(map[string]agent.Check)
	handlers = make(map[string]agent.Handler)
)

func init() {
	log.Print("initializing gamma")

	data, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatalf("failed to read config file: %v")
	}

	if err := agent.Load(data, host, checks, handlers, server); err != nil {
		log.Fatalf("failed to load agent modules: %v", err)
	}

}

func main() {
	log.Print("running gamma")

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
			for _, id := range r.HandlerIDs {
				if err := handlers[id].Handle(r); err != nil {
					log.Printf("handler error: %v", err)
				}
			}
		}(r)
	}
}
