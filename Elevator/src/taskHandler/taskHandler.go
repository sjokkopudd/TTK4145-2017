package taskHandler

import (
	"def"
	"elevatorMap"
	"hardware"
	"time"
)

func TaskHandler(eventChan_toTH chan def.NewHardwareEvent, doorOpenChan chan int) {
	for {
		select{
		case newEvent <- eventChan_toTH:
			currentMap := elevatorMap.GetMap()
			switch newEvent.Type {
			case def.NEWFLOOR:
				if onFloorArrival(currentMap, eventChan_toTH){
					doorOpenChan <- def.DOOR_OPEN
				}
			case def.DOOR:
				if newEvent.Door == def.DOOR_OPEN{
					hardware.SetDoorLight(1)
					time.Sleep(1*time.Second)
					doorOpenChan <- def.DOOR_CLOSE

				}else if newEvent.DOOR == def.DOOR_CLOSE{
					hardware.SetDoorLight(0)
					onDoorTimeout(currentMap)
					//decide which way to go
					//set direction
				}
			case def.BUTTONPUSH:

			}

		}
	}
}
