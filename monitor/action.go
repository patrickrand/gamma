package monitor

import (
	"time"
)

type Action struct {
	Command  string        `json:"command"`
	Resource string        `json:"resource"`
	Interval time.Duration `json:"interval"`
}

func NewAction(command, resource string, interval time.Duration) *Action {
	return &Action{
		Command:  command,
		Resource: resource,
		Interval: interval,
	}
}

func (a *Action) Run() (Result, error) {
	return nil, nil
}
