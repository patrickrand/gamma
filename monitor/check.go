package monitor

import (
	"github.com/patrickrand/gamma/handler"
	"time"
)

var (
	CHECK = "check"
)

type Check struct {
	ID       string `json:"id"`
	Action   `json:"action"`
	Handlers []handler.Handler `json:"handlers"`
}

func NewCheck(id string, action Action, handlers []handler.Handler) *Check {
	return &Check{
		ID:       id,
		Action:   action,
		Handlers: handlers,
	}
}

func (c *Check) Id() string {
	return c.ID
}

func (c *Check) Type() string {
	return CHECK
}

func (c *Check) Exec() (Result, error) {
	return nil, nil
}

func (c *Check) Interval() time.Duration {
	return 0
}
