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

func Fsm(inDataChan chan def.ChannelMessage, outDataChan chan def.ChannelMessage) {
	for {
		select {
		case data := <-inDataChan:

			switch state {

			case IDLE:

				switch data.Event.(def.NewEvent).EventType {
				case def.BUTTONPUSH:

					doorOpen = stopAndOpenDoors(data.Map)

					if doorOpen {

						msg := def.ConstructChannelMessage(nil, def.NewEvent{def.DOOR_EVENT, 1})

						outDataChan <- msg

						state = DOOR_OPEN

						break

					}

					startedMoving, dir = takeOrder(data.Map)

					if startedMoving {

						msg := def.ConstructChannelMessage(nil, def.NewEvent{def.NEWDIR_EVENT, dir})

						outDataChan <- msg

						state = MOVING

						break
					}
				}

			case MOVING:

				switch data.Event.(def.NewEvent).EventType {
				case def.NEWFLOOR_EVENT:

					doorOpen = stopAndOpenDoors(data.Map)

					if doorOpen {

						msg := def.ConstructChannelMessage(nil, def.NewEvent{def.DOOR_EVENT, 1})

						outDataChan <- msg

						state = DOOR_OPEN

						break

					}

				}

			case DOOR_OPEN:

				switch data.Event.(def.NewEvent).EventType {
				case def.BUTTONPUSH_EVENT:

					doorOpen = stopAndOpenDoors(data.Map)

					if doorOpen {

						break

					}
				}

				closeDoors = doorTimeout()

				if closeDoors {
					msg := def.ConstructChannelMessage(nil, def.NewEvent{def.DOOR_EVENT, 0})

					outDataChan <- msg

					startedMoving, dir := takeOrder(data.Map)
					if startedMoving {
						msg := def.ConstructChannelMessage(nil, def.NewEvent{def.NEWDIR_EVENT, dir})

						outDataChan <- msg

						state = MOVING

						break
					}

					msg := def.ConstructChannelMessage(nil, def.NewEvent{def.NEWDIR_EVENT, dir})

					outDataChan <- msg

					state = IDLE
				}
			}
		}
	}
}