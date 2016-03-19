package agent

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// A Server is the local HTTP server running on the agent's host.
// It serves real-time results of the Agent's checks, regardless of their status.
type Server struct {
	// FQDN is the fully qualified domain name of the server's host.
	FQDN string `json:"fqdn"`

	// IP is the IP address of the server's host.
	IP string `json:"ip"`

	// Enabled indicates whether or not this server should be run on the host.
	Enabled bool `json:"enabled"`

	//EntryPoint is the entry subpath to where content will be served on.
	EntryPoint string `json:"entry_point"`

	// BindAddress is the local IP address that this server should bind all incoming HTTP requests to.
	BindAddress string `json:"bind_address"`

	// Port is the local port that this server is listening on.
	Port int `json:"port"`

	// Cache stores the most current results of each of the checks running on the host.
	Cache *Cache
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[http] %s %s %s", r.RemoteAddr, r.Method, r.URL.String())
	paths := strings.Split(strings.TrimSuffix(r.URL.String(), "/"), "/")

	switch len(paths) {
	case 2: // /<entry_point>
		server.resultsHandler(w, r)
	case 3: // /<entry_point>/<check_id>||<status>
		if _, err := strconv.Atoi(paths[2]); err != nil {
			server.checkHandler(w, r)
			return
		}
		server.statusHandler(w, r)
	default:
		httpLogAndRespond(w, r, http.StatusNotFound, errors.New(""))
	}
}

func (server *Server) resultsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := server.Cache.Load(w, jsonFormatterFunc); err != nil {
			httpLogAndRespond(w, r, http.StatusInternalServerError, err)
		}
	case "PUT": // TODO
		httpLogAndRespond(w, r, http.StatusNotImplemented, errors.New(""))
	default:
		httpLogAndRespond(w, r, http.StatusMethodNotAllowed, errors.New(""))
	}
}

func (server *Server) checkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpLogAndRespond(w, r, http.StatusMethodNotAllowed, errors.New(""))
		return
	}

	resultID := strings.Split(r.URL.String(), "/")[2]
	result, ok := server.Cache.Lookup(resultID)
	if !ok {
		httpLogAndRespond(w, r, http.StatusNotFound, errors.New(""))
		return
	}

	httpLogAndRespond(w, r, http.StatusOK, result)
}

func (server *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpLogAndRespond(w, r, http.StatusMethodNotAllowed, errors.New(""))
		return
	}

	var buf bytes.Buffer
	if err := server.Cache.Load(&buf, defaultFormatterFunc); err != nil {
		httpLogAndRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	var results map[string]Result
	if err := Decode(&buf, &results); err != nil {
		httpLogAndRespond(w, r, http.StatusInternalServerError, err)
		return
	}

	status := strings.Split(r.URL.String(), "/")[2]
	requestedResults := make(map[string]Result, 0)
	// NOTE: also implement for CLI
	for id, result := range results {
		if strconv.Itoa(result.StatusCode) == status {
			requestedResults[id] = result
		}
	}

	httpLogAndRespond(w, r, http.StatusOK, requestedResults)
}

func httpLogAndRespond(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	method := r.Method
	url := r.URL.String()

	switch v := v.(type) {
	case error:
		log.Printf("[http] %d %s %s %v", code, method, url, v)
		if code < 400 {
			code = http.StatusInternalServerError
		}
		http.Error(w, http.StatusText(code), code)
	default:
		log.Printf("[http] %d %s %s", code, method, url)
		w.WriteHeader(code)
		if pretty := r.URL.Query().Get("pretty"); pretty == "" || pretty == "true" {
			prettyJSONFormatterFunc(w, v)
			return
		}
		jsonFormatterFunc(w, v)
	}
}
