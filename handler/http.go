package handler

// check := NewCheck()
// result := check.Exec()

// handler := handlerFactory[check.ID()]
// handler.Handle(result)

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
func (h HttpHandler) Handle(res Result) {

}
