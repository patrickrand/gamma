package agent

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type Agent struct {
	Name    string   `json:"name"`
	Version Version  `json:"version"`
	Checks  []Check  `json:"checks"`
	Results []Result `json:"results"`
}

func NewAgent(cfg Config) (*Agent, error) {
	return loadAgentFromConfig(cfg)
}
