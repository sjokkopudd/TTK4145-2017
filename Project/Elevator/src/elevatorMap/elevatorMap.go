package elevatorMap

import (
	"time"
)

const elevators = 3
const floors = 4

type FloorButtons struct {
	ButtonUp    int
	ButtonDown  int
	ButtonPanel int
}

type ElevatorInfo struct {
	ID      int
	Buttons [floors]FloorButtons
	Dir     int
	Pos     int
}

type ElevMap [elevators]ElevatorInfo

func NewMap() [elevators]ElevatorInfo {

	var mapArray [elevators]ElevatorInfo

	for i := 0; i < elevators; i++ {
		mapArray[i].ID = i
		for j := 0; j < floors; j++ {
			mapArray[i].Buttons[j].ButtonUp = j
			mapArray[i].Buttons[j].ButtonDown = j
			mapArray[i].Buttons[j].ButtonPanel = j
		}
		mapArray[i].Dir = 0
		mapArray[i].Pos = 0
	}

	return mapArray
}

func InitMap(passMap chan ElevMap) {
	for {
		mapArray := NewMap()

		mapArray = ReadBackup()

		passMap <- mapArray

		time.Sleep(200 * time.Millisecond)

	}

}
