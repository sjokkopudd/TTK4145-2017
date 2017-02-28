package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	destination_addr, err := net.ResolveTCPAddr("tcp", "129.241.187.43:33546")
	if err != nil {
		log.Fatal(err)
	}

	send_conn, err := net.DialTCP("tcp", nil, destination_addr)
	if err != nil {
		log.Fatal(err)
	}

	defer send_conn.Close()

	for {

		msg := "This is a message \x00"
		fmt.Println("Sending message: ", msg, "\n")
		send_conn.Write([]byte(msg))

		buffer := make([]byte, 1024)
		send_conn.Read(buffer)
		fmt.Println("Recived message: ", string(buffer), "\n")

		time.Sleep(1 * time.Second)
	}

}
