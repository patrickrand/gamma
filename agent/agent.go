package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/patrickrand/gamma"
)

func runChecks(cache, stdout chan gamma.Result) {
	for id, check := range modules.Checks {
		check.ID = id
		modules.Checks[id] = check
		go func(check gamma.Check) {
			for range time.Tick(check.Interval * time.Second) {
				result := check.Run(gamma.NewShellExecuter(-1))
				cache <- result
				stdout <- result
			}
		}(check)
	}
}

func writeResults(stdout chan gamma.Result) {
	go func(stdout chan gamma.Result) {
		for result := range stdout {
			if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
				log.Printf("[error] [main] failed to write result: %v", err)
			}
		}
	}(stdout)
}

func cacheResults(cache *Cache, results chan gamma.Result) {
	go func(cache *Cache, results chan gamma.Result) {
		for result := range results {
			cache.Save(result)
		}
	}(cache, results)
}
