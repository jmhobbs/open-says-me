package main

import (
	"net"
	"os"
	"os/signal"

	"github.com/jmhobbs/open-says-me/internal/firewall"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Knock struct {
	Address net.IP
	Port    int
}

func main() {
	config := getConfig()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if(config.Pretty) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Int("port", config.Port).Msg("locking")

	f := firewall.New(log.Logger)
	err := f.Attach()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer func() {
		if err := f.Detach(); err != nil {
			log.Error().Err(err).Msg("unable to detach")
		}
	}()

	err = f.Block(config.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to lock port")
	}

	log.Info().Ints("knock order", config.KnockPorts).Send()

	ch := make(chan Knock)

	for _, knock_port := range config.KnockPorts {
		go listen(ch, knock_port)
	}

	// TODO: Expiration loop here.

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Info().Msg("shutting down...")
			// TODO: Cleanly kill UDP listeners
			close(ch)
		}
	}()

	connections := map[string]int{}

	for knock := range ch {
		address := knock.Address.String()
		log.Debug().Str("address", address).Int("port", knock.Port).Msg("knock, knock")

		current_step := connections[address]
		if config.KnockPorts[current_step] == knock.Port {
			connections[address] = current_step + 1
		} else {
			log.Warn().
				Str("address", address).
				Int("expected_port", config.KnockPorts[current_step]).
				Int("received_port", knock.Port).
				Int("step", current_step).
				Msg("bad knock")
			// todo: re-lock?
			// todo: temporary ban from knocking?
			delete(connections, address)
		}

		if connections[address] >= len(config.KnockPorts) {
			if err := f.AddException(knock.Address.String(), config.Port); err != nil {
				log.Error().Err(err).Msg("error adding firewall exception")
			} else {
				log.Info().Str("address", knock.Address.String()).Msg("exception added")
			}
			delete(connections, address)
		}
	}
}
