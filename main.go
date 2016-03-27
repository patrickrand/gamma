package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/patrickrand/gamma/agent"
)

var (
	config  = "agent.json"
	modules = struct {
		Checks map[string]agent.Check `json:"checks"`
		Server *agent.Server          `json:"server"`
	}{}
)

func main() {
	log.Print("[main] starting gamma...")

	// load config and decode modules
	cfg, err := os.Open(config)
	if err != nil {
		log.Fatalf("[error] [main] failed to read config file: %v", err)
	}
	log.Printf("[main] read configuration settings from %s", config)

	if err := agent.Decode(cfg, &modules); err != nil {
		log.Fatalf("[error] [main] failed to decode agent modules: %v", err)
	}

	// run checks and output results
	cache, stdout := make(chan agent.Result, len(modules.Checks)), make(chan agent.Result, len(modules.Checks))
	for id, check := range modules.Checks {
		check.ID = id
		modules.Checks[id] = check
		go func(check agent.Check) {
			for range time.Tick(check.Interval * time.Second) {
				result := check.Run(agent.NewShellExecuter(-1))
				cache <- result
				stdout <- result
			}
		}(check)
	}

	// write results to stdout
	go func(stdout chan agent.Result) {
		for result := range stdout {
			if err := agent.Encode(os.Stdout, result); err != nil {
				log.Printf("[error] [main] failed to write result: %v", err)
			}
		}
	}(stdout)

	// save results to server cache
	modules.Server.Cache = agent.NewCache()
	go func(server *agent.Server, cache chan agent.Result) {
		for result := range cache {
			server.Cache.Save(result)
		}
	}(modules.Server, cache)

	// serve results over HTTP
	addr := fmt.Sprintf("%s:%d", modules.Server.BindAddress, modules.Server.Port)
	log.Printf("[main] serving gamma at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, modules.Server))
}
