package hardware

import (
	"errors"
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

func InitHardware(newEventCh chan int) (int, error) {
	if IoInit() != true {
		return -1, errors.New("Unsucsessful init of IO")
	}
	SetMotorDir(-1)
	for readFloor() != 1 {

	}
	SetMotorDir(0)

	for {
		pollButtons()
	}

	return readFloor(), nil
}

func readButton(floor int, button int) bool {

	if ioReadBit(buttonChannelArray[button]) {
		return true
	} else {
		return false
	}
}

func pollButtons(chan newEventCh) {
	for i := 0; i < 6; i++ {
		if readButton(i) {
			newEventCh <- i
		}
	}
}
