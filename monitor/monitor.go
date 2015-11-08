package monitor

import (
	"fmt"
	log "github.com/patrickrand/gamma/log"
	"github.com/patrickrand/gamma/result"
	"time"
)

const MONITOR = "MONI"

type Monitor interface {
	ID() string
	Type() string
	Exec() *result.Result
	Interval() time.Duration
	Handlers() []string
}

func New(monitorType string) (Monitor, error) {
	log.DBUG(MONITOR, "monitor.New => %s", monitorType)
	switch monitorType {
	case CHECK:
		return &Check{}, nil
	default:
		return nil, fmt.Errorf("'%s' is not a valid monitor type", monitorType)
	}
}
