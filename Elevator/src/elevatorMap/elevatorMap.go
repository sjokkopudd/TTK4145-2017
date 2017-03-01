package elevatorMap

import (
	"def"
	"fmt"
)

var localMap def.ElevMap

func NewCleanMap() def.ElevMap {

	newMap := make(def.ElevMap)

	for i := 0; i < def.ELEVATORS; i++ {
		var temp def.ElevatorInfo
		temp.IP = def.IPs[i]
		for j := 0; j < def.FLOORS; j++ {
			for k := 0; k < 3; k++ {
				temp.Buttons[j][k] = 0
			}
		}
		temp.Dir = 0
		temp.Pos = 0
		temp.Door = 0
		newMap[def.IPs[i]] = &temp
	}

	return newMap
}

func InitMap(mapChan chan def.ElevMap /*transmitChan chan def.ElevMap, receiveChan chan def.ElevMap,*/, eventChan chan def.NewHardwareEvent) {

	localMap = NewCleanMap()

	WriteBackup(localMap)

	localMap = ReadBackup()

	mapChan <- localMap

	go updateMap(mapChan /*transmitChan, receiveChan,*/, eventChan)

}

func updateMap(mapChan chan def.ElevMap /* transmitChan chan def.ElevMap, receiveChan chan def.ElevMap, */, eventChan chan def.NewHardwareEvent) {

	for {
		select {
		case event := <-eventChan:
			changeMade := false
			if (event.Pos != -1) && (localMap[def.MY_IP].Pos != event.Pos) {
				localMap[def.MY_IP].Pos = event.Pos
				changeMade = true

			}

			if (event.Floor != -1) && (localMap[def.MY_IP].Buttons[event.Floor][event.Button] == 0) {
				localMap[def.MY_IP].Buttons[event.Floor][event.Button] = 1
				changeMade = true
			}

			if changeMade {
				WriteBackup(localMap)
				//transmitChan <- localMap
				mapChan <- localMap
			}
			/*	case receivedMap := <-receiveChan:
				changeMade := false
				localMap := ReadBackup()
				for e := 0; e < def.ELEVATORS; e++ {
					for f := 0; f < def.FLOORS; f++ {
						for b := 0; b < def.BUTTONS; b++ {
							if receivedMap[def.IPs[e]].Buttons[f][b] == 1 && localMap[def.IPs[e]].Buttons[f][b] != 1 {
								localMap[def.IPs[e]].Buttons[f][b] = 1
								if b != def.PANEL {
									localMap[def.MY_IP].Buttons[f][b] = 1
								}
								changeMade = true
							}
						}
					}
					if receivedMap[def.IPs[e]].Dir != localMap[def.IPs[e]].Dir {
						localMap[def.IPs[e]].Dir = receivedMap[def.IPs[e]].Dir
						changeMade = true
					}
					if receivedMap[def.IPs[e]].Pos != localMap[def.IPs[e]].Pos {
						localMap[def.IPs[e]].Pos = receivedMap[def.IPs[e]].Pos
						changeMade = true
					}
				}

				if changeMade {
					WriteBackup(localMap)
					transmitChan <- localMap
					mapChan <- localMap
				}*/

		}
	}

}

func printMap(localMap def.ElevMap) {
	for e := 0; e < def.ELEVATORS; e++ {
		fmt.Println("IP: " + def.IPs[e])
		fmt.Println(*localMap[def.IPs[e]])
		fmt.Println()

	}
}

func GetMap() def.ElevMap {
	return localMap
}
