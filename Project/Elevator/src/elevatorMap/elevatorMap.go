package elevatorMap

import (
	"def"
	"fmt"
	"time"
)

func NewMap() def.ElevMap {

	mapArray := make(def.ElevMap)

	for i := 0; i < def.Elevators; i++ {
		var temp def.ElevatorInfo
		temp.ID = i
		for j := 0; j < def.Floors; j++ {
			for k := 0; k < 3; k++ {
				temp.Buttons[j][k] = 0
			}
		}
		temp.Dir = 0
		temp.Pos = 0
		mapArray[def.IPs[i]] = &temp
	}

	return mapArray
}

func InitMap( /*passMap chan def.ElevMap, */ eventChan chan def.NewHardwareEvent) {

	mapArray := NewMap()

	WriteBackup(mapArray)

	/*passMap <- mapArray
	<-passMap
	*/
	time.Sleep(200 * time.Millisecond)

	for {
		go updateMap(eventChan)

		time.Sleep(200 * time.Millisecond)
	}

}

func updateMap(eventChan chan def.NewHardwareEvent) {

	for {

		select {
		case event := <-eventChan:
			mapArray := ReadBackup()
			if event.Pos != -1 {
				mapArray[def.MyIP].Pos = event.Pos
			} else if mapArray[def.MyIP].Buttons[event.Floor][event.Button] == 0 {
				mapArray[def.MyIP].Buttons[event.Floor][event.Button] = 1
				WriteBackup(mapArray)

				PrintMap(mapArray)
			}
		}
	}
}

func PrintMap(mapArray def.ElevMap) {
	for i := 0; i < def.Elevators; i++ {
		fmt.Println("IP: " + def.IPs[i])
		fmt.Println(*mapArray[def.IPs[i]])
		fmt.Println()

	}
}
