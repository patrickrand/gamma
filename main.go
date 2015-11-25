package main

import (
	"encoding/json"
	//	"flag"
	"github.com/patrickrand/gamma/agent"
	"log"
	"os"
)

const (
	MAIN        = "MAIN"
	CONFIG_FILE = "templates/agent.json"
)

var Agent *agent.Agent

func main() {
	log.Printf("[%s] Starting Gamma...", MAIN)

	//debug := flag.Bool("debug", false, "Debug mode")
	//if flag.Parse(); *debug {
	//	log.SetLevel(log.DBUG_LVL)
	//}

	cfg, err := agent.NewConfigFromFile(CONFIG_FILE)
	if err != nil {
		log.Printf("[%s] Exiting Gamma... => %s", MAIN, err.Error())
		os.Exit(1)
	}

	Agent, err = agent.New(cfg)
	if err != nil {
		log.Printf("[%s] Exiting Gamma... => %s", MAIN, err.Error())
		os.Exit(1)
	}

	handle(exec())
}

func exec() <-chan agent.Result {
	out := make(chan agent.Result)
	go func() {
		for k, m := range Agent.Checks {
			js, _ := json.Marshal(m)
			log.Printf("[%s] main.exec => %s", MAIN, string(js))
			res := m.Exec()
			res.CheckId = string(k)
			if res.Error != nil {
				res.Status = agent.StatusErr
			}

			out <- *res
		}
		close(out)
	}()
	return out
}

func handle(in <-chan agent.Result) {
	for r := range in {
		log.Printf("[%s] handle => %+v", MAIN, r)
		for _, h := range Agent.Handlers[r.CheckId] {
			h.Handle(r)
			if h.Error != nil {
				log.Printf("[%s] handle => %s", MAIN, h.Error.Error())
			}
		}
	}
}
