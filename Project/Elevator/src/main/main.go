package main

import (
	"def"
	"elevatorMap"

	"hardware"
	"network"
)

func main() {

	transmitChan := make(chan def.ElevMap)
	receiveChan := make(chan def.ElevMap)
	eventChan := make(chan def.NewHardwareEvent)
	mapChan := make(chan def.ElevMap)

	go elevatorMap.InitMap(mapChan, transmitChan, receiveChan, eventChan)

	go network.StartNetworkCommunication(transmitChan, receiveChan)

	go hardware.InitHardware(mapChan, eventChan)

	for {
		//fmt.Println(<-receiveMap)
	}

}
