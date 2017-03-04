package hardware

import (
	"def"
	"errors"
	"fmt"
	"log"
	"time"
)

// -----------------------------------------------------------------
// ----------------------- Interface -------------------------------
// -----------------------------------------------------------------

func InitHardware(mapChan chan def.ElevMap, eventChan chan def.NewHardwareEvent) {
	if IoInit() != true {
		log.Fatal(errors.New("Unsucsessful init of IO"))

	}

	SetMotorDir(0)

	go setLights(mapChan)

	go pollNewEvents(eventChan)

	//go goUpAndDown()
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
		time.Sleep(100 * time.Millisecond)
	}

}

func pollNewEvents(eventChan chan def.NewHardwareEvent) {
	for {
		newPos := readFloor()
		fmt.Println("newPos ", newPos)
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
		time.Sleep(100 * time.Millisecond)

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
		time.Sleep(50 * time.Millisecond)
	}
}
