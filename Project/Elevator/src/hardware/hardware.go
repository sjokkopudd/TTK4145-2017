package hardware

import (
	"def"
	"errors"
	"fmt"
	"log"
	"time"
	"sync"
)

var usingSimulator bool = true
var conn net.TCPConn

var mutex = &sync.Mutex{}

const(
	simServAddr := "localhost:15657"
)

func InitHardware(mapChan chan def.ElevMap, eventChan chan def.NewHardwareEvent) {
	select{
	case usingSimulator != true:
		if IoInit() != true {
			log.Fatal(errors.New("Unsucsessful init of IO"))
		}

		SetMotorDir(0)

		for {
			go pollNewEvents(eventChan)
			time.Sleep(50 * time.Millisecond)
			go setLights(mapChan)
			time.Sleep(50 * time.Millisecond)
		}
	case usingSimulator==true:

	    tcpAddr, err := net.ResolveTCPAddr("tcp", simServAddr)
	    if err != nil {
	        println("ResolveTCPAddr failed:", err.Error())
	        log.Fatal(err)
	    }

	    tempConn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
		    println("Dial failed:", err.Error())
		    log.Fatal(err)
		}

		conn = tempConn

		SetMotorDir(0)

		for {
			go setLights(mapChan)
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func SetMotorDir(dir int) {
	select{
	case usingSimulator != true:
		if dir == 0 {
			IoWriteAnalog(MOTOR, 0)
		} else if dir < 0 {
			IoSetBit(MOTORDIR)
			IoWriteAnalog(MOTOR, 2800)
		} else if dir > 0 {
			IoClearBit(MOTORDIR)
			IoWriteAnalog(MOTOR, 2800)
		}
	case usingSimulator == true:
		mutex.Lock()
		_, err = conn.Write([4]byte(1,dir))
		mutex.Unlock()
		if err != nil {
		    println("Write to server failed:", err.Error())
		    log.Fatal(err)
		}
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
					select{
					case usingSimulator != true:
						IoSetBit(lightChannelMatrix[f][b])

					case usingSimulator == true:
						mutex.Lock()
						_, err = conn.Write([4]byte(2, f, b, 1))
						mutex.Unlock()
						if err != nil {
						    println("Write to server failed:", err.Error())
						    log.Fatal(err)
						}
					}
				} else {					
					select{
					case usingSimulator != true:
						IoClearBit(lightChannelMatrix[f][b])

					case usingSimulator == true:
						mutex.Lock()
						_, err = conn.Write([4]byte(2, f, b, 0))
						mutex.Unlock()
						if err != nil {
						    println("Write to server failed:", err.Error())
						    log.Fatal(err)
						}
						
					}
				}
			}
		}
		select{
		case usingSimulator != true:
			setFloorIndicator(localMap[def.MY_IP].Pos)
		}
	}
}
