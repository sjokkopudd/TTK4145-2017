package taskHandler

import (
	"def"
)

func onFloorArrival(currentMap def.ElevMap, newEvent def.NewEvent) bool {
	for b := 0; b < def.BUTTONS; b++ {
		if currentMap[def.MY_ID].Buttons[newEvent.Data.(int)][b] == 1 {
			return true
		}
	}
	return false
}

func onDoorTimeout(currentMap def.ElevMap) {

}

func onNewButtonEvent() {

}

func chooseDirection(currentMap def.ElevMap) int {
	for f := 0; f < def.FLOORS; f++ {
		if currentMap[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 {
			return isAbove(currentMap, f)
		}
	}
	for f := 0; f < def.FLOORS; f++ {
		for b := 0; b < def.BUTTONS-1; b++ {
			if currentMap[def.MY_ID].Buttons[f][b] == 1 {
				return isAbove(currentMap, f)
			}
		}
	}
	return def.IDLE
}

func isAbove(currentMap def.ElevMap, floor int) int {
	if currentMap[def.MY_ID].Pos-floor > 0 {
		return def.UP
	} else if currentMap[def.MY_ID].Pos-floor < 0 {
		return def.DOWN
	}
	return def.IDLE
}
