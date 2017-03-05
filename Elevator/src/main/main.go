package main

import (
	"def"
	"elevatorMap"
	//"hardware"
	"network"
	"time"
	//"taskHandler"
)

func main() {

	transmitChan := make(chan def.ElevMap)
	receiveChan := make(chan def.ElevMap)

	//eventChan := make(chan def.NewHardwareEvent)
	mapChan := make(chan def.ElevMap)
	//mapChan_toHw := make(chan def.ElevMap)
	//eventChan_toMap := make(chan def.NewHardwareEvent)
	//eventChan_toTH := make(chan def.NewHardwareEvent)
	//openDoorChan := make(chan int)

	go elevatorMap.InitMap(mapChan,  transmitChan, receiveChan)//, eventChan_toMap)

	//go hardware.InitHardware(mapChan_toHw, eventChan)

	//go taskHandler.taskHandler(eventChan_toTH, openDoorChan)

	go network.StartNetworkCommunication(transmitChan,receiveChan)

	for {
/*
		select {
		case updatedMap := <-mapChan:
			mapChan_toHw <- updatedMap
		case newEvent := <-eventChan:
			eventChan_toMap <- newEvent
			eventChan_toTH <- newEvent
		case doorState := <-openDoorChan:
			newDoorEvent := def.NewHardwareEvent{def.DOOR, -1, -1, -1, doorState}
			eventChan_toMap <- newDoorEvent
			eventChan_toTH <- newDoorEvent
		}
*/
		time.Sleep(50*time.Millisecond)
	}
}