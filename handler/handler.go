package handler

import (
	log "github.com/patrickrand/gamma/log"
)

type Handler interface {
	Handle(data interface{}, ch chan<- bool)
}

func New(handlerType string) Handler {
	log.DBUG("handler", "handler.New => %s", handlerType)
	switch handlerType {
	case HTTP:
		return &HttpHandler{}
	default:
		return nil
	}
}
