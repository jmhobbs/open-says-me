package main

import (
	"fmt"
	"log"
	"net"
	"os"
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

	locked_port := os.Args[1]
	ports := os.Args[2:]

	ch := make(chan Knock)

	for _, port := range ports {
		go listen(ch, port)
	}

	connections := map[string]int{}

	for knock := range ch {
		log.Println("knock from", knock.Address, "for port", knock.Port)
		key := string(knock.Address)
		current_step := connections[key]
		log.Println("    step", current_step)
		if ports[current_step] == knock.Port {
			log.Println("    port matched")
			connections[key] = current_step + 1
		} else {
			log.Println("    wrong port")
			delete(connections, key)
		}

		if connections[key] >= len(ports) {
			log.Println("    OPEN SAYS ME")
			log.Println("    Unlocking", locked_port, "for", knock.Address)
			delete(connections, key)
		}
	}
}

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
		ch <- Knock{raddr.IP, port}

		if err != nil {
			log.Println("error: ", err)
		}
	}
}
