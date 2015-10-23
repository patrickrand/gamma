package monitor

import (
	"time"
)

const (
	StatusErr = iota - 1
	StatusOK
	StatusWarning
	StatusCritical
)

type Output struct {
	Status  int
	Message string
}

type Result struct {
	MonitorId string
	Status    int
	StartTime time.Time
	EndTime   time.Time
	*Output
}
