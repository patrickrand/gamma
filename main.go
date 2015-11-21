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
	log.Infof("[%s] Starting Gamma...", MAIN)

	debug := flag.Bool("debug", false, "Debug mode")
	if flag.Parse(); *debug {
		log.SetLevel(log.DBUG_LVL)
	}

	cfg, err := agent.NewConfigFromFile(CONFIG_FILE)
	if err != nil {
		log.Errorf("[%s] Exiting Gamma... => %s", MAIN, err.Error())
		os.Exit(1)
	}

	Agent, err = agent.New(cfg)
	if err != nil {
		log.Errorf("[%s] Exiting Gamma... => %s", MAIN, err.Error())
		os.Exit(1)
	}

	handle(exec())
}

func exec() <-chan result.Result {
	out := make(chan result.Result)
	go func() {
		for k, m := range Agent.Monitors {
			js, _ := json.Marshal(m)
			log.Debugf("[%s] main.exec => %s", MAIN, string(js))
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
		log.Debugf("[%s] handle => %+v", MAIN, r)
		for _, h := range Agent.Handlers[r.MonitorId] {
			h.Handle(r)
			if h.Error != nil {
				log.Errorf("[%s] handle => %s", MAIN, h.Error.Error())
			}
		}
	}
}
