package main

import (
	"def"
	"elevatorMap"
	//"fmt"
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
	msgChan_buttonEvent := make(chan def.ChannelMessage, 100)
	msgChan_fromHardware_buttons := make(chan def.ChannelMessage, 100)
	msgChan_fromHardware_floors := make(chan def.ChannelMessage, 100)
	msgChan_fromFsm := make(chan def.ChannelMessage, 100)
	elevatorMap.InitMap()

	time.Sleep(500 * time.Millisecond)

	go hardware.InitHardware(msgChan_toHardware, msgChan_fromHardware_buttons, msgChan_fromHardware_floors)

	go fsm.InitFsm(msgChan_buttonEvent, msgChan_fromHardware_floors, msgChan_fromFsm, msgChan_deadElevator)

	go network.StartNetworkCommunication(msgChan_toNetwork, msgChan_fromNetwork, msgChan_deadElevator)

	time.Sleep(500 * time.Millisecond)

	transmitTicker := time.NewTicker(100 * time.Millisecond)

	transmitFlag := false

	ligthFlag := false

	var newMsg def.ChannelMessage

	for {

		select {
		case msg := <-msgChan_fromHardware_buttons:
			msgChan_buttonEvent <- msg

		case msg := <-msgChan_fromNetwork:
			receivedMap := msg.Map.(def.ElevMap)

			fsmEvent, currentMap := elevatorMap.GetEventFromNetwork(receivedMap)

			newMsg = def.ConstructChannelMessage(currentMap, fsmEvent)

			msgChan_buttonEvent <- newMsg

			ligthFlag = true

		case msg := <-msgChan_fromFsm:
			receivedMap := msg.Map.(def.ElevMap)

			currentMap, changemade := elevatorMap.AddNewMapChanges(receivedMap, 0)

			newMsg = def.ConstructChannelMessage(currentMap, nil)

			ligthFlag = true

			if changemade {
				transmitFlag = true
			}

			/*case <-transmitTicker.C:
			if ligthFlag {
				fmt.Println("setting lights")
				msgChan_toHardware <- newMsg
				ligthFlag = false
			}
			if transmitFlag {
				fmt.Println("transmitting shit")
				msgChan_toNetwork <- newMsg
				transmitFlag = false
			}*/
		}

		if ligthFlag || transmitFlag {
			select {
			case <-transmitTicker.C:
				if ligthFlag {
					msgChan_toHardware <- newMsg
					ligthFlag = false
				}
				if transmitFlag {
					msgChan_toNetwork <- newMsg
					transmitFlag = false
				}
			}
		}
	}
}
