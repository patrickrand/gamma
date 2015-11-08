package main

import (
	"encoding/json"
	"flag"
	"github.com/patrickrand/gamma/agent"
	log "github.com/patrickrand/gamma/log"
	"github.com/patrickrand/gamma/result"
	"os"
)

const (
	MAIN        = "MAIN"
	CONFIG_FILE = "/home/patrrand/go-ws/src/github.com/patrickrand/gamma/templates/agent.json"
)

var Agent *agent.Agent

func main() {
	log.INFO(MAIN, "Starting Gamma...")

	debug := flag.Bool("debug", false, "Debug mode")
	if flag.Parse(); *debug {
		log.SetLevel(-1)
	}

	cfg, err := agent.NewConfigFromFile(CONFIG_FILE)
	if err != nil {
		log.EROR(MAIN, "Exiting Gamma... => %s", err.Error())
		os.Exit(1)
	}

	Agent, err = agent.New(cfg)
	if err != nil {
		log.EROR(MAIN, "Exiting Gamma... => %s", err.Error())
		os.Exit(1)
	}

	handle(exec())
}

func exec() <-chan result.Result {
	out := make(chan result.Result)
	go func() {
		for k, m := range Agent.Monitors {
			js, _ := json.Marshal(m)
			log.DBUG(MAIN, "main.exec => %s", string(js))
			res := m.Exec()
			res.MonitorId = string(k)
			if res.Error != nil {
				res.Status = result.StatusErr
			}

			out <- *res
		}
		close(out)
	}()
	return out
}

func handle(in <-chan result.Result) {
	for r := range in {
		log.DBUG(MAIN, "handle => %+v", r)
		for _, h := range Agent.Handlers[r.MonitorId] {
			h.Handle(r)
			if h.Error != nil {
				log.EROR(MAIN, "handle => %s", h.Error.Error())
			}
		}
	}
}
