package elevatorMap

import (
	"def"
	"fmt"
	"sync"
)

const (
	FSM     = 0
	NETWORK = 1
)

var mapMutex = &sync.Mutex{}
var localMap *def.ElevMap

func InitMap() {

	localMap = new(def.ElevMap)
	localMap = def.NewCleanElevMap()

	WriteBackup(*localMap)

}

func AddNewMapChanges(receivedMap def.ElevMap, user int) (def.NewEvent, def.ElevMap, bool, bool) {
	currentMap := GetMap()
	var fsmEvent def.NewEvent
	allAgree := true
	changeMade := false

	for e := 0; e < def.ELEVATORS; e++ {
		if receivedMap[e].Door != currentMap[e].Door {
			if user == FSM {
				currentMap[e].Door = receivedMap[e].Door
				changeMade = true
			} else if e != def.MY_ID {
				currentMap[e].Door = receivedMap[e].Door
				changeMade = true
			} else {
				receivedMap[e].Door = currentMap[e].Door
			}

		}
		for f := 0; f < def.FLOORS; f++ {
			for b := 0; b < def.BUTTONS; b++ {

				if receivedMap[e].Buttons[f][b] == 1 && currentMap[e].Buttons[f][b] != 1 {
					if b != def.PANEL_BUTTON {
						currentMap[e].Buttons[f][b] = receivedMap[e].Buttons[f][b]
						currentMap[def.MY_ID].Buttons[f][b] = receivedMap[e].Buttons[f][b]
						changeMade = true
						fsmEvent = def.NewEvent{def.BUTTON_PUSH, []int{f, b}}
					}
				} else if receivedMap[e].Buttons[f][b] == 0 && receivedMap[e].Door == f {
					if b != def.PANEL_BUTTON {
						currentMap[e].Buttons[f][b] = receivedMap[e].Buttons[f][b]
					} else if e == def.MY_ID {
						currentMap[e].Buttons[f][def.PANEL_BUTTON] = receivedMap[e].Buttons[f][def.PANEL_BUTTON]
					}
				}
			}
		}

		//  ?
		if receivedMap[e].Dir != currentMap[e].Dir {
			currentMap[e].Dir = receivedMap[e].Dir
			changeMade = true
		}
		if receivedMap[e].Pos != currentMap[e].Pos && e != def.MY_ID {
			currentMap[e].Pos = receivedMap[e].Pos
			changeMade = true
		}

	}

	setMap(currentMap)
	WriteBackup(currentMap)

	allAgree = allButtonsAgree(currentMap)

	return fsmEvent, currentMap, changeMade, allAgree

}

func deleteOrders(receivedMap def.ElevMap, currentMap def.ElevMap) (def.ElevMap, def.ElevMap) {
	for e := 0; e < def.ELEVATORS; e++ {
		if receivedMap[e].Door != -1 {
			f := receivedMap[e].Door
			//localMap[e].Door = def.DOOR_OPEN
			currentMap[def.MY_ID].Buttons[f][def.UP_BUTTON] = 0
			currentMap[def.MY_ID].Buttons[f][def.DOWN_BUTTON] = 0
			receivedMap[def.MY_ID].Buttons[f][def.UP_BUTTON] = 0
			receivedMap[def.MY_ID].Buttons[f][def.DOWN_BUTTON] = 0

			currentMap[e].Buttons[f][def.UP_BUTTON] = 0
			currentMap[e].Buttons[f][def.DOWN_BUTTON] = 0
			receivedMap[e].Buttons[f][def.UP_BUTTON] = 0
			receivedMap[e].Buttons[f][def.DOWN_BUTTON] = 0

			currentMap[e].Buttons[f][def.PANEL_BUTTON] = 0
			receivedMap[e].Buttons[f][def.PANEL_BUTTON] = 0

			if e == def.MY_ID {
				currentMap[def.MY_ID].Buttons[f][def.PANEL_BUTTON] = 0
				receivedMap[def.MY_ID].Buttons[f][def.PANEL_BUTTON] = 0
			}
		}
	}

	return currentMap, receivedMap

}

func allButtonsAgree(m def.ElevMap) bool {

	for f := 0; f < def.FLOORS; f++ {
		for b := 0; b < def.BUTTONS-1; b++ {
			prevVal := -1
			for e := 0; e < def.ELEVATORS; e++ {
				if prevVal == -1 {
					prevVal = m[e].Buttons[f][b]
				} else if m[e].Buttons[f][b] != prevVal {
					return false
				}
			}
		}
	}
	return true
}

func AddNewEvent(newEvent def.NewEvent) (def.ElevMap, bool, bool) {
	changeMade := false
	allAgree := true
	currentMap := GetMap()
	switch newEvent.EventType {

	case def.FLOOR_ARRIVAL:
		if currentMap[def.MY_ID].Pos != newEvent.Data.(int) {
			currentMap[def.MY_ID].Pos = newEvent.Data.(int)
			changeMade = true
		}

	case def.BUTTON_PUSH:
		data := newEvent.Data.([]int)
		if currentMap[def.MY_ID].Buttons[data[0]][data[1]] != 1 {
			currentMap[def.MY_ID].Buttons[data[0]][data[1]] = 1
			changeMade = true
		}
	}
	setMap(currentMap)

	if changeMade {
		WriteBackup(currentMap)
	}

	return currentMap, changeMade, allAgree

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

	case def.FLOOR_ARRIVAL:
		fmt.Println("Event: elevator arrival at floor: ", event.Data)

	case def.BUTTON_PUSH:
		fmt.Println("Event: button pressed: ", event.Data)

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
