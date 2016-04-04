package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/patrickrand/gamma"
)

var (
	config  = "agent.json"
	modules = struct {
		Checks map[string]gamma.Check `json:"checks"`
		Server *Server                `json:"server"`
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

	if err := json.NewDecoder(cfg).Decode(&modules); err != nil {
		log.Fatalf("[error] [main] failed to decode agent modules: %v", err)
	}

	// run checks and and write results
	cache, stdout := make(chan gamma.Result, len(modules.Checks)), make(chan gamma.Result, len(modules.Checks))
	runChecks(cache, stdout)
	writeResults(stdout)

	// save results to server cache
	modules.Server.Cache = NewCache()
	cacheResults(modules.Server.Cache, cache)

	// serve results over HTTP
	addr := fmt.Sprintf("%s:%d", modules.Server.BindAddress, modules.Server.Port)
	log.Printf("[main] serving gamma at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, modules.Server))
}
