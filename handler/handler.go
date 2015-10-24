package handler

type Handler interface {
	Handle(data interface{}, ch chan<- bool)
}

func New(handlerType string) Handler {
	switch handlerType {
	case HTTP:
		return &HttpHandler{}
	default:
		return nil
	}
}
