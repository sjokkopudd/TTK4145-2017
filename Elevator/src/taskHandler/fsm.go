package taskHandler

import (
	"def"

)

func onFloorArrival(currentMap def.ElevMap, newEvent def.NewHardwareEvent) bool {
	for b := 0; b < def.BUTTONS; b++{
		if currentMap[def.MY_IP].Buttons[newEvent.Pos][i] == 1{
			return true
		}
	}
	return false
}

func onDoorTimeout(currentMap def.ElevMap) {

}

func onNewButtonEvent() {

}
