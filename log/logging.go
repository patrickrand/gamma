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

type LogLevel int

const (
	DBUG_LVL LogLevel = iota - 1
	INFO_LVL
	EROR_LVL
)

var l = logger{
	timestamp: time.RFC3339,
	out:       os.Stderr,
	level:     INFO_LVL,
}

func init() {
	l.Logger = log.New(l.out, "", 0)
}

func Debugf(tag, format string, val ...interface{}) {
	l.write(DBUG_LVL, tag, format, val...)
}

func Infof(tag, format string, val ...interface{}) {
	l.write(INFO_LVL, tag, format, val...)
}

func Errorf(tag, format string, val ...interface{}) {
	l.write(EROR_LVL, tag, format, val...)
}

func PrintJson(data interface{}) string {
	js, _ := json.Marshal(data)
	return string(js)
}

func (l logger) write(lvl LogLevel, tag, format string, val ...interface{}) {
	l.Lock()
	defer l.Unlock()
	if lvl >= l.level {
		if tag != "" {
			format = fmt.Sprintf("%s  [%s] [%s] %s", time.Now().Format(l.timestamp), lvl.String(), tag, format)
		} else {
			format = fmt.Sprintf("%s [%s] %s", time.Now().Format(l.timestamp), lvl.String(), format)
		}

		l.Logger.Printf(format, val...)
	}
}

func Level() LogLevel {
	l.Lock()
	defer l.Unlock()
	return l.level
}

func SetLevel(lvl LogLevel) {
	l.Lock()
	defer l.Unlock()
	l.level = lvl
}

func (lvl LogLevel) String() string {
	l.Lock()
	defer l.Unlock()
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
