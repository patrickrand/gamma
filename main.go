package main

import (
	//"fmt"
	"github.com/patrickrand/gamma/agent"
	"log"
	//"os"
)

var (
	ConfigFile = "templates/agent.json"
	Agent      *agent.Agent
)

// NOTE: This is a temporary main.go file, and is only being used
// 		 for testing and development purposes.
func main() {
	log.Printf("Starting Gamma...")

	var err error
	Agent, err = agent.LoadFromFile(ConfigFile)
	if err != nil {
		log.Panicf("Exiting Gamma... %v", err)
	}

	if Agent.Server.IsActive {
		go func() { log.Fatal(Agent.Serve()) }()
	}

	checks := make(chan string)
	abort := make(chan struct{}, 2)
	/*
		go func() {
			os.Stdin.Read(make([]byte, 1))
			for id := range checks {
				fmt.Printf("abort: %s\n", id)
				abort <- struct{}{}
			}
		}()
	*/
	for id := range Agent.Checks {
		check := Agent.Checks[id]
		go check.Run(checks, abort)
	}

	for id := range checks {
		check := Agent.Checks[id]
		go func(string) {
			for _, hid := range check.HandlerIDs {
				if err := Agent.Handlers[hid].Handle(check.Result); err != nil {
					log.Printf("ERROR: handle => %#v", err)

				}
			}
		}(id)
	}

}

/*
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
			if Agent.Checks[i].ShouldAlert(result.Status) {
				out <- i
			}
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
*/
