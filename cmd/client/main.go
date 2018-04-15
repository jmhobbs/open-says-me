package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("usage:", os.Args[0], "<server>", "<port>...")
		return
	}

	hostname := os.Args[1]
	ports := os.Args[2:]

	for _, port := range ports {
		address := fmt.Sprintf("%s:%s", hostname, port)
		log.Println("Knocking on", address)
		addr, err := net.ResolveUDPAddr("udp", address)
		if err != nil {
			panic(err)
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			panic(err)
		}

		conn.Write([]byte("ok"))

		conn.Close()
	}
}
