package main

import (
	"github.com/patrickrand/gamma/agent"
	"log"
)

const (
	MAIN        = "MAIN"
	CONFIG_FILE = "templates/agent.json"
)

var (
	Agent *agent.Agent
)

func main() {
	log.Printf("Starting Gamma...")

	var err error
	Agent, err = agent.LoadFromFile(CONFIG_FILE)
	if err != nil {
		log.Panicf("Exiting Gamma... %v", err)
	}
	handle(exec())
}

func exec() <-chan *agent.Result {
	out := make(chan *agent.Result)
	go func() {
		for _, check := range Agent.Checks {
			result := check.Exec()
			if result.Error != "" {
				result.Status = new(int)
				*result.Status = agent.StatusErr
			}
			out <- result
		}
		close(out)
	}()
	return out
}

func handle(in <-chan *agent.Result) {
	for r := range in {
		for _, hid := range r.Check.HandlerIDs {
			if err := Agent.Handlers[hid].Handle(r); err != nil {
				log.Printf("ERROR: handle => %#v", err)
			}
		}
	}
}
