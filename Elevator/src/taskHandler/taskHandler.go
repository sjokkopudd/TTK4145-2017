package taskHandler

import (
	"def"
	"elevatorMap"
)

func TaskHandler(eventChan chan def.NewHardwareEvent) {
	for {
		select{
		case newEvent <- eventChan:
			currentMap := elevatorMap.GetMap()
			switch newEvent.Type {
			case NEWFLOOR:
				go onFloorArrival()

			}
		}
	}
}
