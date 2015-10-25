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
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	agent := &Agent{
		Name:    cfg.AgentName,
		Version: cfg.AgentVersion,
	}

	agent.Monitors = make(map[monitor.ID]monitor.Monitor, 0)
	for i, m := range cfg.Monitors {
		//m := m.(map[string]interface{})
		monit := monitor.New(m.(map[string]interface{})["type"].(string))

		log.Infof("test", "hello, %s!", "world")
		data, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, monit)
		if err != nil {
			return nil, err
		}

		agent.Monitors[monitor.ID(i)] = monit
	}

	cfgHandlers := make(map[string]handler.Handler, 0)
	for i, h := range cfg.Handlers {
		h := h.(map[string]interface{})
		handlr := handler.New(h["type"].(string))
		data, err := json.Marshal(h)
		err = json.Unmarshal(data, handlr)
		if err != nil {
			return nil, err
		}
		cfgHandlers[i] = handlr
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
