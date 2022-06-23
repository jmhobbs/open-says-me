package firewall

type Exception struct {
	Host string
	Port int
}

type Firewall interface {
	Attach() error
	Detach() error
	Block(port int) error
	AddException(host string, port int) error
	RemoveException(host string, port int) error
	ListExceptions() ([]Exception, error)
}
