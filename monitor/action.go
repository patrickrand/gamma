package monitor

type Action struct {
	Command  string `json:"command"`
	Resource string `json:"resource"`
}

func NewAction(command, resource string) *Action {
	return &Action{Command: command, Resource: resource}
}
