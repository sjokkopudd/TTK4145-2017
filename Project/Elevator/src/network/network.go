package network

import (
	"elevatorMap"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func StartNetworkCommunication(s chan elevatorMap.ElevMap, r chan elevatorMap.ElevMap) {

	fmt.Println("Trying to setup nettwork connection")
	go transmitMap(s)
	go reciveMap(r)

}

func transmitMap(mapArray chan elevatorMap.ElevMap) {

	destination_addr, err := net.ResolveUDPAddr("udp", "129.241.187.143:20005")
	if err != nil {
		log.Fatal(err)
	}

	send_conn, err := net.DialUDP("udp", nil, destination_addr)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	for {

		msg := <-mapArray
		json_buffer, _ := json.Marshal(msg)

		if len(json_buffer) > 0 {
			send_conn.Write(json_buffer)
			var m elevatorMap.ElevMap
			err = json.Unmarshal(json_buffer, &m)
			if err != nil {
				log.Fatal(err)
			}

		}

		time.Sleep(200 * time.Millisecond)
	}

}

func reciveMap(r chan elevatorMap.ElevMap) {

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
		recive_conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, _, err := recive_conn.ReadFromUDP(r_buffer)

		if err != nil {
			log.Fatal(err)
		}

		if n > 0 {

			var m elevatorMap.ElevMap
			err = json.Unmarshal(r_buffer[0:n], &m)
			if err != nil {
				log.Fatal(err)
			}
			r <- m

		}
	}
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
