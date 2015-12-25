package agent

type Host struct {
	ID   string `json:"id"`
	FQDN string `json:"fqdn"`
	IP   string `json:"ip"`
	SSL  bool   `json:"ssl"`
}

var host = new(Host)
