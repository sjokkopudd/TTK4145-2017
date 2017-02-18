package elevatorMap

import (
	"hardware"
	"time"
)

const Elevators = 1
const Floors = 4

type ElevatorInfo struct {
	ID      int
	Buttons [Floors][3]int
	Dir     int
	Pos     int
}

type ElevMap [Elevators]ElevatorInfo

func NewMap() [Elevators]ElevatorInfo {

	var mapArray [Elevators]ElevatorInfo

	for i := 0; i < Elevators; i++ {
		mapArray[i].ID = i
		for j := 0; j < Floors; j++ {
			for k := 0; k < 3; k++ {
				mapArray[i].Buttons[j][k] = 0
			}
		}
		mapArray[i].Dir = 0
		mapArray[i].Pos = 0
	}

	return mapArray
}

func InitMap(passMap chan ElevMap) {

	mapArray := NewMap()

	mapArray = ReadBackup()

	passMap <- mapArray

	time.Sleep(200 * time.Millisecond)

	for {

	}

}

func updateMap(eventChan chan hardware.NewHardwareEvent) {

	for {
		select {
		case event := <-eventChan:
			mapArray := ReadBackup()
			if event.Pos {

			} else {
				mapArray[0].Buttons[event.Floor][event.Button]
			}
		}
	}
}
