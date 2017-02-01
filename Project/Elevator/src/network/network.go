package network

import (
	"fmt"
	"log"
	"net"
	"time"
)

func StartNetworkCommunication(r chan string, s chan string) {

	fmt.Println("Trying to setup nettwork connection")
	go recive(r)
	go transmit(s)

}

func transmit(s chan string) {

	destination_addr, err := net.ResolveUDPAddr("udp", "129.241.187.141:20005")
	if err != nil {
		log.Fatal(err)
	}

	send_conn, err := net.DialUDP("udp", nil, destination_addr)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	for {
		s_buffer := <-s

		if len(s_buffer) > 0 {
			send_conn.Write([]byte(s_buffer))
			fmt.Println("Message send: ", s_buffer)
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func recive(r chan string) {

	local_addr, err := net.ResolveUDPAddr("udp", ":20005")
	if err != nil {
		log.Fatal(err)
	}
	recive_conn, err := net.ListenUDP("udp", local_addr)
	if err != nil {
		log.Fatal(err)
	}

	defer recive_conn.Close()

	r_buffer := make([]byte, 1024)

	time.Sleep(500 * time.Millisecond)

	for {
		recive_conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		n, _, err := recive_conn.ReadFromUDP(r_buffer)

		if n > 0 {
			r <- (string(r_buffer[0:n]))
			if err != nil {
				log.Fatal(err)
			}

		}
	}
}
