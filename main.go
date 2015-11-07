package main

import (
	"encoding/json"
	"flag"
	"github.com/patrickrand/gamma/agent"
	log "github.com/patrickrand/gamma/log"
	"github.com/patrickrand/gamma/result"
)

const CONFIG_FILE = "/home/patrrand/go-ws/src/github.com/patrickrand/gamma/templates/agent.json"

func main() {
	log.INFO("main", "Starting Gamma...")

	cfg, err := agent.NewConfigFromFile(CONFIG_FILE)
	if err != nil {
		log.EROR("main", "Exiting Gamma... => %s", err.Error())
		panic("")
	}

	agent, err := agent.New(cfg)
	if err != nil {
		log.EROR("main", "Exiting Gamma... => %s", err.Error())
		panic("")
	}

	logLevel := flag.String("log", "info", "Log level")
	flag.Parse()
	switch *logLevel {
	case "debug":
		log.SetLevel(-1)
	}

	handle(exec(agent))
}

func exec(agent *agent.Agent) <-chan result.Result {
	out := make(chan result.Result)
	go func() {
		for k, m := range agent.Monitors {
			js, _ := json.Marshal(m)
			log.INFO("main", "exec => %s", string(js))
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
		log.INFO("main", "handle => %+v", r)
	}
}
