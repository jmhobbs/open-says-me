package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/jmhobbs/open-says-me/internal/firewall"
)

type Knock struct {
	Address net.IP
	Port    string
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage:", os.Args[0], "<locked port>", "<port>...")
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Locking port", port)

	f := firewall.IPTables{}
	err = f.Attach()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Detach(); err != nil {
			log.Println(err)
		}
	}()

	err = f.Block(port)
	if err != nil {
		log.Fatal(err)
	}

	knock_ports := os.Args[2:]
	log.Println("Knock order:", knock_ports)

	ch := make(chan Knock)

	for _, knock_port := range knock_ports {
		go listen(ch, knock_port)
	}

	// TODO: Expiration loop here.

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Println("Shutting down...")
			// TODO: Cleanly kill UDP listeners
			close(ch)
		}
	}()

	connections := map[string]int{}

	for knock := range ch {
		key := string(knock.Address)
		current_step := connections[key]
		if knock_ports[current_step] == knock.Port {
			connections[key] = current_step + 1
		} else {
			log.Println("bad knock from:", knock.Address.String())
			delete(connections, key)
		}

		if connections[key] >= len(knock_ports) {
			if err := f.Add(knock.Address.String(), port); err != nil {
				log.Println("error adding firewall exception:", err)
			}
			delete(connections, key)
		}
	}
}

// Listens on the UDP port for a knock packet
// and then pushes those knocks up the channel
func listen(ch chan Knock, port string) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
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
			log.Println("error: ", err)
		}
		ch <- Knock{raddr.IP, port}
	}
}
