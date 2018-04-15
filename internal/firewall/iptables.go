package firewall

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"
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
	log.Println("IPTables.Detach Called")

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

func (i *IPTables) Add(host string, port int) error {
	// iptables -I open-says-me 1 -p tcp -s <host> --dport <port> -j ACCEPT -m comment --comment <time>
	cmd := exec.Command("iptables", "-I", "open-says-me", "1", "-p", "tcp", "-s", host, "--dport", strconv.Itoa(port), "-j", "ACCEPT", "-m", "comment", "--comment", time.Now().String())
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error adding rule for %s: %s", host, err)
	}
	return nil
}

func (i *IPTables) Remove(host string, port int) error {
	return errors.New("Not Implemented")
}

func (i *IPTables) List() ([]Exception, error) {
	return []Exception{}, errors.New("Not Implemented")
}
