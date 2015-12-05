package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type Check struct {
	ID         string        `json:"id"`
	Command    string        `json:"command"`
	Interval   time.Duration `json:"interval"`
	HandlerIds []string      `json:"handler_ids"`
}

func (c *Check) Exec() *Result {
	res := NewResult(time.Now())
	data, err := exec.Command(c.Command).Output()
	if err != nil {
		res.Error = err
	} else if err = json.NewDecoder(bytes.NewReader(data)).Decode(res.Output); err != nil {
		res.Error = err
	} else if res.Status == nil {
		res.Error = fmt.Errorf("output status is nil")
	}
	res.EndTime = time.Now()

	return res
}
