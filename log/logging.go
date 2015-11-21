package log

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type logger struct {
	timestamp string
	out       *os.File
	level     LogLevel
	*log.Logger
	sync.Mutex
}

var l = logger{
	timestamp: time.RFC3339,
	out:       os.Stderr,
	level:     INFO_LVL,
}

func init() {
	l.Logger = log.New(l.out, "", 0)
}

func Debugf(format string, val ...interface{}) {
	l.write(DBUG_LVL, format, val...)
}

func Infof(format string, val ...interface{}) {
	l.write(INFO_LVL, format, val...)
}

func Errorf(format string, val ...interface{}) {
	l.write(EROR_LVL, format, val...)
}

func PrintJson(data interface{}) string {
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

func Level() LogLevel {
	return l.logLevel()
}

func (l *logger) logLevel() LogLevel {
	l.Lock()
	defer l.Unlock()
	return l.level
}

func SetLevel(lvl LogLevel) {
	l.setLevel(lvl)
}

func (l *logger) setLevel(lvl LogLevel) {
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
