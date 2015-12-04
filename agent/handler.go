package agent

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type HandlerType string

var (
	HttpHandler HandlerType = "http"
	FileHandler HandlerType = "file"
	SmtpHandler HandlerType = "smtp"
)

type Handler struct {
	ID          string                 `json:"id"`
	Type        HandlerType            `json:"handler_type"`
	Destination string                 `json:"destination"`
	Parameters  map[string]interface{} `json:"parameters"`
}

func (h *Handler) Handle(data interface{}) error {
	handlerFunc, ok := HandlerFuncIndex[h.Type]
	if !ok {
		return fmt.Errorf("invalid handler type %s: no matching handlerfunc", h.Type)
	}
	return handlerFunc(data, h.Destination, h.Parameters)
}

type HandlerFunc func(data interface{}, dest string, params map[string]interface{}) error

func HttpHandlerFunc(data interface{}, dest string, params map[string]interface{}) error {
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

func FileHandlerFunc(data interface{}, dest string, params map[string]interface{}) error {
	return json.NewEncoder(os.Stdout).Encode(data)
}

var HandlerFuncIndex = map[HandlerType]HandlerFunc{
	HttpHandler:   HttpHandlerFunc,
	StdoutHandler: FileHandlerFunc,
}
