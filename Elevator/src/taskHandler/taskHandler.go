package taskHandler

import (
	"def"
	"elevatorMap"
	"hardware"
	"time"
)

const (
	IDLE      = 0
	MOVING    = 1
	DOOR_OPEN = 2
)

var state int

func Fsm(dataChan chan def.Data) {
	select {
	case data := <-dataChan:

		switch state {

		case IDLE:
			// button press
			// button on same floor? yes: state = DOOR_OPEN; break
			// should take? state = MOVING; break

		case MOVING:
			// floor arrival
			// should stop?
			// yes: state = DOOR_OPEN
			// no: ignore

		case DOOR_OPEN:
			// button press on same floor
			// reset door timer

			// door timeout
			// more in same dir? yes: state = MOVING same dir
			// more in other dir?yes: state = MOVING other dir
			// no: state = IDLE

		}
	}
}

func EventHandler(eventChan_toTH chan def.NewEvent, evenChan_fromtTH chan def.NewEvent) {
	for {
		select {
		case newEvent := <-eventChan_toTH:
			currentMap := elevatorMap.GetMap()
			switch newEvent.EventType {
			case def.NEWFLOOR:
				if onFloorArrival(currentMap, newEvent) {
					evenChan_fromtTH <- def.NewEvent{def.DOOR, def.DOOR_OPEN}
				}
			case def.DOOR:
				hardware.SetMotorDir(def.IDLE)
				if newEvent.Data.(int) == def.DOOR_OPEN {
					time.Sleep(1 * time.Second)
					//onDoorTimeout(currentMap, newEvent)
					evenChan_fromtTH <- def.NewEvent{def.DOOR, def.DOOR_CLOSE}

				} else if newEvent.Data.(int) == def.DOOR_CLOSE {
					dir := chooseDirection(currentMap, newEvent)
					evenChan_fromtTH <- def.NewEvent{def.NEWDIR, dir}
				}
			case def.BUTTONPUSH:
				data := newEvent.Data.([]int)
				if data[0] == currentMap[def.MY_ID].Pos {
					evenChan_fromtTH <- def.NewEvent{def.DOOR, def.DOOR_OPEN}
				} else if currentMap[def.MY_ID].Dir == def.IDLE {
					dir := chooseDirection(currentMap, newEvent)
					if dir != def.IDLE {
						evenChan_fromtTH <- def.NewEvent{def.NEWDIR, dir}
					}
				}
			case def.NEWDIR:
				hardware.SetMotorDir(newEvent.Data.(int))
			case def.OTHERELEVATOR:

			case def.ELEVATORDEAD:

			}
		}
	}

}
