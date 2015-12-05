package main

import (
	"encoding/json"
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

func exec() <-chan agent.Result {
	out := make(chan agent.Result)
	go func() {
		for k, m := range Agent.Checks {
			js, _ := json.Marshal(m)
			log.Printf("main.exec => %s", string(js))
			res := m.Exec()
			res.CheckId = string(k)
			if res.Error != nil {
				res.Status = new(int)
				*res.Status = agent.StatusErr
			}

			out <- *res
		}
		close(out)
	}()
	return out
}

func handle(in <-chan agent.Result) {
	for r := range in {
		log.Printf("handle => %+v", r)
		for _, h := range Agent.Handlers {
			if err := h.Handle(r); err != nil {
				log.Printf("handle => %v", err)
			}
		}
	}
}
