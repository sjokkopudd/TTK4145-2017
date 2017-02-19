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

	go elevatorMap.InitMap(transmitChan, receiveChan, eventChan)

	go network.StartNetworkCommunication(transmitChan, receiveChan)

	go hardware.InitHardware(eventChan)

	for {
		//fmt.Println(<-receiveMap)
	}

}
