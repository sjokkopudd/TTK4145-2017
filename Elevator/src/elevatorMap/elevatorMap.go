package elevatorMap

import (
	"def"
	"sync"
)

var mapMutex = &sync.Mutex{}
var localMap *def.ElevMap

func InitMap(backup bool) {
	mapMutex.Lock()
	localMap = new(def.ElevMap)
	if backup {
		*localMap = readBackup()
	} else {
		localMap = def.NewCleanElevMap()
	}

	writeBackup(*localMap)
	mapMutex.Unlock()

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
		if currentMap[e].IsAlive == 1 && receivedMap[e].IsAlive != 1 {
			currentMap[e].IsAlive = 0
		}
		if currentMap[e].IsAlive == 1 {
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

	}

	setMap(currentMap)
	writeBackup(currentMap)
	return currentMap, changeMade
}

func GetEventFromNetwork(receivedMap def.ElevMap) (def.NewEvent, def.ElevMap) {
	currentMap := GetMap()
	var fsmEvent def.NewEvent
	floorWithDoorOpen := -1

	for e := 0; e < def.ELEVATORS; e++ {
		if receivedMap[e].Door != -1 {
			floorWithDoorOpen = receivedMap[e].Door
		}
	}

	for e := 0; e < def.ELEVATORS; e++ {
		if receivedMap[e].IsAlive == 1 {
			for f := 0; f < def.FLOORS; f++ {
				for b := 0; b < def.BUTTONS; b++ {

					if receivedMap[e].Buttons[f][b] == 1 && currentMap[e].Buttons[f][b] != 1 {
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
			if currentMap[e].IsAlive != 1 {
				currentMap[e].IsAlive = 1
			}
		}
	}

	setMap(currentMap)
	writeBackup(currentMap)
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

	if changeMade {
		setMap(currentMap)
		writeBackup(currentMap)
	}

	return currentMap, changeMade

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
