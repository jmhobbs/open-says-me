package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/peterbourgon/ff/v3"
)

type knockPortFlag struct {
	Ports []int
}

func (i *knockPortFlag) String() string {
	buf := []string{}
	for _, port := range i.Ports {
		buf = append(buf, strconv.Itoa(port))
	}
	return strings.Join(buf, ", ")
}

func (i *knockPortFlag) Set(value string) error {
	port, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	i.Ports = append(i.Ports, port)
	return nil
}

type config struct {
	Port int
	KnockPorts []int
	Debug bool
	Pretty bool
}

func getConfig() config {
	fs := flag.NewFlagSet("open-says-me", flag.ExitOnError)
	var (
		port = fs.Int("port", 9000, "port to protect")
		pretty = fs.Bool("pretty", false, "pretty print logs")
		debug      = fs.Bool("debug", false, "log debug information")
		_          = fs.String("config", "", "config file (optional)")
	)
	knockFlag := knockPortFlag{}
	fs.Var(&knockFlag, "knock", "knock port (multiple supported)")

	err := ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	)
	if err != nil {
		fmt.Fprintf(fs.Output(), "error: %v\n\n", err)
		fs.Usage()
		os.Exit(1);
	}

	return config{
		Port: *port,
		KnockPorts: knockFlag.Ports,
		Debug: *debug,
		Pretty: *pretty,
	}
}