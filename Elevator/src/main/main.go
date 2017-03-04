package main

import (
	"def"
	"elevatorMap"
	"hardware"
	//"network"
	"time"
)

func main() {

	//transmitChan := make(chan def.ElevMap)
	//receiveChan := make(chan def.ElevMap)

	eventChan := make(chan def.NewHardwareEvent)
	mapChan := make(chan def.ElevMap)
	mapChan_toHw := make(chan def.ElevMap)
	eventChan_toMap := make(chan def.NewHardwareEvent)
	//eventChan_toTH := make(chan def.NewHardwareEvent)

	go elevatorMap.InitMap(mapChan /* transmitChan, receiveChan,*/, eventChan_toMap)

	go hardware.InitHardware(mapChan_toHw, eventChan)

	for {

		select {
		case updatedMap := <-mapChan:
			mapChan_toHw <- updatedMap
		case newEvent := <-eventChan:
			eventChan_toMap <- newEvent
			//eventChan_toTH <- newEvent
		}
		time.Sleep(50*time.Millisecond)
	}
}