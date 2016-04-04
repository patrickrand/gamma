package gamma

import (
	"encoding/json"
	"fmt"
	"io"
)

type FormatterFunc func(w io.Writer, v interface{}) error

var DefaultFormatterFunc = func(w io.Writer, v interface{}) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("gamma.DefaultFormatterFunc failed to encode value: %v", err)
	}
	return nil
}

var JSONFormatterFunc = func(w io.Writer, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("gamma.JSONFormatterFunc failed to marshal value: %v", err)
	}
	if _, err := fmt.Fprint(w, string(data)); err != nil {
		return fmt.Errorf("gamma.JSONFormatterFunc failed to write data: %v", err)
	}

	return nil
}

var PrettyJSONFormatterFunc = func(w io.Writer, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return fmt.Errorf("gamma.PrettyJSONFormatterFunc failed to marshal (w/ indent) value: %v", err)
	}
	if _, err := fmt.Fprint(w, string(data)); err != nil {
		return fmt.Errorf("gamma.PrettyJSONFormatterFunc failed to write data: %v", err)
	}
	return nil
}
