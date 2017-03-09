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
	localMap := GetMap()
	var fsmEvent def.NewEvent
	allAgree := true
	changeMade := false
	for e := 0; e < def.ELEVATORS; e++ {
		for f := 0; f < def.FLOORS; f++ {
			for b := 0; b < def.BUTTONS; b++ {
				if receivedMap[e].Buttons[f][b] != localMap[e].Buttons[f][b] {
					if receivedMap[e].Buttons[f][b] == 1 && localMap[e].Buttons[f][b] != 1 {
						if b != def.PANEL_BUTTON {
							localMap[e].Buttons[f][b] = 1
							localMap[def.MY_ID].Buttons[f][b] = 1
							changeMade = true
							fsmEvent = def.NewEvent{def.BUTTON_PUSH, []int{f, b}}

						} else if e == def.MY_ID{
							localMap[e].Buttons[f][b] = 1
							changeMade = true
						}else{
							localMap[e].Buttons[f][b] = 1
						}

					}
				}
				if receivedMap[e].Door == def.DOOR_OPEN && receivedMap[e].Pos == f {
					//localMap[e].Door = def.DOOR_OPEN
					localMap[def.MY_ID].Buttons[f][def.UP_BUTTON] = 0
					localMap[def.MY_ID].Buttons[f][def.DOWN_BUTTON] = 0
					receivedMap[def.MY_ID].Buttons[f][def.UP_BUTTON] = 0
					receivedMap[def.MY_ID].Buttons[f][def.DOWN_BUTTON] = 0

					localMap[e].Buttons[f][def.UP_BUTTON] = 0
					localMap[e].Buttons[f][def.DOWN_BUTTON] = 0
					receivedMap[e].Buttons[f][def.UP_BUTTON] = 0
					receivedMap[e].Buttons[f][def.DOWN_BUTTON] = 0

					localMap[e].Buttons[f][def.PANEL_BUTTON] = 0
					receivedMap[e].Buttons[f][def.PANEL_BUTTON] = 0

				}

			}
		}

		//  ?
		if receivedMap[e].Dir != localMap[e].Dir {
			localMap[e].Dir = receivedMap[e].Dir
			changeMade = true
		}
		if receivedMap[e].Pos != localMap[e].Pos && e != def.MY_ID{
			localMap[e].Pos = receivedMap[e].Pos
			changeMade = true
		}
	}
	if receivedMap[def.MY_ID].Door != localMap[def.MY_ID].Door {
			localMap[def.MY_ID].Door = receivedMap[def.MY_ID].Door
			changeMade = true
	}


	setMap(localMap)
	WriteBackup(localMap)

	return fsmEvent, localMap, changeMade, allAgree

}

func AddNewEvent(newEvent def.NewEvent) (def.ElevMap, bool, bool) {
	changeMade := false
	allAgree := true
	localMap := GetMap()
	switch newEvent.EventType {

	case def.FLOOR_ARRIVAL:
		if localMap[def.MY_ID].Pos != newEvent.Data.(int) {
			localMap[def.MY_ID].Pos = newEvent.Data.(int)
			changeMade = true
		}

	case def.BUTTON_PUSH:
		data := newEvent.Data.([]int)
		if localMap[def.MY_ID].Buttons[data[0]][data[1]] != 1 {
			localMap[def.MY_ID].Buttons[data[0]][data[1]] = 1
			changeMade = true
		}
	}
	setMap(localMap)


	if changeMade{
		WriteBackup(localMap)
	}

	return localMap, changeMade, allAgree

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
