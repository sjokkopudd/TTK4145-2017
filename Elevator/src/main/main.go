package main

import (
	"def"
	"elevatorMap"
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

	go hardware.InitHardware(msgChan_toHardware, msgChan_fromHardware)

	go fsm.InitFsm(msgChan_toFsm, msgChan_fromFsm)

	go network.StartNetworkCommunication(msgChan_toNetwork, msgChan_fromNetwork, msgChan_deadElevator)

	for {
		select {
		case msg := <-msgChan_fromHardware:

			msgChan_toFsm <- msg

		case msg := <-msgChan_fromNetwork:

			receivedMap := msg.Map.(def.ElevMap)

			fsmEvent, currentMap := elevatorMap.GetEventFromNetwork(receivedMap)
			// AddNewMapChanges() skal luke ut om det er gjort en fms_trigger event
			// og returnere et event, det nye mappet og om alle er eninge

			newMsg := def.ConstructChannelMessage(currentMap, fsmEvent)

			msgChan_toHardware <- newMsg

			msgChan_toFsm <- newMsg

		case msg := <-msgChan_fromFsm:

			receivedMap := msg.Map.(def.ElevMap)

			currentMap, changemade := elevatorMap.AddNewMapChanges(receivedMap, 0)

			newMsg := def.ConstructChannelMessage(currentMap, nil)

			msgChan_toHardware <- newMsg

			if changemade {
				msgChan_toNetwork <- newMsg
			}
		}
	}
}
