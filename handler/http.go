package handler

// check := NewCheck()
// result := check.Exec()

// handler := handlerFactory[check.ID()]
// handler.Handle(result)

var HTTP = "http"

type HttpHandler struct {
	Id          string `json:"id"`
	Destination `json:"destination"`
	Parameters  `json:"parameters"`
}

func NewHttpHandler(id string, dest Destination, params map[string]interface{}) *HttpHandler {
	return &HttpHandler{
		Id:         id,
		Parameters: params,
	}
}

func (h HttpHandler) Handle(data interface{}, ch chan<- bool) {

}
