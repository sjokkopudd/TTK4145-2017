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
var direction int

func Fsm(inDataChan chan def.ChannelMessage, outDataChan chan def.ChannelMessage) {
	for {
		select {
		case data := <-inDataChan:

			switch state {

			case IDLE:

				switch data.Event.(def.NewEvent).EventType {
				case def.BUTTONPUSH:

					var doorOpen bool

					doorOpen, direction = stopAndOpenDoors(data.Map)

					if doorOpen {

						msg := def.ConstructChannelMessage(nil, def.NewEvent{def.DOOR_EVENT, 1})

						outDataChan <- msg

						state = DOOR_OPEN

						break

					}

					var startedMoving bool

					startedMoving, direction = takeOrder(data.Map)

					if startedMoving {

						msg := def.ConstructChannelMessage(nil, def.NewEvent{def.NEWDIR_EVENT, direction})

						outDataChan <- msg

						state = MOVING

						break
					}
				}

			case MOVING:

				switch data.Event.(def.NewEvent).EventType {
				case def.NEWFLOOR_EVENT:

					var doorOpen bool

					doorOpen, direction = stopAndOpenDoors(data.Map)

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

					var doorOpen bool
					doorOpen, direction = stopAndOpenDoors(data.Map)

					if doorOpen {

						break

					}
				}

				closeDoors = doorTimeout()

				if closeDoors {
					msg := def.ConstructChannelMessage(nil, def.NewEvent{def.DOOR_EVENT, 0})

					outDataChan <- msg

					var startedMoving bool

					startedMoving, direction = takeOrder(data.Map)
					if startedMoving {
						msg = def.ConstructChannelMessage(nil, def.NewEvent{def.NEWDIR_EVENT, direction})

						outDataChan <- msg

						state = MOVING

						break
					}

					msg = def.ConstructChannelMessage(nil, def.NewEvent{def.NEWDIR_EVENT, direction})

					outDataChan <- msg

					state = IDLE
				}
			}
		}
	}
}