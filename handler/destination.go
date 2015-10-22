package handler

import (
	"fmt"
)

type Destination struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func NewDestination(host string, port int) *Destination {
	return &Destination{
		Host: host,
		Port: port,
	}
}

func (d *Destination) String() string {
	return fmt.Sprintf("%s:%d", d.Host, d.Port)
}
