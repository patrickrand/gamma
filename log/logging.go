package log

import (
	"fmt"
	"log"
)

func Infof(tag, format string, val ...interface{}) {
	if tag != "" {
		format = fmt.Sprintf("[%s] %s", tag, format)
	}
	log.Printf("[INFO] "+format, val...)
}
