package hardware

import (
	"elevatorMap"
	"errors"
)

var buttonChannelMatrix = [elevatorMap.Floors][3]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

type NewHardwareEvent struct {
	Pos    int
	Floor  int
	Button int
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

	if floor < 0 || floor >= elevatorMap.Floors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if button < 0 || button >= 3 {
		log.Printf("Error: Button %d out of range!\n", button)
		return false
	}
	if button == UP && floor == elevatorMap.Floors-1 {
		log.Println("Button up from top floor does not exist!")
		return false
	}
	if button == DOWN && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return false
	}

	if ioReadBit(buttonChannelMatrix[floor][button]) {
		return true
	} else {
		return false
	}
}

func pollAllButtons(eventChan chan NewHardwareEvent) {

	for i := 0; i < elevatorMap.Floors; i++ {
		for j := 0; j < 3; j++ {
			if readButton(i, j) {
				e := NewHardwareEvent(-1, i, j)
				eventChan <- e
			}
		}
	}
}
