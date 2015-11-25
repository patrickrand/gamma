package agent

import (
	"time"
)

const (
	StatusErr = iota - 1
	StatusOK
	StatusWarning
	StatusCritical
)

type Result struct {
	CheckId   string    `json:"check_id"`
	Status    int       `json:"status"`
	Message   string    `json:"message,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Error     error     `json:"error,omitempty"`
}

func NewResult(startTime time.Time) *Result {
	log.Debugf("result.New => %s", startTime.String())
	return &Result{StartTime: startTime}
}
