package handler

import (
	"bytes"
	"encoding/json"
	log "github.com/patrickrand/gamma/log"
	"io/ioutil"
	"net/http"
	"os"
)

const HANDLER = "HNDL"

type HandlerFunc func(data interface{}, dest string, params Parameters) error

var HttpHandler = func(data interface{}, dest string, params Parameters) error {
	if dest == "" {
		dest = "127.0.0.1"
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", dest, bytes.NewBuffer(b))
	if err != nil {
		return nil
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

	log.Infof(HANDLER, "%s %s", resp.Status, string(body))
	return nil
}

var StdoutHandler = func(data interface{}, dest string, params Parameters) error {
	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(data)

}

var Handlers = map[string]HandlerFunc{
	"http":   HttpHandler,
	"stdout": StdoutHandler,
}
