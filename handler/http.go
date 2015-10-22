package handler

type HttpHandler struct {
	Id           string `json:"id"`
	*Destination `json:"destination"`
	Parameters   `json:"parameters"`
}

func NewHttpHandler(id string, dest *Destination, params map[string]interface{}) *HttpHandler {
	return &HttpHandler{
		Id:          id,
		Destination: dest,
		Parameters:  params,
	}
}

func (h HttpHandler) Handle(channel <-chan bool) {

}
