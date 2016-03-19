package agent

import (
	"encoding/json"
	"fmt"
	"io"
)

type FormatterFunc func(w io.Writer, v interface{}) error

var defaultFormatterFunc = func(w io.Writer, v interface{}) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("agent.defaultFormatterFunc failed to encode value: %v", err)
	}
	return nil
}

var jsonFormatterFunc = func(w io.Writer, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("agent.jsonFormatterFunc failed to marshal value: %v", err)
	}
	if _, err := fmt.Fprint(w, string(data)); err != nil {
		return fmt.Errorf("agent.jsonFormatterFunc failed to write data: %v", err)
	}

	return nil
}

var prettyJSONFormatterFunc = func(w io.Writer, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return fmt.Errorf("agent.prettyJSONFormatterFunc failed to marshal (w/ indent) value: %v", err)
	}
	if _, err := fmt.Fprint(w, string(data)); err != nil {
		return fmt.Errorf("agent.prettyJSONFormatterFunc failed to write data: %v", err)
	}
	return nil
}
