package taskHandler

import (
	"def"
)

func onFloorArrival(currentMap def.ElevMap, newEvent def.NewEvent) bool {
	if currentMap[def.MY_ID].Pos == 3 || currentMap[def.MY_ID].Pos == 0 {
		return true
	}
	for b := 0; b < def.BUTTONS; b++ {
		if currentMap[def.MY_ID].Buttons[newEvent.Data.(int)][b] == 1 {
			if shouldStop(b, currentMap) {
				return true
			}
		}
	}
	return false
}

func shouldStop(button int, currentMap def.ElevMap) bool {
	if button == def.PANEL_BUTTON {
		return true
	} else if button == def.UP_BUTTON {
		if currentMap[def.MY_ID].Dir == def.UP {
			return true
		} else if currentMap[def.MY_ID].Dir == def.DOWN {
			if isOrderBelow(currentMap) {
				return false
			}
		} else {
			return true
		}
	} else if button == def.DOWN_BUTTON {
		if currentMap[def.MY_ID].Dir == def.DOWN {
			return true
		} else if currentMap[def.MY_ID].Dir == def.UP {
			if isOrderAbove(currentMap) {
				return false
			}
		} else {
			return true
		}
	}
	return true
}

func isOrderAbove(currentMap def.ElevMap) bool {
	for f := currentMap[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
		for b := 0; b < def.BUTTONS; b++ {
			if currentMap[def.MY_ID].Buttons[f][b] == 1 {
				return true
			}
		}
	}
	return false
}

func isOrderBelow(currentMap def.ElevMap) bool {
	for f := 0; f < currentMap[def.MY_ID].Pos; f++ {
		for b := 0; b < def.BUTTONS; b++ {
			if currentMap[def.MY_ID].Buttons[f][b] == 1 {
				return true
			}
		}
	}
	return false
}

func onDoorTimeout(currentMap def.ElevMap) {

}

func onNewButtonEvent() {

}

func chooseDirection(currentMap def.ElevMap) int {
	for e := 0; e < def.MY_ID; e++ {
		if currentMap[e].Dir == def.IDLE {
			for f := 0; f < def.FLOORS; f++ {
				if currentMap[e].Buttons[f][def.PANEL_BUTTON] == 0 {
					return def.IDLE
				}
			}
		}
	}

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
	if currentMap[def.MY_ID].Pos-floor < 0 {
		return def.UP
	} else if currentMap[def.MY_ID].Pos-floor > 0 {
		return def.DOWN
	}
	return def.IDLE
}
