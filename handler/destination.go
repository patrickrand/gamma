package handler

type Destination struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Subpath string `json:"subpath"`
}

func NewDestination(host string, port int, subpath string) *Destination {
	return &Destination{
		Host:    host,
		Port:    port,
		Subpath: subpath,
	}
}

func (d *Destination) Endpoint(protocol string) string {
	endpoint := protocol + "://" + d.Host
	if d.Port > 0 {
		endpoint = endpoint + ":" +  strconv.Itoa(d.Port))
	}
	return endpoint + d.Subpath
}
