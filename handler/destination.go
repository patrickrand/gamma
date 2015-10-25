package handler

import (
	log "github.com/patrickrand/gamma/log"
	"strconv"
)

type Destination struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Subpath string `json:"subpath"`
}

func (d *Destination) Endpoint(protocol string) string {
	log.DBUG("destination", "(*Destination).Endpoint => (%s:%d%s).%s", d.Host, d.Port, d.Subpath, protocol)

	endpoint := protocol + "://" + d.Host
	if d.Port > 0 {
		endpoint = endpoint + ":" + strconv.Itoa(d.Port)
	}
	return endpoint + d.Subpath
}
