package firewall

type Exception struct {
	Host string
	Port int
}

type Firewall interface {
	Attach() error
	Detach() error
	Block(port int) error
	Add(host string, port int) error
	Remove(host string, port int) error
	List() ([]Exception, error)
}
