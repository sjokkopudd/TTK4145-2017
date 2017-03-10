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

func AddNewMapChanges(receivedMap def.ElevMap, user int) (def.ElevMap, bool) {
	currentMap := GetMap()
	changeMade := false
	floorWithDoorOpen := -1

	if receivedMap[def.MY_ID].Door != -1 {
		floorWithDoorOpen = receivedMap[def.MY_ID].Door
		currentMap[def.MY_ID].Door = receivedMap[def.MY_ID].Door
		changeMade = true
	}

	if receivedMap[def.MY_ID].Dir != currentMap[def.MY_ID].Dir {
		currentMap[def.MY_ID].Dir = receivedMap[def.MY_ID].Dir
		changeMade = true
	}

	if receivedMap[def.MY_ID].Pos != currentMap[def.MY_ID].Pos {
		currentMap[def.MY_ID].Pos = receivedMap[def.MY_ID].Pos
		changeMade = true
	}

	for e := 0; e < def.ELEVATORS; e++ {
		for f := 0; f < def.FLOORS; f++ {
			for b := 0; b < def.BUTTONS; b++ {

				if receivedMap[e].Buttons[f][b] == 1 && currentMap[e].Buttons[f][b] != 1 {
					if b != def.PANEL_BUTTON {
						currentMap[e].Buttons[f][b] = receivedMap[e].Buttons[f][b]
						currentMap[def.MY_ID].Buttons[f][b] = receivedMap[e].Buttons[f][b]
						changeMade = true
					} else if e == def.MY_ID {
						currentMap[e].Buttons[f][def.PANEL_BUTTON] = receivedMap[e].Buttons[f][def.PANEL_BUTTON]
						changeMade = true
					}
				} else if receivedMap[e].Buttons[f][b] == 0 && floorWithDoorOpen == f {
					if b != def.PANEL_BUTTON {
						currentMap[e].Buttons[f][b] = receivedMap[e].Buttons[f][b]
						changeMade = true
					} else if e == def.MY_ID {
						currentMap[e].Buttons[f][def.PANEL_BUTTON] = receivedMap[e].Buttons[f][def.PANEL_BUTTON]
						changeMade = true
					}
				}
			}
		}
	}

	setMap(currentMap)
	WriteBackup(currentMap)
	return currentMap, changeMade
}

func GetEventFromNetwork(receivedMap def.ElevMap) (def.NewEvent, def.ElevMap) {
	currentMap := GetMap()
	var fsmEvent def.NewEvent
	floorWithDoorOpen := -1

	for e := 0; e < def.ELEVATORS; e++ {
		if receivedMap[e].Door != -1 {
			floorWithDoorOpen = receivedMap[e].Door
			if e == def.MY_ID {
				currentMap[e].Door = receivedMap[e].Door
			}
		}
	}

	for e := 0; e < def.ELEVATORS; e++ {
		for f := 0; f < def.FLOORS; f++ {
			for b := 0; b < def.BUTTONS; b++ {

				if receivedMap[e].Buttons[f][b] == 1 && currentMap[e].Buttons[f][b] != 1 /*&& f != currentMap[def.MY_ID].Door */ {
					if b != def.PANEL_BUTTON {
						currentMap[e].Buttons[f][b] = receivedMap[e].Buttons[f][b]
						fsmEvent = def.NewEvent{def.BUTTON_PUSH, []int{f, b}}
					} else {
						currentMap[e].Buttons[f][b] = receivedMap[e].Buttons[f][b]
					}
				} else if receivedMap[e].Buttons[f][b] == 0 && floorWithDoorOpen == f {
					currentMap[e].Buttons[f][b] = receivedMap[e].Buttons[f][b]
				}
			}
		}

		if receivedMap[e].Dir != currentMap[e].Dir && e != def.MY_ID {
			currentMap[e].Dir = receivedMap[e].Dir
		}
		if receivedMap[e].Pos != currentMap[e].Pos && e != def.MY_ID {
			currentMap[e].Pos = receivedMap[e].Pos
		}
	}

	setMap(currentMap)
	WriteBackup(currentMap)
	return fsmEvent, currentMap
}

func AddNewEvent(newEvent def.NewEvent) (def.ElevMap, bool) {
	changeMade := false
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

	return currentMap, changeMade

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
