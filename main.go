package main

import (
	"github.com/patrickrand/gamma/agent"
	"github.com/patrickrand/gamma/result"
	"log"
)

const CONFIG_FILE = "/home/patrrand/go-ws/src/github.com/patrickrand/gamma/templates/agent.json"

func main() {
	log.Println("Starting Gamma...")

	cfg, err := agent.NewConfigFromFile(CONFIG_FILE)
	if err != nil {
		log.Printf("Error loading config file => %s", err.Error())
		panic("Exiting Gamma...")
	}
	log.Printf("Config loaded from file: %s", CONFIG_FILE)

	agent, err := agent.New(cfg)
	if err != nil {
		log.Printf("Error creating Agent from Config => %s", err.Error())
		panic("Exiting Gamma...")
	}
	log.Printf("Agent initialized: %+v", *agent)

	for k, m := range agent.Monitors {
		res := m.Exec()
		res.MonitorId = string(k)
		if res.Error != nil {
			res.Status = result.StatusErr
		}
		log.Printf("%+v", *res)

	}
}
