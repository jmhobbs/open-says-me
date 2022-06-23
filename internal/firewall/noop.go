// +build !linux

package firewall

import (
	"errors"

	"github.com/rs/zerolog"
)

type NoOp struct {}

func (n *NoOp) Attach() error {
	return nil
}

func (n *NoOp) Detach()  error {
	return nil
}

func (n *NoOp) Block(port int) error {
	return nil
}

func (n *NoOp) AddException(host string, port int) error {
	return nil
}

func (n *NoOp) RemoveException(host string, port int) error {
	return nil
}

func (n *NoOp) ListExceptions() ([]Exception, error) {
	return nil, errors.New("Not Implemented")
}

func New(logger zerolog.Logger) Firewall {
	logger.Warn().Msg("Using no-op firewall, this will not block traffic!")
	logger.Warn().Msg("This is only intended for testing on currently unsupported platforms!")
	return &NoOp{}
}