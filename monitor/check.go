package monitor

import (
	"time"
)

type Check struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Interval int    `json:"interval"`
	Action   `json:"action"`
	Handler  `json:"handler"`
}

func NewCheck() *Check {
	return &Check{}
}

func (c *Check) Id() string {
	return c.ID
}

func (c *Check) Exec() (Result, error) {
	return nil, nil
}

func (c *Check) RuntimeInterval() time.Duration {
	return 0
}
