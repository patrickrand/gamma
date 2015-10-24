package result

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
	MonitorId string    `json:"monitor_id"`
	Status    int       `json:"status"`
	Message   string    `json:"message,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Error     error     `json:"error,omitempty"`
}

func New(startTime time.Time) *Result {
	return &Result{StartTime: startTime}
}
