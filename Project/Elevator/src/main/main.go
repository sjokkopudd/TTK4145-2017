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



	go elevatorMap.InitMap(mapChan/*, transmitChan, receiveChan*/, eventChan)

	time.Sleep(100 * time.Millisecond)

	go hardware.InitHardware(mapChan, eventChan)

	for{

	time.Sleep(100 * time.Millisecond)
		
	}
}
