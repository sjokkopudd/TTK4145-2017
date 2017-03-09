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

	time.Sleep(50 * time.Millisecond)

	go hardware.InitHardware(msgChan_toHardware, msgChan_fromHardware)
	go fsm.Fsm(msgChan_toFsm, msgChan_fromFsm)

	go network.StartNetworkCommunication(msgChan_toNetwork, msgChan_fromNetwork, msgChan_deadElevator)

	for {
		select {
		case msg := <-msgChan_fromHardware:

			newEvent := msg.Event.(def.NewEvent)

			currentMap, changeMade, allAgree := elevatorMap.UpdateMap(newEvent)

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

			newEvent := elevatorMap.ReceivedMapFromNetwork(receivedMap)

			currentMap, changemade, allAgree := elevatorMap.UpdateMap(newEvent)

			elevatorMap.PrintMap(currentMap)

			newMsg := def.ConstructChannelMessage(currentMap, newEvent)

			msgChan_toHardware <- newMsg

			msgChan_toFsm <- newMsg

			if changemade {
				msgChan_toNetwork <- newMsg
			}
		}
	}
}
