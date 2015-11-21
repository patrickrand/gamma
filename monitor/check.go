package monitor

import (
	"fmt"
	log "github.com/patrickrand/gamma/log"
	"github.com/patrickrand/gamma/result"
	"time"
)

const CHECK = "check"

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

func (c *Check) Type() string {
	return CHECK
}

func (c *Check) Exec() *result.Result {
	log.Debugf(MONITOR, "(*Check).Exec => (%s)", log.PrintJson(c))

	res := result.New(time.Now())

	output, err := c.Action.Run()
	if err != nil {
		res.Error = err
	} else if output.Status == nil {
		res.Error = fmt.Errorf("Output status is nil")
	} else {
		switch *output.Status {
		case result.StatusOK, result.StatusWarning, result.StatusCritical:
			res.Status = *output.Status
			res.Message = output.Message
		default:
			res.Error = fmt.Errorf("Invalid output status: %d", output.Status)
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
