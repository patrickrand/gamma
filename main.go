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
	Agent, err = agent.LoadFromFile(ConfigFile)
	if err != nil {
		log.Panicf("Exiting Gamma... %v", err)
	}
	Agent.Initialize()

	results := make(chan *agent.Result, len(Agent.Checks))
	for id := range Agent.Checks {
		go func(check agent.Check) {
			check.Run(results)
		}(Agent.Checks[id])
	}

	for res := range results {
		Agent.Results[res.CheckID] = res
		go func(result *agent.Result) {
			check := Agent.Checks[result.CheckID]
			for _, id := range check.HandlerIDs {
				if err := Agent.Handlers[id].Handle(result); err != nil {
					log.Printf("handler error: %v", err)

				}
			}
		}(res)
	}
}
