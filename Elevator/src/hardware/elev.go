package hardware

import (
	"def"
	"log"
)

// ----------------------------------------------------------
// -------------------- Inputs ------------------------------
// ----------------------------------------------------------

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

// ----------------------------------------------------------
// -------------------- Outputs -----------------------------
// ----------------------------------------------------------

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

func setFloorIndicator(floor int) {
	if floor < 0 || floor > 3 {
		IoClearBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	} else if floor == 0 {
		IoClearBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	} else if floor == 2 {
		IoSetBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	} else if floor == 1 {
		IoClearBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	} else {
		IoSetBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	}

}

func setOrderLight(f byte, b byte, val byte) {

	if val == 1 {
		IoSetBit(lightChannelMatrix[f][b])
	} else {
		IoClearBit(lightChannelMatrix[f][b])
	}

}

func setDoorLight(val int) {

	if val == 1 {
		IoSetBit(LIGHT_DOOR_OPEN)
	} else {
		IoClearBit(LIGHT_DOOR_OPEN)
	}

}
