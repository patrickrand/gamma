package monitor

import (
	log "github.com/patrickrand/gamma/log"
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
	log.DBUG("monitor", "monitor.New => %s", monitorType)
	switch monitorType {
	case CHECK:
		return &Check{}
	default:
		return nil
	}
}
