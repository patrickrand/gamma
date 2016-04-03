package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

func ExecuteCheck(executer Executer, check Check) Result {
	start := time.Now()
	code, message := executer.Execute(check.Command, check.Args...)
	return Result{
		ID:        check.ID,
		Command:   check.Command,
		Args:      check.Args,
		Interval:  check.Interval,
		StartTime: start,
		EndTime:   time.Now(),
		Code:      code,
		Message:   message,
	}
}

// Decode decodes a reader into an interface, returning an error if the operation fails.
func Decode(r io.Reader, modules interface{}) error {
	if err := json.NewDecoder(r).Decode(modules); err != nil {
		return fmt.Errorf("agent.Decode failed decoding modules: %v", err)
	}
	return nil
}

// Encode encodes a writer into an interface, returning an error if the operation fails.
func Encode(w io.Writer, results interface{}) error {
	if err := json.NewEncoder(w).Encode(results); err != nil {
		return fmt.Errorf("agent.Encode failed to encode result: %v", err)
	}
	return nil
}
