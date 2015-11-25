package agent

import (
	"fmt"
	"time"
)

type Check struct {
	id         string
	Action     `json:"action"`
	HandlerIds []string `json:"handler_ids"`
}

func NewCheck(id string) *Check {
	return &Check{id: id}
}

func (c *Check) ID() string {
	return c.id
}

func (c *Check) Exec() *Result {
	log.Debugf("(*Check).Exec => %s", log.PrintJson(c))

	res := NewResult(time.Now())
	if output, err := c.Action.Run(); err != nil {
		res.Error = err
	} else if output.Status == nil {
		res.Error = fmt.Errorf("output status is nil")
	} else {
		switch *output.Status {
		case StatusOK, StatusWarning, StatusCritical:
			res.Status = *output.Status
			res.Message = output.Message
		default:
			res.Error = fmt.Errorf("invalid output status: %d", output.Status)
		}
	}
	res.EndTime = time.Now()

	return res
}

func (c *Check) Interval() time.Duration {
	return c.Action.Interval
}

func (c *Check) Handlers() []string {
	return c.HandlerIds
}
