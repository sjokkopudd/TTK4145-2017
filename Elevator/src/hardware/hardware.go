package hardware

import (
	"def"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var conn net.Conn
var mutex = &sync.Mutex{}

// -----------------------------------------------------------------
// ----------------------- Interface -------------------------------
// -----------------------------------------------------------------

func InitHardware(msgChan_toHW chan def.ChannelMessage, msgChan_fromHW chan def.ChannelMessage) {
	if def.USING_SIMULATOR {

		fmt.Println("Mode: def.USING_SIMULATOR")

		tcpAddr, err := net.ResolveTCPAddr("tcp", def.SIM_SERV_ADDR)
		if err != nil {
			fmt.Println("ResolveTCPAddr failed:", err.Error())
			log.Fatal(err)
		}
		fmt.Println("ResolveTCPAddr success")

		conn, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Println("Dial failed:", err.Error())
			log.Fatal(err)
		}
		fmt.Println("Dial success")

		goToNearestFloor()

		go pollNewEvents(msgChan_fromHW)

		go setLights(msgChan_toHW)

	}

	if !def.USING_SIMULATOR {

		if IoInit() != true {
			log.Fatal(errors.New("Unsucsessful init of IO"))
		}

		goToNearestFloor()

		go pollNewEvents(msgChan_fromHW)

		go setLights(msgChan_toHW)
	}
}

func goToNearestFloor() {
	if readFloor() == -1 {
		SetMotorDir(-1)
	}
	for {
		if readFloor() != -1 {
			SetMotorDir(0)
			break
		}
	}
}

// -------------------------------------------------------------------------
// ----------------------------- LOOPS -------------------------------------
// -------------------------------------------------------------------------

func setLights(msgChan_toHW chan def.ChannelMessage) {
	for {
		select {
		case msg := <-msgChan_toHW:
			currentMap := msg.Map.(def.ElevMap)
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
			setFloorIndicator(currentMap[def.MY_ID].Pos)
		}
		time.Sleep(100 * time.Millisecond)
	}

}

func pollNewEvents(msgChan_fromHW chan def.ChannelMessage) {
	lastPos := -1
	var buttonState [def.FLOORS][def.BUTTONS]bool
	for {
		newPos := readFloor()
		if (newPos != -1) && (newPos != lastPos) {
			newEvent := def.NewEvent{def.FLOOR_ARRIVAL, newPos}
			newMsg := def.ConstructChannelMessage(nil, newEvent)
			msgChan_fromHW <- newMsg
			lastPos = newPos
		}
		for f := 0; f < def.FLOORS; f++ {
			for b := 0; b < def.BUTTONS; b++ {
				if !((f == 0) && (b == 1)) && !((f == def.FLOORS-1) && (b == 0)) {
					if readButton(f, b) && buttonState[f][b] == false {
						newEvent := def.NewEvent{def.BUTTON_PUSH, []int{f, b}}
						newMsg := def.ConstructChannelMessage(nil, newEvent)
						msgChan_fromHW <- newMsg
						buttonState[f][b] = true

					} else if !readButton(f, b) {
						buttonState[f][b] = false

					}
				}
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func goUpAndDown() {
	SetMotorDir(-1)
	dir := -1

	for {
		if readFloor() == 0 && dir == -1 {
			SetMotorDir(1)
			dir = 1
		}
		if readFloor() == 3 && dir == 1 {
			SetMotorDir(-1)
			dir = -1
		}
	}
}
