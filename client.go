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

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	for _, port := range ports {
		log.Println("Knocking on", port)
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", hostname, port))
		if err != nil {
			panic(err)
		}

		conn, err := net.DialUDP("udp", localAddr, addr)
		if err != nil {
			panic(err)
		}

		conn.Write([]byte(""))

		conn.Close()
	}
}
