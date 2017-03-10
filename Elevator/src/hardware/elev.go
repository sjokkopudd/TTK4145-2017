package hardware

import (
	"def"
	"fmt"
	"log"
	"time"
)

// ----------------------------------------------------------
// -------------------- Inputs ------------------------------
// ----------------------------------------------------------

func readFloor() int {
	if def.USING_SIMULATOR {
		mutex.Lock()
		_, err := conn.Write([]byte{7, byte(0), byte(0), byte(0)})
		mutex.Unlock()
		if err != nil {
			fmt.Println("Write to server failed:", err.Error())
			log.Fatal(err)
		}

		buffer := make([]byte, 4)
		mutex.Lock()
		conn.Read(buffer)
		mutex.Unlock()
		time.Sleep(5 * time.Millisecond)
		if buffer[1] == 1 {
			return int(buffer[2])
		} else {
			return -1
		}
	}

	if !def.USING_SIMULATOR {
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
	return -1
}

func readButton(floor int, button int) bool {
	if def.USING_SIMULATOR {
		mutex.Lock()
		_, err := conn.Write([]byte{6, byte(button), byte(floor), byte(0)})
		mutex.Unlock()
		if err != nil {
			fmt.Println("Write to server failed:", err.Error())
			log.Fatal(err)
		}

		buffer := make([]byte, 4)
		mutex.Lock()
		conn.Read(buffer)
		mutex.Unlock()
		time.Sleep(5 * time.Millisecond)
		if buffer[1] == 1 {
			return true
		} else {
			return false
		}
	}

	if !def.USING_SIMULATOR {
		if floor < 0 || floor >= def.FLOORS {
			log.Printf("Error: Floor %d out of range!\n", floor)
			return false
		}
		if button < 0 || button >= 3 {
			log.Printf("Error: Button %d out of range!\n", button)
			return false
		}
		if button == def.UP_BUTTON && floor == def.FLOORS-1 {
			log.Println("Button up from top floor does not exist!")
			return false
		}
		if button == def.DOWN_BUTTON && floor == 0 {
			log.Println("Button down from ground floor does not exist!")
			return false
		}

		if IoReadBit(buttonChannelMatrix[floor][button]) {
			return true
		} else {
			return false
		}
	}
	return false
}

// ----------------------------------------------------------
// -------------------- Outputs -----------------------------
// ----------------------------------------------------------

func SetMotorDir(dir int) {
	if def.USING_SIMULATOR {
		mutex.Lock()
		_, err := conn.Write([]byte{1, byte(dir), byte(0), byte(0)})
		mutex.Unlock()
		time.Sleep(5 * time.Millisecond)
		if err != nil {
			fmt.Println("Write to server failed:", err.Error())
			log.Fatal(err)
		}
	}

	if !def.USING_SIMULATOR {
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
}

func setFloorIndicator(floor int) {
	if def.USING_SIMULATOR {
		mutex.Lock()
		_, err := conn.Write([]byte{3, byte(floor), byte(0), byte(0)})
		mutex.Unlock()
		time.Sleep(5 * time.Millisecond)
		if err != nil {
			fmt.Println("Write to server failed:", err.Error())
			log.Fatal(err)
		}
	}

	if !def.USING_SIMULATOR {
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

}

func setOrderLight(f byte, b byte, val byte) {

	if def.USING_SIMULATOR {
		mutex.Lock()
		_, err := conn.Write([]byte{2, byte(b), byte(f), byte(val)})
		mutex.Unlock()
		time.Sleep(5 * time.Millisecond)
		if err != nil {
			fmt.Println("Write to server failed:", err.Error())
			log.Fatal(err)
		}
	}
	if !def.USING_SIMULATOR {
		if val == 1 {
			IoSetBit(lightChannelMatrix[f][b])
		} else {
			IoClearBit(lightChannelMatrix[f][b])
		}
	}

}

func SetDoorLight(val int) {

	if def.USING_SIMULATOR {
		mutex.Lock()
		_, err := conn.Write([]byte{4, byte(val), byte(0), byte(0)})
		mutex.Unlock()
		time.Sleep(5 * time.Millisecond)
		if err != nil {
			fmt.Println("Write to server failed:", err.Error())
			log.Fatal(err)
		}
	}
	if !def.USING_SIMULATOR {
		if val == 1 {
			IoSetBit(LIGHT_DOOR_OPEN)
		} else {
			IoClearBit(LIGHT_DOOR_OPEN)
		}
	}

}
