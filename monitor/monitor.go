package monitor

type Monitor interface {
	Id() string
	Type() string
	Exec() (Result, error)
	RuntimeInterval() time.Duration
}
