package main

import (
	"github.com/patrickrand/gamma/agent"
	"log"
)

const (
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

func exec() <-chan int {
	out := make(chan int)
	go func() {
		for i := range Agent.Checks {
			result := Agent.Checks[i].Exec()
			if result.Error != "" {
				result.Status = new(int)
				*result.Status = agent.StatusErr
			}
			Agent.Checks[i].Result = result
			out <- i
		}
		close(out)
	}()
	return out
}

func handle(in <-chan int) {
	for i := range in {
		for _, hid := range Agent.Checks[i].HandlerIDs {
			if err := Agent.Handlers[hid].Handle(Agent.Checks[i].Result); err != nil {
				log.Printf("ERROR: handle => %#v", err)
			}
		}
	}
}
