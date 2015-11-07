package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type logger struct {
	timestamp string
	out       io.Writer
	level     int
}

const (
	DBUG_LVL = iota - 1
	INFO_LVL
	WARN_LVL
	EROR_LVL
)

var (
	l = newLogger(time.RFC3339, os.Stderr, INFO_LVL)
)

func newLogger(timestamp string, out io.Writer, level int) *logger {
	log.SetFlags(0)
	return &logger{timestamp: timestamp, out: out, level: level}
}

func DBUG(tag, format string, val ...interface{}) {
	write(DBUG_LVL, tag, format, val...)
}

func INFO(tag, format string, val ...interface{}) {
	write(INFO_LVL, tag, format, val...)
}

func WARN(tag, format string, val ...interface{}) {
	write(WARN_LVL, tag, format, val...)
}

func EROR(tag, format string, val ...interface{}) {
	write(EROR_LVL, tag, format, val...)
}

func PrintJson(data interface{}) string {
	js, _ := json.Marshal(data)
	return string(js)
}

func write(lvl int, tag, format string, val ...interface{}) {
	if lvl >= l.level {
		if tag != "" {
			format = fmt.Sprintf("%s  [%s] [%s] %s", time.Now().Format(l.timestamp), PrintLevel(lvl), tag, format)
		} else {
			format = fmt.Sprintf("%s [%s] %s", time.Now().Format(l.timestamp), PrintLevel(lvl), format)
		}

		log.Printf(format, val...)
	}
}

func Level() int {
	return l.level
}

func SetLevel(lvl int) {
	if lvl >= DBUG_LVL && lvl <= EROR_LVL {
		l.level = lvl
	}
}

func PrintLevel(lvl int) string {
	switch lvl {
	case DBUG_LVL:
		return "DBUG"
	case INFO_LVL:
		return "INFO"
	case WARN_LVL:
		return "WARN"
	case EROR_LVL:
		return "EROR"
	}
	return ""
}
