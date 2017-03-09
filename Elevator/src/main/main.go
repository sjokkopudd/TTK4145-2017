package main

import (
	"def"
	"elevatorMap"
	"fmt"
	"fsm"
	"hardware"
	"network"
	"time"
)

func main() {

	msgChan_toNetwork := make(chan def.ChannelMessage, 100)
	msgChan_fromNetwork := make(chan def.ChannelMessage, 100)
	msgChan_deadElevator := make(chan def.ChannelMessage, 100)
	msgChan_toHardware := make(chan def.ChannelMessage, 100)
	msgChan_fromHardware := make(chan def.ChannelMessage, 100)
	msgChan_toFsm := make(chan def.ChannelMessage, 100)
	msgChan_fromFsm := make(chan def.ChannelMessage, 100)

	elevatorMap.InitMap()

	time.Sleep(50 * time.Millisecond)

	go hardware.InitHardware(msgChan_toHardware, msgChan_fromHardware)

	go fsm.InitFsm(msgChan_toFsm, msgChan_fromFsm)

	go network.StartNetworkCommunication(msgChan_toNetwork, msgChan_fromNetwork, msgChan_deadElevator)

	for {
		select {
		case msg := <-msgChan_fromHardware:

			newEvent := msg.Event.(def.NewEvent)

			currentMap, changeMade, allAgree := elevatorMap.AddNewEvent(newEvent)

			newMsg := def.ConstructChannelMessage(currentMap, newEvent)

			msgChan_toHardware <- newMsg

			if allAgree {
				msgChan_toFsm <- newMsg
			}

			if changeMade {
				msgChan_toNetwork <- newMsg

			}

		case msg := <-msgChan_fromNetwork:

			receivedMap := msg.Map.(def.ElevMap)

			fsmEvent, currentMap, changemade, allAgree := elevatorMap.AddNewMapChanges(receivedMap, 1)
			// AddNewMapChanges() skal luke ut om det er gjort en fms_trigger event
			// og returnere et event, det nye mappet og om alle er eninge

			newMsg := def.ConstructChannelMessage(currentMap, fsmEvent)

			msgChan_toHardware <- newMsg

			if allAgree {
				msgChan_toFsm <- newMsg
			}

			if changemade {
				fmt.Println("Sending this map, came from network")
				elevatorMap.PrintMap(currentMap)
				msgChan_toNetwork <- newMsg
			}

		case msg := <-msgChan_fromFsm:

			receivedMap := msg.Map.(def.ElevMap)

			newEvent, currentMap, changemade, _ := elevatorMap.AddNewMapChanges(receivedMap, 0)

			newMsg := def.ConstructChannelMessage(currentMap, newEvent)

			if changemade {
				fmt.Println("Sending this map, came from FSM")
				elevatorMap.PrintMap(currentMap)
				msgChan_toNetwork <- newMsg
				msgChan_toHardware <- newMsg
			}
		}
	}
}
