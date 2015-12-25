package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// A Handler wraps the data and logic needed to push Check results from
// the Agent.
type Handler struct {
	// Type is the HandlerType of this Handler.
	Type HandlerType `json:"type"`

	// Destination is the destination where the Handler pushes results to.
	// The interpretation of this string value depends on the HandlerFunc
	// used by the Handler (which is determined by its HandlerType). Destination
	// could be a file name, HTTP URL, database connection string, etc.
	Destination string `json:"destination"`

	// Parameters is a semi-arbitrary set of key-value pairs that will be
	// passed to this Handler's HandlerFunc. Similiar to Destination, the
	// interpretation of these pairs are determined by the HandlerType/HandlerFunc
	// of the Handler. Examples of pairs are HTTP Header parameters, file permissions,
	// and authentication credentials.
	Parameters map[string]interface{} `json:"parameters"`
}

// Handlers is the map of configured Handlers available to the Agent
// for pushing Check results.
var handlers = make(map[string]Handler)

// Handle is a factory-style method that handles an Agent's Check results according
// to its HandlerType.
func (h Handler) Handle(data interface{}) error {
	handlerFunc, ok := HandlerFuncIndex[h.Type]
	if !ok {
		return fmt.Errorf("invalid handler type %s: no matching handlerfunc", h.Type)
	}
	return handlerFunc(data, h.Destination, h.Parameters)
}

// A HandlerType encapsulates the identifiers of the various Handler types
// available to an Agent.
type HandlerType string

var (
	// HttpClient maps the Handler's Handle method to the HttpClientFunc.
	HttpClient HandlerType = "http_client"

	// FileWriter maps the Handler's Handle method to the FileWriterFunc.
	FileWriter HandlerType = "file_writer"

	// SmtpClient maps the Handler's Handle method to the SmtpClientFunc.
	SmtpClient HandlerType = "smtp_client"

	// HandlerFuncIndex is the lookup table that maps a HandlerType to its
	// corresponding HandlerFunc.
	HandlerFuncIndex = map[HandlerType]HandlerFunc{
		HttpClient: HttpClientFunc,
		FileWriter: FileWriterFunc,
	}
)

// A HandlerFunc contains specialized logic for handling the pushing of an Agent's Check results.
type HandlerFunc func(data interface{}, dest string, params map[string]interface{}) error

// HttpClientFunc is a HandlerFunc that pushes data via an HTTP POST request to an endpoint.
func HttpClientFunc(data interface{}, dest string, params map[string]interface{}) error {
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

	log.Printf("%s %s", resp.Status, string(body))
	return nil
}

// FileWriterFunc is a HandlerFunc that pushes data via writing to a file.
func FileWriterFunc(data interface{}, dest string, params map[string]interface{}) error {
	var f *os.File
	switch dest {
	case "stdin", "/dev/stdin":
		f = os.Stdin
	case "stdout", "/dev/stdout":
		f = os.Stdout
	case "stderr", "/dev/stderr":
		f = os.Stderr
	default:
		var err error
		if f, err = os.Open(dest); err != nil {
			return err
		}
	}
	return json.NewEncoder(f).Encode(data)
}
