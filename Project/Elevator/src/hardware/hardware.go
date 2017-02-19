package hardware

import (
	"def"
	"elevatorMap"
	"errors"
	"log"
	"time"
)

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

func InitHardware(eventChan chan def.NewHardwareEvent) {
	if IoInit() != true {
		log.Fatal(errors.New("Unsucsessful init of IO"))

	}
	/*SetMotorDir(-1)
	for readFloor() != 1 {

	}*/
	SetMotorDir(0)

	for {
		go pollAllButtons(eventChan)
		time.Sleep(100 * time.Millisecond)
		go setLights()
		time.Sleep(100 * time.Millisecond)
	}
}

func readButton(floor int, button int) bool {

	if floor < 0 || floor >= def.Floors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if button < 0 || button >= 3 {
		log.Printf("Error: Button %d out of range!\n", button)
		return false
	}
	if button == UP && floor == def.Floors-1 {
		log.Println("Button up from top floor does not exist!")
		return false
	}
	if button == DOWN && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return false
	}

	if IoReadBit(buttonChannelMatrix[floor][button]) {
		return true
	} else {
		return false
	}
}

func pollAllButtons(eventChan chan def.NewHardwareEvent) {

	for i := 0; i < def.Floors; i++ {
		for j := 0; j < 3; j++ {
			if !((i == 0) && (j == 1)) && !((i == def.Floors-1) && (j == 0)) {
				if readButton(i, j) {
					e := def.NewHardwareEvent{-1, i, j}
					eventChan <- e
				}
			}
		}
	}
}

func setLights() {
	mapArray := elevatorMap.ReadBackup()
	for i := 0; i < 3; i++ {
		for j := 0; j < def.Floors; j++ {
			setLight := true
			for k := 0; k < def.Elevators; k++ {
				if mapArray[def.IPs[k]].Buttons[j][i] != 1 {
					setLight = false
				}
			}
			if setLight {
				IoSetBit(lightChannelMatrix[j][i])
			} else {
				IoClearBit(lightChannelMatrix[j][i])
			}
		}
	}
}
