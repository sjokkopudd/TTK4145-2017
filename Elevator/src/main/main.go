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

	time.Sleep(500 * time.Millisecond)

	go hardware.InitHardware(msgChan_toHardware, msgChan_fromHardware)

	go fsm.InitFsm(msgChan_toFsm, msgChan_fromFsm)

	go network.StartNetworkCommunication(msgChan_toNetwork, msgChan_fromNetwork, msgChan_deadElevator)

	time.Sleep(500 * time.Millisecond)

	transmitTicker := time.NewTicker(100 * time.Millisecond)

	transmitFlag := false

	ligthFlag := false

	var newMsg def.ChannelMessage

	for {

		select {
		case msg := <-msgChan_fromHardware:
			fmt.Println("got hardware event")

			msgChan_toFsm <- msg

		case msg := <-msgChan_fromNetwork:

			fmt.Println("got network event")

			receivedMap := msg.Map.(def.ElevMap)

			fsmEvent, currentMap := elevatorMap.GetEventFromNetwork(receivedMap)

			newMsg = def.ConstructChannelMessage(currentMap, fsmEvent)

			msgChan_toFsm <- newMsg

			ligthFlag = true

		case msg := <-msgChan_fromFsm:
			fmt.Println("got fsm event")

			receivedMap := msg.Map.(def.ElevMap)

			currentMap, changemade := elevatorMap.AddNewMapChanges(receivedMap, 0)

			newMsg = def.ConstructChannelMessage(currentMap, nil)

			ligthFlag = true

			if changemade {
				transmitFlag = true
			}
		case msg := <-msgChan_deadElevator:
			fmt.Println("got dead event")
			msgChan_toFsm <- msg

		case <-transmitTicker.C:
			fmt.Println("tickerout")
			if ligthFlag {
				fmt.Println("setting lights")
				msgChan_toHardware <- newMsg
				ligthFlag = false
			}
			if transmitFlag {
				fmt.Println("transmitting shit")
				msgChan_toNetwork <- newMsg
				fmt.Println("am i full?")
				transmitFlag = false
			}
		}

		/*if ligthFlag {
			select {
			case <-transmitTicker.C:
				fmt.Println("setting lights")
				msgChan_toHardware <- newMsg
				ligthFlag = false
			}
		}
		if transmitFlag {
			select {
			case <-transmitTicker.C:
				fmt.Println("transmitting shit")
				msgChan_toNetwork <- newMsg
				fmt.Println("am i full?")
				transmitFlag = false
			}
		}*/

	}
}
