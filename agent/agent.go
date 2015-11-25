package agent

import (
	"encoding/json"
	"fmt"
)

type Agent struct {
	Name     string
	Version  string
	Checks   map[string]*Check
	Handlers map[string][]Context
}

func New(cfg Config) (*Agent, error) {
	log.Infof("Creating new agent from config %s", cfg.FilePath)
	agent := &Agent{}
	err := loadAgentFromConfig(cfg, agent)
	return agent, err
}

func loadAgentFromConfig(cfg Config, agent *Agent) (err error) {
	log.Debugf("agent.loadAgentFromConfig => %s", log.PrintJson(cfg))

	if agent.Name = cfg.AgentName; agent.Name == "" {
		return fmt.Errorf("Agent name is required")
	}
	agent.Version = cfg.AgentVersion

	agent.Checks = make(map[string]*Check, 0)
	for id, c := range cfg.Checks {
		c, ok := c.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unable to parse JSON for check: %s", id)
		}

		data, err := json.Marshal(c)
		if err != nil {
			return err
		}

		ch := new(Check)
		if err = json.Unmarshal(data, ch); err != nil {
			return err
		}
		agent.Checks[id] = ch
		log.Infof("Added check %s: %s", id, log.PrintJson(agent.Checks[id]))
	}

	contexts := make(map[string]*Context, 0)
	for id, h := range cfg.Handlers {
		h, ok := h.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unable to parse JSON for handler: %s", id)
		}

		contexts[id] = NewContext()
		dest, ok := h["destination"].(string)
		if !ok {
			return fmt.Errorf("destination not specified for handler: %s", id)
		}
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
