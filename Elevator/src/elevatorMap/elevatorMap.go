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

func AddNewMapChanges(receivedMap def.ElevMap) (def.NewEvent, def.ElevMap, bool, bool) {
	currentMap := GetMap()
	var fsmEvent def.NewEvent
	allAgree := true
	changeMade := false
	for e := 0; e < def.ELEVATORS; e++ {
		for f := 0; f < def.FLOORS; f++ {
			for b := 0; b < def.BUTTONS; b++ {
				if receivedMap[e].Buttons[f][b] != currentMap[e].Buttons[f][b] {
					if receivedMap[e].Buttons[f][b] == 1 && currentMap[e].Buttons[f][b] != 1 {
						if b != def.PANEL_BUTTON {
							currentMap[e].Buttons[f][b] = 1
							currentMap[def.MY_ID].Buttons[f][b] = 1
							changeMade = true
							fsmEvent = def.NewEvent{def.BUTTON_PUSH, []int{f, b}}

						} else if e == def.MY_ID {
							currentMap[e].Buttons[f][b] = 1
							changeMade = true
						} else {
							currentMap[e].Buttons[f][b] = 1
						}

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

	currentMap = deleteOrders(receivedMap, currentMap)

	if receivedMap[def.MY_ID].Door != currentMap[def.MY_ID].Door {
		currentMap[def.MY_ID].Door = receivedMap[def.MY_ID].Door
		changeMade = true
	}

	setMap(currentMap)
	WriteBackup(currentMap)

	allAgree = allButtonsAgree(currentMap)

	return fsmEvent, currentMap, changeMade, allAgree

}

func deleteOrders(receivedMap def.ElevMap, localMap def.ElevMap) def.ElevMap{
	for e:= 0; e < def.ELEVATORS; e ++{
		for f:= 0; f < def.ELEVATORS; f++{
			if receivedMap[e].Door == def.DOOR_OPEN && receivedMap[e].Pos == f {
				//localMap[e].Door = def.DOOR_OPEN
				localMap[def.MY_ID].Buttons[f][def.UP_BUTTON] = 0
				localMap[def.MY_ID].Buttons[f][def.DOWN_BUTTON] = 0

				localMap[e].Buttons[f][def.UP_BUTTON] = 0
				localMap[e].Buttons[f][def.DOWN_BUTTON] = 0

				localMap[e].Buttons[f][def.PANEL_BUTTON] = 0

			}
		}
	}

	return localMap

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
