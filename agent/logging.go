package agent

import (
	"encoding/json"
	"fmt"
	stdlog "log"
	"os"
	"sync"
	"time"
)

type logger struct {
	timestamp string
	out       *os.File
	level     LogLevel
	*stdlog.Logger
	sync.Mutex
}

var log = logger{
	timestamp: time.RFC3339,
	out:       os.Stderr,
	level:     INFO_LVL,
}

func init() {
	log.Logger = stdlog.New(log.out, "", 0)
}

func (l *logger) Debugf(format string, val ...interface{}) {
	l.write(DBUG_LVL, format, val...)
}

func (l *logger) Infof(format string, val ...interface{}) {
	l.write(INFO_LVL, format, val...)
}

func (l *logger) Errorf(format string, val ...interface{}) {
	l.write(EROR_LVL, format, val...)
}

func (l *logger) PrintJson(data interface{}) string {
	js, _ := json.Marshal(data)
	return string(js)
}

func (l *logger) write(lvl LogLevel, format string, val ...interface{}) {
	l.Lock()
	defer l.Unlock()
	if lvl >= l.level {
		format = fmt.Sprintf("%s [%s] %s", time.Now().Format(l.timestamp), lvl.String(), format)
		l.Logger.Printf(format, val...)
	}
}

func (l *logger) Level() LogLevel {
	l.Lock()
	defer l.Unlock()
	return l.level
}

func (l *logger) SetLevel(lvl LogLevel) {
	l.Lock()
	defer l.Unlock()
	l.level = lvl
}

type LogLevel int

const (
	DBUG_LVL LogLevel = iota - 1
	INFO_LVL
	EROR_LVL
)

func (lvl LogLevel) String() string {
	switch lvl {
	case DBUG_LVL:
		return "DBUG"
	case INFO_LVL:
		return "INFO"
	case EROR_LVL:
		return "EROR"
	}
	return ""
}
