package elevatorMap

import (
	"def"
	"fmt"
	"sync"
)

var localMap def.ElevMap
var mapMutex = &sync.Mutex{}

func InitMap() {

	localMap = def.NewCleanElevMap()

	WriteBackup(localMap)

}

func ReceivedMapFromNetwork(receivedMap def.ElevMap) def.NewEvent {
	newMap := GetMap()
	for e := 0; e < def.ELEVATORS; e++ {
		if e != def.MY_ID {
			for f := 0; f < def.FLOORS; f++ {
				for b := 0; b < def.BUTTONS; b++ {
					if receivedMap[e].Buttons[f][b] != newMap[e].Buttons[f][b] {
						if receivedMap[e].Buttons[f][b] == 1 && newMap[e].Buttons[f][b] != 1 {
							if b != def.PANEL_BUTTON {
								newMap[def.MY_ID].Buttons[f][b] = 1
								changes := def.NewEvent{def.BUTTONPUSH, []int{f, b}}
								setMap(newMap)
								return changes

							} else {
								newMap[e].Buttons[f][b] = 1
								changes := def.NewEvent{def.OTHERELEVATOR, -1}
								setMap(newMap)
								return changes
							}

						}
					} else {
						if receivedMap[e].Door == 1 && receivedMap[e].Pos == f {
							newMap[e].Door = 1
							newMap[def.MY_ID].Buttons[f][def.UP_BUTTON] = 0
							newMap[def.MY_ID].Buttons[f][def.DOWN_BUTTON] = 0
							changes := def.NewEvent{def.OTHERELEVATOR, -1}
							setMap(newMap)
							return changes
						}
					}

				}
			}

			if receivedMap[e].Dir != newMap[e].Dir {
				newMap[e].Dir = receivedMap[e].Dir
				changes := def.NewEvent{def.OTHERELEVATOR, -1}
				setMap(newMap)
				return changes
			}
			if receivedMap[e].Pos != newMap[e].Pos {
				newMap[e].Pos = receivedMap[e].Pos
				changes := def.NewEvent{def.OTHERELEVATOR, -1}
				setMap(newMap)
				return changes
			}
			if receivedMap[e].Door != newMap[e].Door {
				newMap[e].Door = receivedMap[e].Door
				changes := def.NewEvent{def.OTHERELEVATOR, -1}
				setMap(newMap)
				return changes
			}
		}
	}

	changes := def.NewEvent{def.NEWFLOOR, newMap[def.MY_ID].Pos}
	return changes

}

func UpdateMap(newEvent def.NewEvent) (def.ElevMap, bool) {
	changeMade := false
	newMap := GetMap()
	switch newEvent.EventType {

	case def.NEWFLOOR:
		if newMap[def.MY_ID].Pos != newEvent.Data.(int) {
			newMap[def.MY_ID].Pos = newEvent.Data.(int)
			changeMade = true
		}

	case def.BUTTONPUSH:
		data := newEvent.Data.([]int)
		if newMap[def.MY_ID].Buttons[data[0]][data[1]] != 1 {
			newMap[def.MY_ID].Buttons[data[0]][data[1]] = 1
			changeMade = true
		}

	case def.DOOR:
		newMap[def.MY_ID].Door = newEvent.Data.(int)
		newMap[def.MY_ID].Dir = def.IDLE
		for b := 0; b < def.BUTTONS; b++ {
			newMap[def.MY_ID].Buttons[newMap[def.MY_ID].Pos][b] = 0
		}
		changeMade = true

	case def.OTHERELEVATOR:
		changeMade = true

	case def.NEWDIR:
		if newMap[def.MY_ID].Dir != newEvent.Data.(int) {
			newMap[def.MY_ID].Dir = newEvent.Data.(int)
			changeMade = true
		}
	}
	if changeMade {
		WriteBackup(newMap)
		setMap(newMap)
	}

	return newMap, changeMade

}

func PrintMap(elevatorMap def.ElevMap) {
	for e := 0; e < def.ELEVATORS; e++ {
		fmt.Println("ID: ", e)
		fmt.Println(elevatorMap[e])
		fmt.Println()

	}
}

func GetMap() def.ElevMap {
	mapMutex.Lock()
	currentMap := localMap
	mapMutex.Unlock()
	return currentMap
}

func setMap(newMap def.ElevMap) {
	mapMutex.Lock()
	localMap = newMap
	mapMutex.Unlock()
}
