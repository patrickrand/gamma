package handler

type Handler interface {
	Handle(monitorType, monitorName string, result interface{}) error
}
