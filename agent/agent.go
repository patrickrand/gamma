package agent

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Agent struct {
	Name     string
	Version  string
	Checks   []Check
	Handlers []Handler
}

func LoadFromFile(file string) (*Agent, error) {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(absPath)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	a := new(Agent)
	if err = json.NewDecoder(f).Decode(a); err != nil {
		return nil, err
	}

	log.Printf("Loaded new Agent from file: %s", absPath)
	return a, nil
}
