package taskHandler

import (
	"def"

)

func onFloorArrival(currentMap def.ElevMap, newEvent def.NewHardwareEvent) {
	for b := 0; b < def.BUTTONS; b++{
		if currentMap[def.MY_IP].Buttons[newEvent.Pos][i] == 1{
			
		}
	}
}

func onDoorTimeout() {

}

func onNewButtonEvent() {

}

func onNetworkRecvMsg() {

}
