package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type Cache struct {
	*sync.RWMutex
	results map[string]Result
}

// NewCache initializes and returns a pointer to a new cache.
func NewCache() *Cache {
	return &Cache{RWMutex: new(sync.RWMutex), results: make(map[string]Result)}
}

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

// Load writes a cache's results map to the given writer.
func (cache Cache) Load(w io.Writer, fmtFunc FormatterFunc) error {
	cache.RLock()
	defer cache.RUnlock()

	if fmtFunc == nil {
		fmtFunc = defaultFormatterFunc
	}

	if err := fmtFunc(w, cache.results); err != nil {
		return fmt.Errorf("agent.Cache.Load failed to load results map: %v", err)
	}
	return nil
}

// Save sets the cache's map with the given ID and result pair.
func (cache Cache) Save(result Result) {
	cache.Lock()
	cache.results[result.ID] = result
	cache.Unlock()
}
