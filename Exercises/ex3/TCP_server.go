package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	//
	destination_addr, err := net.ResolveTCPAddr("tcp", "129.241.187.38:33546")
	if err != nil {
		log.Fatal(err)
	}

	server_conn, err := net.DialTCP("tcp", nil, destination_addr)
	if err != nil {
		log.Fatal(err)
	}

	connect_order := "Connect to: 10.24.38.226:20005\x00"

	//
	local_addr, err := net.ResolveTCPAddr("tcp", "10.24.38.226:20005")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", local_addr)
	if err != nil {
		log.Fatal(err)
	}

	//
	_, err = server_conn.Write([]byte(connect_order))
	if err != nil {
		log.Fatal(err)
	}

	client_conn, err := listener.AcceptTCP()
	if err != nil {
		log.Fatal(err)
	}

	for {

		msg := "Test studass\x00"
		fmt.Println("Sending message: ", msg, "\n")
		client_conn.Write([]byte(msg))

		buffer := make([]byte, 1024)
		client_conn.Read(buffer)
		fmt.Println("Recived message: ", string(buffer), "\n")

		time.Sleep(1 * time.Second)
	}

}
