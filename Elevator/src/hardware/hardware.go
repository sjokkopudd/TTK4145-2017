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

var USING_SIMULATOR bool = false

const (
	simServAddr = "127.0.0.1:15657"
)

var conn net.Conn
var mutex = &sync.Mutex{}

// -----------------------------------------------------------------
// ----------------------- Interface -------------------------------
// -----------------------------------------------------------------

func InitHardware(mapChan chan def.ElevMap, eventChan chan def.NewHardwareEvent) {
	if USING_SIMULATOR {

		fmt.Println("Mode: USING_SIMULATOR")

		tcpAddr, err := net.ResolveTCPAddr("tcp", simServAddr)
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

		go setLights(mapChan)

		go pollNewEvents(eventChan)

		go goUpAndDown()

	}

	if !USING_SIMULATOR {
		if IoInit() != true {
			log.Fatal(errors.New("Unsucsessful init of IO"))

		}

		go setLights(mapChan)

		go pollNewEvents(eventChan)

		go goUpAndDown()
	}
}

// -------------------------------------------------------------------------
// ----------------------------- LOOPS -------------------------------------
// -------------------------------------------------------------------------

func setLights(mapChan chan def.ElevMap) {
	for {
		select {
		case localMap := <-mapChan:
			for b := 0; b < def.BUTTONS; b++ {
				for f := 0; f < def.FLOORS; f++ {
					ligthVal := 1
					for e := 0; e < def.ELEVATORS; e++ {
						if localMap[def.IPs[e]].Buttons[f][b] != 1 {
							ligthVal = 0
						}
					}

					setOrderLight(byte(f), byte(b), byte(ligthVal))

				}
			}
			setFloorIndicator(localMap[def.MY_IP].Pos)
		}
		time.Sleep(10 * time.Millisecond)
	}

}

func pollNewEvents(eventChan chan def.NewHardwareEvent) {
	for {
		newPos := readFloor()
		for f := 0; f < def.FLOORS; f++ {
			for b := 0; b < def.BUTTONS; b++ {
				if !((f == 0) && (b == 1)) && !((f == def.FLOORS-1) && (b == 0)) {
					if readButton(f, b) {
						e := def.NewHardwareEvent{newPos, f, b}
						eventChan <- e
					} else if newPos != -1 {
						e := def.NewHardwareEvent{newPos, -1, -1}
						eventChan <- e
					}
				}
			}
		}
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
