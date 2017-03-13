package main

import (
	"def"
	"elevatorMap"
	"encoding/json"
	"fmt"
	"fsm"
	"hardware"
	"log"
	"net"
	"network"
	"time"
)

func main() {

	backup := amIBackup()

	msgChan_toNetwork := make(chan def.ChannelMessage, 100)
	msgChan_fromNetwork := make(chan def.ChannelMessage, 100)
	msgChan_deadElevator := make(chan def.ChannelMessage, 100)
	msgChan_toHardware := make(chan def.ChannelMessage, 100)
	msgChan_buttonEvent := make(chan def.ChannelMessage, 100)
	msgChan_fromHardware_buttons := make(chan def.ChannelMessage, 100)
	msgChan_fromHardware_floors := make(chan def.ChannelMessage, 100)
	msgChan_fromFsm := make(chan def.ChannelMessage, 100)

	elevatorMap.InitMap(backup)

	go elevatorMap.SoftwareBackup()

	go hardware.InitHardware(msgChan_toHardware, msgChan_fromHardware_buttons, msgChan_fromHardware_floors)

	go fsm.Fsm(msgChan_buttonEvent, msgChan_fromHardware_floors, msgChan_fromFsm, msgChan_deadElevator)

	go network.StartNetworkCommunication(msgChan_toNetwork, msgChan_fromNetwork, msgChan_deadElevator)

	time.Sleep(500 * time.Millisecond)

	transmitTicker := time.NewTicker(100 * time.Millisecond)

	transmitFlag := false

	var newMsg def.ChannelMessage

	for {

		select {
		case msg := <-msgChan_fromHardware_buttons:
			msgChan_buttonEvent <- msg

		case msg := <-msgChan_fromNetwork:
			receivedMap := msg.Map.(elevatorMap.ElevMap)

			buttonPushes, currentMap := elevatorMap.GetEventFromNetwork(receivedMap)

			newMsg = def.ConstructChannelMessage(currentMap, nil)

			msgChan_toHardware <- newMsg

			for _, push := range buttonPushes {

				fsmEvent := def.NewEvent{def.BUTTON_PUSH, []int{push[0], push[1]}}

				newMsg = def.ConstructChannelMessage(currentMap, fsmEvent)

				msgChan_buttonEvent <- newMsg

			}

		case msg := <-msgChan_fromFsm:
			receivedMap := msg.Map.(elevatorMap.ElevMap)

			currentMap, changemade := elevatorMap.AddNewMapChanges(receivedMap, 0)

			newMsg = def.ConstructChannelMessage(currentMap, nil)

			msgChan_toHardware <- newMsg

			if changemade {
				transmitFlag = true
			}
		default:

		}

		select {
		case <-transmitTicker.C:
			if transmitFlag {

				if newMsg.Map != nil {
					msgChan_toNetwork <- newMsg
					transmitFlag = false
				}

			}

		}
	}
}

func amIBackup() bool {

	var msg bool

	addr, err := net.ResolveUDPAddr("udp", def.BACKUP_PORT)
	if err != nil {
		log.Fatal(err)
	}

	listenCon, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer listenCon.Close()

	buffer := make([]byte, 16)

	for {
		listenCon.SetReadDeadline(time.Now().Add(600 * time.Millisecond))
		n, _, err := listenCon.ReadFromUDP(buffer[:])
		if n > 0 {
			json.Unmarshal(buffer[0:n], &msg)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Elevator not alive, I'm taking over")
			return msg
		}
	}

}
