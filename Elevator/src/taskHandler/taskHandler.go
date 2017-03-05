package taskHandler

import (
	"def"
	"elevatorMap"
	"hardware"
	"time"
)

func EventHandler(eventChan_toTH chan def.NewEvent, evenChan_fromtTH chan def.NewEvent) {
	for {
		select {
		case newEvent := <-eventChan_toTH:
			currentMap := elevatorMap.GetMap()
			switch newEvent.EventType {
			case def.NEWFLOOR:
				if onFloorArrival(currentMap, newEvent) {
					evenChan_fromtTH <- def.NewEvent{def.DOOR, def.DOOR_OPEN}
				} /*else {
					dir := chooseDirection(currentMap)

					evenChan_fromtTH <- def.NewEvent{def.NEWDIR, dir}
				}*/
			case def.DOOR:
				hardware.SetMotorDir(def.IDLE)
				if newEvent.Data.(int) == def.DOOR_OPEN {
					time.Sleep(1 * time.Second)
					onDoorTimeout(currentMap)
					evenChan_fromtTH <- def.NewEvent{def.DOOR, def.DOOR_CLOSE}

				} else if newEvent.Data == def.DOOR_CLOSE {
					dir := chooseDirection(currentMap)
					evenChan_fromtTH <- def.NewEvent{def.NEWDIR, dir}
				}
			case def.BUTTONPUSH:
				if currentMap[def.MY_ID].Dir == def.IDLE {
					dir := chooseDirection(currentMap)
					evenChan_fromtTH <- def.NewEvent{def.NEWDIR, dir}
				}
			case def.NEWDIR:
				hardware.SetMotorDir(newEvent.Data.(int))

			case def.OTHERELEVATOR:

			case def.ELEVATORDEAD:

			}
		}
	}

}
