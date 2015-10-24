package monitor

import (
	"github.com/patrickrand/gamma/result"
	"time"
)

type ID string

type Monitor interface {
	Type() string
	Exec() *result.Result
	Interval() time.Duration
	Handlers() []string
}

func New(monitorType string) Monitor {
	switch monitorType {
	case CHECK:
		return &Check{}
	default:
		return nil
	}
}
