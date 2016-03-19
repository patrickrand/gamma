package agent

import (
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

// Lookup returns the result associated when the the give result ID, if one exists.
func (cache Cache) Lookup(id string) (Result, bool) {
	cache.RLock()
	defer cache.RUnlock()
	result, ok := cache.results[id]
	return result, ok
}

// Save sets the cache's map with the given ID and result pair.
func (cache Cache) Save(result Result) {
	cache.Lock()
	cache.results[result.ID] = result
	cache.Unlock()
}
