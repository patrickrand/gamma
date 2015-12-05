package agent

import (
	"time"
)

type Output struct {
	Status  *int   `json:"status"`
	Message string `json:"message,omitempty"`
}

type Result struct {
	CheckID   string `json:"check_id"`
	Output    `json:"output"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Error     string    `json:"error,omitempty"`
}

const (
	StatusErr = iota - 1
	StatusOK
	StatusWarning
	StatusCritical
)

func NewResult(checkID string, startTime time.Time) *Result {
	return &Result{
		CheckID:   checkID,
		StartTime: startTime,
	}
}
