package main

import (
	"def"
	"elevatorMap"
	"fmt"
	"hardware"
	"network"
	"taskHandler"
	"time"
)

func main() {

	transmitChan := make(chan def.ElevMap, 10)
	receiveChan := make(chan def.ElevMap, 10)
	deadElevatorChan := make(chan def.NewEvent, 10)

	eventChan_fromHW := make(chan def.NewEvent)
	eventChan_fromTH := make(chan def.NewEvent)
	eventChan_toTH := make(chan def.NewEvent, 4)
	mapChan_toHW := make(chan def.ElevMap)

	elevatorMap.InitMap()

	go hardware.InitHardware(mapChan_toHW, eventChan_fromHW)
	go taskHandler.EventHandler(eventChan_toTH, eventChan_fromTH)

	go network.StartNetworkCommunication(transmitChan, receiveChan, deadElevatorChan)

	for {
		select {
		case newEvent := <-eventChan_fromHW:
			fmt.Println("NEW HW EVENT: ", newEvent)
			currentMap, changeMade := elevatorMap.UpdateMap(newEvent)
			if changeMade {
				transmitChan <- currentMap
				mapChan_toHW <- currentMap
				eventChan_toTH <- newEvent
			}

		case receivedMap := <-receiveChan:

			fmt.Println("RECIEVED MAP: ", receivedMap)
			newEvent := elevatorMap.ReceivedMapFromNetwork(receivedMap)
			currentMap, changemade := elevatorMap.UpdateMap(newEvent)
			if changemade {
				transmitChan <- currentMap
			} else {
				mapChan_toHW <- currentMap
				eventChan_toTH <- newEvent

			}
		case newEvent := <-eventChan_fromTH:
			fmt.Println("NEW TH EVENT: ", newEvent)
			currentMap, changeMade := elevatorMap.UpdateMap(newEvent)
			if changeMade {
				transmitChan <- currentMap
				mapChan_toHW <- currentMap
				eventChan_toTH <- newEvent
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}
