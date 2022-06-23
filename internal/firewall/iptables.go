// +build linux

package firewall

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type IPTables struct{}

func (i *IPTables) Attach() error {
	// Add our chain
	cmd := exec.Command("iptables", "-N", "open-says-me")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error creating chain: %s", err)
	}

	// Add return from our chain
	cmd = exec.Command("iptables", "-A", "open-says-me", "-j", "RETURN")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error setting return: %s", err)
	}

	// Add jump to our chain
	cmd = exec.Command("iptables", "-t", "filter", "-I", "INPUT", "1", "-j", "open-says-me")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error attaching chain: %s", err)
	}

	return nil
}

func (i *IPTables) Detach() error {
	// Delete jump to our chain
	cmd := exec.Command("iptables", "-D", "INPUT", "-j", "open-says-me")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error removing jump to chain: %s", err)
	}

	// Delete rules in our chain
	cmd = exec.Command("iptables", "-F", "open-says-me")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error removing rules in chain: %s", err)
	}

	// Delete our chain
	cmd = exec.Command("iptables", "-X", "open-says-me")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error removing chain: %s", err)
	}

	return nil
}

func (i *IPTables) Block(port int) error {
	cmd := exec.Command("iptables", "-I", "open-says-me", "1", "-p", "tcp", "--dport", strconv.Itoa(port), "-j", "REJECT", "-m", "comment", "--comment", "BLOCKED")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error adding block for %d: %s", port, err)
	}
	return nil
}

func (i *IPTables) AddException(host string, port int) error {
	// iptables -I open-says-me 1 -p tcp -s <host> --dport <port> -j ACCEPT -m comment --comment <time>
	cmd := exec.Command("iptables", "-I", "open-says-me", "1", "-p", "tcp", "-s", host, "--dport", strconv.Itoa(port), "-j", "ACCEPT", "-m", "comment", "--comment", time.Now().String())
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error adding rule for %s: %s", host, err)
	}
	return nil
}

func (i *IPTables) RemoveException(host string, port int) error {
	return errors.New("Not Implemented")
}

func (i *IPTables) ListExceptions() ([]Exception, error) {
	return nil, errors.New("Not Implemented")
}

func New(logger zerolog.Logger) Firewall {
	return &IPTables{}
}