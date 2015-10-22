package monitor

import (
	"time"
)

type Monitor interface {
	Id() string
	Type() string
	Exec() (Result, error)
	Interval() time.Duration
}
