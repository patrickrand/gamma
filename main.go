package main

import (
	"github.com/patrickrand/gamma/agent"
	"log"
)

var (
	ConfigFile = "templates/agent.json"
	Agent      *agent.Agent
)

func main() {
	var err error
	log.Printf("Starting Gamma...")
	err = agent.LoadFromFile(ConfigFile)
	if err != nil {
		log.Panicf("Exiting Gamma... %v", err)
	}

	results := make(chan *agent.Result)
	for id := range agent.Checks {
		go func(check agent.Check) {
			check.Run(results)
		}(agent.Checks[id])
	}

	for res := range results {
		agent.Results[res.CheckID] = res
		go func(result *agent.Result) {
			check := agent.Checks[result.CheckID]
			for _, id := range check.HandlerIDs {
				if err := agent.Handlers[id].Handle(result); err != nil {
					log.Printf("handler error: %v", err)
				}
			}
		}(res)
	}
}
