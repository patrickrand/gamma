package agent

import (
	"fmt"
)

type Agent struct {
	Name     string
	Version  string
	Checks   map[string]*Check
	Handlers map[string]*Handler
	Handlers map[string][]Context
}

func New(cfg *Config) (*Agent, error) {
	agent := &Agent{}
	err := loadAgentFromConfig(cfg, agent)
	return agent, err
}

func loadAgentFromConfig(cfg *Config, agent *Agent) (err error) {
	if err := cfg.Get("agent_name", &agent.Name); err != nil {
		return err
	}

	if err := cfg.Get("checks", &agent.Checks); err != nil {
		return err
	}

	log.Infof("loaded checks from config: %#s", log.PrintJson(agent.Checks))

	if err := cfg.Get("handlers", &agent.Handlers); err != nil {
		return err
	}

	log.Infof("loaded handlers from config: %#s", log.PrintJson(agent.Handlers))

	contexts := make(map[string]*Context, 0)
	for id, h := range agent.Handlers {
		contexts[id] = NewContext()
		contexts[id].Destination = dest

		if p, ok := h["parameters"].(map[string]interface{}); ok {
			contexts[id].Parameters = Parameters(p)
		}

		if contexts[id].HandlerFunc, ok = Handlers[h["type"].(string)]; !ok {
			return fmt.Errorf("invalid handler type: %s", id)
		}

		log.Infof("Added agent handler context %s: %s", id, contexts[id])
	}

	agent.Handlers = make(map[string][]Context)
	for id, c := range agent.Checks {
		ctxList := make([]Context, 0)
		for _, h := range c.Handlers() {
			ctxList = append(ctxList, *contexts[h])
		}
		agent.Handlers[id] = ctxList
	}

	return nil
}
