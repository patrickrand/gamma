package handler

import (
	log "github.com/patrickrand/gamma/log"
)

var HTTP = "http"

type HttpHandler struct {
	Id          string `json:"id"`
	Destination `json:"destination"`
	Parameters  `json:"parameters"`
}

func NewHttpHandler(id string, dest Destination, params map[string]interface{}) *HttpHandler {
	log.DBUG("http", "handler.NewHttpHandler => %s, %+v, %+v", id, dest, params)
	return &HttpHandler{
		Id:          id,
		Destination: dest,
		Parameters:  params,
	}
}

func (h HttpHandler) Handle(data interface{}, ch chan<- bool) {
	log.DBUG("http", "(HttpHandler).Handle => (%+v).%+v, %+v", h, data, ch)
}
