package elevatorMap

import (
	"def"
	"fmt"
	"sync"
)

var mapMutex = &sync.Mutex{}
var localMap *def.ElevMap

func InitMap() {

	localMap = new(def.ElevMap)
	localMap = def.NewCleanElevMap()

	WriteBackup(*localMap)

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
								newMap[e].Buttons[f][b] = 1
								changes := def.NewEvent{def.BUTTONPUSH_EVENT, []int{f, b}}
								setMap(newMap)
								return changes

							} else {
								newMap[e].Buttons[f][b] = 1
								changes := def.NewEvent{def.OTHERELEVATOR_EVENT, -1}
								setMap(newMap)
								return changes
							}

						}
					} else {
						if receivedMap[e].Door == 1 && receivedMap[e].Pos == f {
							newMap[e].Door = 1
							newMap[def.MY_ID].Buttons[f][def.UP_BUTTON] = 0
							newMap[def.MY_ID].Buttons[f][def.DOWN_BUTTON] = 0

							newMap[e].Buttons[f][def.UP_BUTTON] = 0
							newMap[e].Buttons[f][def.DOWN_BUTTON] = 0

							newMap[e].Buttons[f][def.PANEL_BUTTON] = 0

							changes := def.NewEvent{def.OTHERELEVATOR_EVENT, -1}
							setMap(newMap)
							return changes
						}
					}

				}
			}

			if receivedMap[e].Dir != newMap[e].Dir {
				newMap[e].Dir = receivedMap[e].Dir
				changes := def.NewEvent{def.OTHERELEVATOR_EVENT, -1}
				setMap(newMap)
				return changes
			}
			if receivedMap[e].Pos != newMap[e].Pos {
				newMap[e].Pos = receivedMap[e].Pos
				changes := def.NewEvent{def.OTHERELEVATOR_EVENT, -1}
				setMap(newMap)
				return changes
			}
			if receivedMap[e].Door != newMap[e].Door {
				newMap[e].Door = receivedMap[e].Door
				changes := def.NewEvent{def.OTHERELEVATOR_EVENT, -1}
				setMap(newMap)
				return changes
			}
		}
	}

	changes := def.NewEvent{def.NEWFLOOR_EVENT, newMap[def.MY_ID].Pos}
	return changes

}

func UpdateMap(newEvent def.NewEvent) (def.ElevMap, bool) {
	changeMade := false
	newMap := GetMap()
	switch newEvent.EventType {

	case def.NEWFLOOR_EVENT:
		if newMap[def.MY_ID].Pos != newEvent.Data.(int) {
			newMap[def.MY_ID].Pos = newEvent.Data.(int)
			changeMade = true
		}

	case def.BUTTONPUSH_EVENT:
		data := newEvent.Data.([]int)
		if newMap[def.MY_ID].Buttons[data[0]][data[1]] != 1 {
			newMap[def.MY_ID].Buttons[data[0]][data[1]] = 1
			changeMade = true
		}

	case def.DOOR_EVENT:
		newMap[def.MY_ID].Door = newEvent.Data.(int)
		newMap[def.MY_ID].Dir = def.IDLE
		for b := 0; b < def.BUTTONS; b++ {
			newMap[def.MY_ID].Buttons[newMap[def.MY_ID].Pos][b] = 0
		}
		changeMade = true

	case def.OTHERELEVATOR_EVENT:
		changeMade = true

	case def.NEWDIR_EVENT:
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
		if e == def.MY_ID {
			fmt.Println(elevatorMap[e].ID, " - My Map ")

		} else {
			fmt.Println(elevatorMap[e].ID)
		}
		for f := 0; f < def.FLOORS; f++ {
			fmt.Println(elevatorMap[e].Buttons[f])
		}
		fmt.Println(elevatorMap[e].Dir)
		fmt.Println(elevatorMap[e].Pos)
		fmt.Println(elevatorMap[e].Door)
		fmt.Println(elevatorMap[e].IsAlive)

	}
}

func PrintEvent(event def.NewEvent) {
	switch event.EventType {

	case def.NEWFLOOR_EVENT:
		fmt.Println("Event: elevator arrival at floor: ", event.Data)

	case def.BUTTONPUSH_EVENT:
		fmt.Println("Event: button pressed: ", event.Data)

	case def.DOOR_EVENT:
		fmt.Println("Event: door open/close")

	case def.OTHERELEVATOR_EVENT:
		fmt.Println("Event: another elevator did something")

	case def.NEWDIR_EVENT:
		fmt.Println("Event: elevator changed direction")

	}

}

func GetMap() def.ElevMap {
	mapMutex.Lock()
	currentMap := *localMap
	mapMutex.Unlock()
	return currentMap
}

func setMap(newMap def.ElevMap) {
	mapMutex.Lock()
	*localMap = newMap
	mapMutex.Unlock()

}
