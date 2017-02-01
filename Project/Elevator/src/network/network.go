package network

import (
	"fmt"
	"log"
	"net"
	"time"
)

func StartNetworkCommunication(r chan string) {

	fmt.Println("Trying to get nettwork conn")

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

		if len(buffer) > 0 {
			r <- (string(buffer[0:n]), " recived from addr: ", addr)
			if err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}
