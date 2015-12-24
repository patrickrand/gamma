package agent

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var testData = []byte(`
{
	"agent_id": "test_agent",
	"host_id": "localhost",
	"server": {
		"active": true,
		"entrypoint": "/test_agent",
		"bind_addr": "0.0.0.0",
		"port": 7070
	},
	"checks": {
		"test_check-1": {
			"command": "/usr/bin/bash ./test_check -c 25 -w 50",
			"interval": 30,
			"alert_on": "critical",
			"handler_ids": ["test_handler-1"]
		},
		"test_check-2": {
			"command": "/usr/bin/bash ./test_check -c 10 -w 40",
			"interval": 15,
			"alert_on": "warning",
			"handler_ids": ["test_handler-1", "test_handler-2"]
		}
	},
    "handlers": { 
        "test_handler-1": {
            "type": "http_client",
            "destination": "localhost:8080",                
            "parameters": {}
        },
        "test_handler-2": {
            "type": "file_writer",
            "destination": "/dev/stdout",
            "parameters": {}
        }
    }
}`)

func TestLoadFromFile(t *testing.T) {
	// create expected Agent struct from test data
	var expected Agent
	if err := json.Unmarshal(testData, &expected); err != nil {
		t.Fatal(err.Error())
	}

	// setup test environment
	filename := filepath.Join(os.TempDir(), "agent.json")
	if err := ioutil.WriteFile(filename, testData, 0755); err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(filename)

	// test
	got, err := LoadFromFile(filename)
	if err != nil {
		t.Errorf("expected: %v, got: %q", nil, err.Error())
	}

	// verify success
	if *got != expected {
		t.Errorf("expected: %#v\ngot: %#v", expected, *got)
	}
}
