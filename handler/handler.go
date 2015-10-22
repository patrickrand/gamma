package handler

type Handler interface {
	Handle(channel <-chan bool)
}
