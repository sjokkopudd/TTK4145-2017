package main

import (
	"fmt"
	"network"
)

func main() {

	r := make(chan string)

	go network.StartNetworkCommunication()

	for {
		msg := <-r
		fmt.Println(msg)
	}
}
