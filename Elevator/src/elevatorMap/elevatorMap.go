package elevatorMap

import (
	"def"
	"fmt"
	"sync"
)

var mapMutex = &sync.Mutex{}
var localMap *ElevMap

func InitMap(backup bool) {
	mapMutex.Lock()
	localMap = new(ElevMap)
	if backup {
		*localMap = readBackup()
	} else {
		localMap = NewCleanElevMap()
	}

	writeBackup(*localMap)
	mapMutex.Unlock()

}

func AddNewMapChanges(receivedMap ElevMap, user int) (ElevMap, bool) {
	currentMap := GetLocalMap()
	changeMade := false
	floorWithDoorOpen := -1

	if receivedMap[def.MY_ID].Door != currentMap[def.MY_ID].Door {
		if receivedMap[def.MY_ID].Door != -1 {
			floorWithDoorOpen = receivedMap[def.MY_ID].Door
		}
		currentMap[def.MY_ID].Door = receivedMap[def.MY_ID].Door
		changeMade = true
	}

	if receivedMap[def.MY_ID].Direction != currentMap[def.MY_ID].Direction {
		currentMap[def.MY_ID].Direction = receivedMap[def.MY_ID].Direction
		changeMade = true
	}

	if receivedMap[def.MY_ID].Position != currentMap[def.MY_ID].Position {
		currentMap[def.MY_ID].Position = receivedMap[def.MY_ID].Position
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
							fmt.Println("hello")
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

	setLocalMap(currentMap)
	writeBackup(currentMap)
	return currentMap, changeMade
}

func GetEventFromNetwork(receivedMap ElevMap) (def.NewEvent, ElevMap) {
	currentMap := GetLocalMap()
	var fsmEvent def.NewEvent
	floorWithDoorOpen := -1

	for e := 0; e < def.ELEVATORS; e++ {
		if receivedMap[e].Door != -1 {
			floorWithDoorOpen = receivedMap[e].Door
		}
	}

	for e := 0; e < def.ELEVATORS; e++ {
		if currentMap[e].IsAlive != 1 {
			currentMap[e].IsAlive = 1
		}
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

			if receivedMap[e].Direction != currentMap[e].Direction && e != def.MY_ID {
				currentMap[e].Direction = receivedMap[e].Direction
			}
			if receivedMap[e].Position != currentMap[e].Position && e != def.MY_ID {
				currentMap[e].Position = receivedMap[e].Position
			}

		}
	}

	setLocalMap(currentMap)
	writeBackup(currentMap)
	return fsmEvent, currentMap
}

func GetLocalMap() ElevMap {
	mapMutex.Lock()
	currentMap := *localMap
	mapMutex.Unlock()
	return currentMap
}

type ElevMap [def.ELEVATORS]def.ElevatorInfo

func NewCleanElevMap() *ElevMap {

	newMap := new(ElevMap)

	for e := 0; e < def.ELEVATORS; e++ {
		newMap[e].ID = e
		for f := 0; f < def.FLOORS; f++ {
			for b := 0; b < def.BUTTONS; b++ {
				newMap[e].Buttons[f][b] = 0
			}
		}
		newMap[e].Direction = def.STILL
		newMap[e].Position = 0
		newMap[e].Door = -1
		newMap[e].IsAlive = 1
	}
	return newMap
}

func setLocalMap(newMap ElevMap) {
	mapMutex.Lock()
	*localMap = newMap
	mapMutex.Unlock()
}
