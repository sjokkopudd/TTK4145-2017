package main

import (
	//"bufio"
	"elevatorMap"
	"fmt"
	"network"
	//"os"
)

func main() {

	transmitMap := make(chan elevatorMap.ElevMap)
	receiveMap := make(chan elevatorMap.ElevMap)
	go elevatorMap.InitMap(transmitMap)

	go network.StartNetworkCommunication(transmitMap, receiveMap)

	for {
		fmt.Println(<-receiveMap)
	}

}
