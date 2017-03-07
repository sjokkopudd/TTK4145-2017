package main

import (
	"def"
	"elevatorMap"
	"fmt"
	"hardware"
	"network"
	//"taskHandler"
	"time"
)

func main() {

	transmitChan := make(chan def.ElevMap, 100)
	receiveChan := make(chan def.ElevMap, 100)
	deadElevatorChan := make(chan def.NewEvent, 100)

	eventChan_fromHW := make(chan def.NewEvent, 100)
	eventChan_fromTH := make(chan def.NewEvent, 100)
	//eventChan_toTH := make(chan def.NewEvent, 4)
	mapChan_toHW := make(chan def.ElevMap)

	elevatorMap.InitMap()

	go hardware.InitHardware(mapChan_toHW, eventChan_fromHW)
	//go taskHandler.EventHandler(eventChan_toTH, eventChan_fromTH)

	go network.StartNetworkCommunication(transmitChan, receiveChan, deadElevatorChan)

	for {
		select {
		case newEvent := <-eventChan_fromHW:
			fmt.Println("FROM HW CHAN")
			elevatorMap.PrintEvent(newEvent)
			currentMap, changeMade := elevatorMap.UpdateMap(newEvent)
			if changeMade {
				transmitChan <- currentMap
				mapChan_toHW <- currentMap
				//eventChan_toTH <- newEvent
			}

		case receivedMap := <-receiveChan:
			fmt.Println("FROM RECEIVE CHAN")
			newEvent := elevatorMap.ReceivedMapFromNetwork(receivedMap)
			elevatorMap.PrintEvent(newEvent)
			currentMap, changemade := elevatorMap.UpdateMap(newEvent)
			elevatorMap.PrintMap(currentMap)
			if changemade {
				transmitChan <- currentMap
			} else {
				mapChan_toHW <- currentMap
				//eventChan_toTH <- newEvent

			}
		case newEvent := <-eventChan_fromTH:
			fmt.Println("NEW TH EVENT: ", newEvent)
			currentMap, changeMade := elevatorMap.UpdateMap(newEvent)
			if changeMade {
				transmitChan <- currentMap
				mapChan_toHW <- currentMap
				//eventChan_toTH <- newEvent
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}
