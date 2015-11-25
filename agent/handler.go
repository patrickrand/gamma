package agent

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type HandlerFunc func(data interface{}, dest string, params Parameters) error

func HttpHandler(data interface{}, dest string, params Parameters) error {
	if dest == "" {
		dest = "127.0.0.1"
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", dest, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Infof("%s %s", resp.Status, string(body))
	return nil
}

var StdoutHandler = func(data interface{}, dest string, params Parameters) error {
	return json.NewEncoder(os.Stdout).Encode(data)
}

var Handlers = map[string]HandlerFunc{
	"http":   HttpHandler,
	"stdout": StdoutHandler,
}
