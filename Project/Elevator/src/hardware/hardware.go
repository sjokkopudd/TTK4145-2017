package hardware

import (
	"def"
	"errors"
	"fmt"
	"log"
	"time"
)

func InitHardware(mapChan chan def.ElevMap, eventChan chan def.NewHardwareEvent) {
	if IoInit() != true {
		log.Fatal(errors.New("Unsucsessful init of IO"))

	}
	/*SetMotorDir(-1)
	for readFloor() != 1 {

	}*/
	SetMotorDir(0)

	for {
		go pollNewEvents(eventChan)
		time.Sleep(50 * time.Millisecond)
		go setLights(mapChan)
		time.Sleep(50 * time.Millisecond)
	}
}

func SetMotorDir(dir int) {
	if dir == 0 {
		IoWriteAnalog(MOTOR, 0)
	} else if dir < 0 {
		IoSetBit(MOTORDIR)
		IoWriteAnalog(MOTOR, 2800)
	} else if dir > 0 {
		IoClearBit(MOTORDIR)
		IoWriteAnalog(MOTOR, 2800)
	}
}

func readFloor() int {
	if IoReadBit(SENSOR_FLOOR1) {
		return 0
	} else if IoReadBit(SENSOR_FLOOR2) {
		return 1
	} else if IoReadBit(SENSOR_FLOOR3) {
		return 2
	} else if IoReadBit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func readButton(floor int, button int) bool {

	if floor < 0 || floor >= def.FLOORS {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if button < 0 || button >= 3 {
		log.Printf("Error: Button %d out of range!\n", button)
		return false
	}
	if button == def.UP && floor == def.FLOORS-1 {
		log.Println("Button up from top floor does not exist!")
		return false
	}
	if button == def.DOWN && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return false
	}

	if IoReadBit(buttonChannelMatrix[floor][button]) {
		return true
	} else {
		return false
	}
}

func pollNewEvents(eventChan chan def.NewHardwareEvent) {

	newPos := readFloor()
	fmt.Println("newPos: ", newPos)

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

func setFloorIndicator(floor int) {
	fmt.Println("Floor: ", floor)

	if floor < 0 || floor > 3 {
		IoClearBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	} else if floor == 0 {
		IoClearBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	} else if floor == 1 {
		IoSetBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	} else if floor == 2 {
		IoClearBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	} else {
		IoSetBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	}
}

func setLights(mapChan chan def.ElevMap) {
	select {
	case localMap := <-mapChan:
		for b := 0; b < def.BUTTONS; b++ {
			for f := 0; f < def.FLOORS; f++ {
				setLight := true
				for e := 0; e < def.ELEVATORS; e++ {
					if localMap[def.IPs[e]].Buttons[f][b] != 1 {
						setLight = false
					}
				}
				if setLight {
					IoSetBit(lightChannelMatrix[f][b])
				} else {
					IoClearBit(lightChannelMatrix[f][b])
				}
			}
		}
		setFloorIndicator(localMap[def.MY_IP].Pos)
	}
}
