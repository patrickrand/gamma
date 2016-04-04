package main

import (
	"fmt"
	"io"
	"sync"

	"github.com/patrickrand/gamma"
)

type Cache struct {
	*sync.RWMutex
	results map[string]gamma.Result
}

// NewCache initializes and returns a pointer to a new cache.
func NewCache() *Cache {
	return &Cache{RWMutex: new(sync.RWMutex), results: make(map[string]gamma.Result)}
}

// Load writes a cache's results map to the given writer.
func (cache Cache) Load(w io.Writer, fmtFunc gamma.FormatterFunc) error {
	cache.RLock()
	defer cache.RUnlock()

	if fmtFunc == nil {
		fmtFunc = gamma.DefaultFormatterFunc
	}

	if err := fmtFunc(w, cache.results); err != nil {
		return fmt.Errorf("agent.Cache.Load failed to load results map: %v", err)
	}
	return nil
}

// Lookup returns the result associated with the given result ID, if one exists.
func (cache Cache) Lookup(id string) (gamma.Result, bool) {
	cache.RLock()
	defer cache.RUnlock()
	result, ok := cache.results[id]
	return result, ok
}

// Save sets the cache's map with the given ID and result pair.
func (cache Cache) Save(result gamma.Result) {
	cache.Lock()
	cache.results[result.ID] = result
	cache.Unlock()
}
