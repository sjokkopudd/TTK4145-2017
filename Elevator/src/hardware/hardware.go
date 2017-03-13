package hardware

import (
	"def"
	"elevatorMap"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var conn *net.Conn
var mutex = &sync.Mutex{}

func InitHardware(msgChan_toHW chan def.ChannelMessage, msgChan_fromHW_buttons chan def.ChannelMessage, msgChan_fromHW_floors chan def.ChannelMessage) {
	if def.USING_SIMULATOR {

		conn = new(net.Conn)

		fmt.Println("Mode: def.USING_SIMULATOR")

		tcpAddr, err := net.ResolveTCPAddr("tcp", def.SIM_SERV_ADDR)
		if err != nil {
			fmt.Println("ResolveTCPAddr failed:", err.Error())
			log.Fatal(err)
		}
		fmt.Println("ResolveTCPAddr success")

		*conn, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Println("Dial failed:", err.Error())
			log.Fatal(err)
		}
		fmt.Println("Dial success")

		GoToNearestFloor()

		go pollButtonEvents(msgChan_fromHW_buttons)

		go pollFloorEvents(msgChan_fromHW_floors)

		go setLights(msgChan_toHW)

	}

	if !def.USING_SIMULATOR {

		if InitIO() != true {
			log.Fatal(errors.New("Unsucsessful init of IO"))
		}

		GoToNearestFloor()

		go pollButtonEvents(msgChan_fromHW_buttons)

		go pollFloorEvents(msgChan_fromHW_floors)

		go setLights(msgChan_toHW)
	}
}

func setLights(msgChan_toHW chan def.ChannelMessage) {
	for {
		select {
		case msg := <-msgChan_toHW:
			currentMap := msg.Map.(elevatorMap.ElevMap)
			for b := 0; b < def.BUTTONS; b++ {
				for f := 0; f < def.FLOORS; f++ {
					lightVal := 1
					if b != def.PANEL_BUTTON {
						for e := 0; e < def.ELEVATORS; e++ {
							if (currentMap[e].Buttons[f][b] != 1) && (currentMap[e].IsAlive == 1) {
								lightVal = 0
							}

						}
					} else {
						if currentMap[def.MY_ID].Buttons[f][b] != 1 {
							lightVal = 0
						}
					}

					setOrderLight(byte(f), byte(b), byte(lightVal))

				}
			}
			setFloorIndicator(currentMap[def.MY_ID].Position)
		}
		time.Sleep(10 * time.Millisecond)
	}

}

func pollButtonEvents(msgChan_fromHW_buttons chan def.ChannelMessage) {
	buttonTicker := time.NewTicker(200 * time.Millisecond)

	var buttonState [def.FLOORS][def.BUTTONS]bool
	for {
		select {
		case <-buttonTicker.C:
			for f := 0; f < def.FLOORS; f++ {
				for b := 0; b < def.BUTTONS; b++ {
					if !((f == 0) && (b == 1)) && !((f == def.FLOORS-1) && (b == 0)) {
						if readButton(f, b) && buttonState[f][b] == false {
							newEvent := def.NewEvent{def.BUTTON_PUSH, []int{f, b}}
							newMsg := def.ConstructChannelMessage(nil, newEvent)
							msgChan_fromHW_buttons <- newMsg
							buttonState[f][b] = true

						} else if !readButton(f, b) {
							buttonState[f][b] = false

						}
					}
				}
			}
		}
	}
}

func pollFloorEvents(msgChan_fromHW_floors chan def.ChannelMessage) {
	lastPos := -1
	floorTicker := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case <-floorTicker.C:
			newPos := readFloor()
			if (newPos != -1) && (newPos != lastPos) {
				newEvent := def.NewEvent{def.FLOOR_ARRIVAL, newPos}
				newMsg := def.ConstructChannelMessage(nil, newEvent)
				msgChan_fromHW_floors <- newMsg
				lastPos = newPos
			}
		}
	}
}
