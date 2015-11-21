package agent

import (
	"encoding/json"
	"fmt"
	"github.com/patrickrand/gamma/handler"
	log "github.com/patrickrand/gamma/log"
	"github.com/patrickrand/gamma/monitor"
)

const AGENT = "AGNT"

type Agent struct {
	Name     string
	Version  string
	Monitors map[string]monitor.Monitor
	Handlers map[string][]*handler.Context
}

func New(cfg Config) (*Agent, error) {
	log.Infof(AGENT, "Creating new agent from config %s", cfg.FilePath)
	agent := &Agent{}
	err := loadAgentFromConfig(cfg, agent)
	return agent, err
}

func loadAgentFromConfig(cfg Config, agent *Agent) (err error) {
	log.Debugf(AGENT, "agent.loadAgentFromConfig => %s", log.PrintJson(cfg))

	if agent.Name = cfg.AgentName; agent.Name == "" {
		return fmt.Errorf("Agent name is required")
	}
	agent.Version = cfg.AgentVersion

	agent.Monitors = make(map[string]monitor.Monitor, 0)
	for id, m := range cfg.Monitors {
		m, ok := m.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Unable to parse JSON for monitor %s", id)
		}

		agent.Monitors[id], err = monitor.New(m["type"].(string))
		if err != nil {
			return err
		}

		data, err := json.Marshal(m)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, agent.Monitors[id])
		if err != nil {
			return err
		}

		log.Infof(AGENT, "Added new monitor %v: %s", id, log.PrintJson(agent.Monitors[id]))
	}

	contexts := make(map[string]*handler.Context, 0)
	for id, h := range cfg.Handlers {
		h, ok := h.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Unable to parse JSON for handler %s", id)
		}
		contexts[id] = handler.NewContext()

		if contexts[id].Destination, ok = h["destination"].(string); !ok {
			return fmt.Errorf("Destination not specified for handler %v", id)
		}

		if p, ok := h["parameters"].(map[string]interface{}); ok {
			contexts[id].Parameters = handler.Parameters(p)
		}

		if contexts[id].HandlerFunc, ok = handler.Handlers[h["type"].(string)]; !ok {
			return fmt.Errorf("Invalid handler type for handler %s", id)
		}

		log.Infof(AGENT, "Added new agent handler context %s: %v", id, contexts[id])
	}

	agent.Handlers = make(map[string][]*handler.Context)
	for id, m := range agent.Monitors {
		ctxList := make([]*handler.Context, 0)
		for _, h := range m.Handlers() {
			ctxList = append(ctxList, contexts[h])
		}
		agent.Handlers[id] = ctxList
	}

	return nil
}
