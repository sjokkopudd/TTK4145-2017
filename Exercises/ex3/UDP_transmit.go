package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	//send
	destination_addr, err := net.ResolveUDPAddr("udp", "129.241.187.38:20005")
	if err != nil {
		log.Fatal(err)
	}

	send_conn, err := net.DialUDP("udp", nil, destination_addr)
	if err != nil {
		log.Fatal(err)
	}

	//listen
	local_addr, err := net.ResolveUDPAddr("udp", ":20005")
	if err != nil {
		log.Fatal(err)
	}

	recive_conn, err := net.ListenUDP("udp", local_addr)
	if err != nil {
		log.Fatal(err)
	}

	//close sockets at end
	defer send_conn.Close()
	defer recive_conn.Close()

	for {
		msg := "Test studass"
		fmt.Println("Sending: ", msg, "\n")
		send_conn.Write([]byte(msg))

		buffer := make([]byte, 1024)
		n, addr, err := recive_conn.ReadFromUDP(buffer)
		fmt.Println("Recived: ", string(buffer[0:n]), " from ", addr, "\n")
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(1 * time.Second)
	}
}
