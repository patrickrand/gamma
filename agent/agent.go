package agent

import (
	"encoding/json"
	"github.com/patrickrand/gamma/handler"
	log "github.com/patrickrand/gamma/log"
	"github.com/patrickrand/gamma/monitor"
)

type Agent struct {
	Name     string
	Version  string
	Monitors map[monitor.ID]monitor.Monitor
	Handlers map[monitor.ID][]handler.Handler
}

func New(cfg Config) (*Agent, error) {
	return loadAgentFromConfig(cfg)
}

func loadAgentFromConfig(cfg Config) (*Agent, error) {
	log.DBUG("agent", "loadAgentFromConfig => %s", log.PrintJson(cfg))

	agent := &Agent{
		Name:    cfg.AgentName,
		Version: cfg.AgentVersion,
	}

	log.INFO("agent", "Creating new agent => %s %s", agent.Name, agent.Version)

	agent.Monitors = make(map[monitor.ID]monitor.Monitor, 0)
	for i, m := range cfg.Monitors {
		monit := monitor.New(m.(map[string]interface{})["type"].(string))

		data, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, monit)
		if err != nil {
			return nil, err
		}

		agent.Monitors[monitor.ID(i)] = monit
		log.INFO("agent", "Added new agent monitor (%v) => %s", monitor.ID(i), log.PrintJson(monit))
	}

	cfgHandlers := make(map[string]handler.Handler, 0)
	for i, h := range cfg.Handlers {
		handlr := handler.New(h.(map[string]interface{})["type"].(string))

		data, err := json.Marshal(h)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, handlr)
		if err != nil {
			return nil, err
		}

		cfgHandlers[i] = handlr
		log.INFO("agent", "Added new agent handler (%s) => %s", i, log.PrintJson(handlr))
	}

	agent.Handlers = make(map[monitor.ID][]handler.Handler)
	for i, m := range agent.Monitors {
		hList := make([]handler.Handler, 0)
		for _, h := range m.Handlers() {
			hList = append(hList, cfgHandlers[h])
		}
		agent.Handlers[monitor.ID(i)] = hList
	}

	return agent, nil
}
