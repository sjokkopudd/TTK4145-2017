package main

import (
	"def"
	"elevatorMap"

	"hardware"
	//"network"
)

func main() {

	//transmitMap := make(chan def.ElevMap)
	//receiveMap := make(chan elevatorMap.ElevMap)
	eventChan := make(chan def.NewHardwareEvent)

	go elevatorMap.InitMap( /*transmitMap, */ eventChan)

	//go network.StartNetworkCommunication(transmitMap, receiveMap)

	go hardware.InitHardware(eventChan)

	for {
		//fmt.Println(<-receiveMap)
	}

}
