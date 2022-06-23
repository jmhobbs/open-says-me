package main

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
)

// Listens on the UDP port for a knock packet
// and then pushes those knocks up the channel
func listen(ch chan Knock, port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 2)

	for {
		_, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Error().Err(err).Send()
		}
		ch <- Knock{raddr.IP, port}
	}
}
