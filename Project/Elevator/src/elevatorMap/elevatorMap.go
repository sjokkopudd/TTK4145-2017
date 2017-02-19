package elevatorMap

import (
	"def"
	"fmt"
	"time"
)

func NewCleanMap() def.ElevMap {

	mapArray := make(def.ElevMap)

	for i := 0; i < def.Elevators; i++ {
		var temp def.ElevatorInfo
		temp.IP = def.IPs[i]
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

func InitMap(transmitChan chan def.ElevMap, receiveChan chan def.ElevMap, eventChan chan def.NewHardwareEvent) {

	mapArray := NewCleanMap()

	WriteBackup(mapArray)

	time.Sleep(200 * time.Millisecond)

	go updateMap(transmitChan, eventChan)
	go receivedMap(receiveChan)

}

func updateMap(transmitChan chan def.ElevMap, eventChan chan def.NewHardwareEvent) {

	for {
		select {
		case event := <-eventChan:
			mapArray := ReadBackup()
			if event.Pos != -1 {
				mapArray[def.MyIP].Pos = event.Pos
			} else if mapArray[def.MyIP].Buttons[event.Floor][event.Button] == 0 {
				mapArray[def.MyIP].Buttons[event.Floor][event.Button] = 1
				WriteBackup(mapArray)
				transmitChan <- mapArray
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func receivedMap(receiveChan chan def.ElevMap) {
	for {
		select {
		case receivedMap := <-receiveChan:

			oldMap := ReadBackup()
			newMap := NewCleanMap()

			for i := 0; i < def.Elevators; i++ {
				for j := 0; j < def.Floors; j++ {
					for k := 0; k < 3; k++ {
						if receivedMap[def.IPs[i]].Buttons[j][k] != oldMap[def.IPs[i]].Buttons[j][k] {
							newMap[def.MyIP].Buttons[j][k] = receivedMap[def.IPs[i]].Buttons[j][k]
							newMap[def.IPs[i]].Buttons[j][k] = receivedMap[def.IPs[i]].Buttons[j][k]
						} else {
							newMap[def.MyIP].Buttons[j][k] = 0
							newMap[def.IPs[i]].Buttons[j][k] = 0
						}
					}
				}
				if receivedMap[def.IPs[i]].Dir != oldMap[def.IPs[i]].Dir {
					newMap[def.IPs[i]].Dir = receivedMap[def.IPs[i]].Dir
				}
				if receivedMap[def.IPs[i]].Pos != oldMap[def.IPs[i]].Pos {
					newMap[def.IPs[i]].Pos = receivedMap[def.IPs[i]].Pos
				}
			}

			WriteBackup(newMap)

		}
		time.Sleep(200 * time.Millisecond)
	}
}

func printMap(mapArray def.ElevMap) {
	for i := 0; i < def.Elevators; i++ {
		fmt.Println("IP: " + def.IPs[i])
		fmt.Println(*mapArray[def.IPs[i]])
		fmt.Println()

	}
}
