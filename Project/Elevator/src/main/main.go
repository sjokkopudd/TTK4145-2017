package main

import (
	"fmt"
	"network"
)

func main() {

	r := make(chan string)

	s := make(chan string)

	go network.StartNetworkCommunication(r, s)

	for {

		fmt.Println("Send a message:")
		s_msg := "This is a message"
		fmt.Println(s_msg)
		s <- s_msg

		msg := <-r
		fmt.Println("Recived message: ", msg)
	}
}
