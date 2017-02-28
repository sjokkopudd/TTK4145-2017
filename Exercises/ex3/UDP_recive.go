package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":30000")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		fmt.Println("Recived ", string(buffer[0:n]), " from ", addr)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}
