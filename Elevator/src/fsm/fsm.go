package fsm

import (
	"def"
	"elevatorMap"
	"fmt"
)

const (
	IDLE      = 0
	MOVING    = 1
	DOOR_OPEN = 2
)

var state int
var direction int

func Fsm(inDataChan chan def.ChannelMessage, outDataChan chan def.ChannelMessage) {

	timeoutChan := make(chan bool)

	for {
		select {
		case data := <-inDataChan:

			switch state {

			case IDLE:

				fmt.Println("State: IDLE")

				switch data.Event.(def.NewEvent).EventType {
				case def.BUTTONPUSH_EVENT:
					fmt.Println("Button pushed")

					var doorOpen bool

					doorOpen, direction = stopAndOpenDoors(data.Map.(def.ElevMap), timeoutChan)

					if doorOpen {

						fmt.Println("Door open")

						msg := def.ConstructChannelMessage(nil, def.NewEvent{def.DOOR_EVENT, def.DOOR_OPEN})

						outDataChan <- msg

						state = DOOR_OPEN

						break

					}

					var startedMoving bool

					startedMoving, direction = takeOrder(data.Map.(def.ElevMap))

					if startedMoving {

						msg := def.ConstructChannelMessage(nil, def.NewEvent{def.NEWDIR_EVENT, direction})

						outDataChan <- msg

						state = MOVING

						break
					}
				}

			case MOVING:

				fmt.Println("State: MOVING")

				switch data.Event.(def.NewEvent).EventType {
				case def.NEWFLOOR_EVENT:

					fmt.Println("Reached floor")

					var doorOpen bool

					doorOpen, direction = stopAndOpenDoors(data.Map.(def.ElevMap), timeoutChan)

					if doorOpen {

						msg := def.ConstructChannelMessage(nil, def.NewEvent{def.DOOR_EVENT, def.DOOR_OPEN})

						outDataChan <- msg

						state = DOOR_OPEN

						break

					}

				}

			case DOOR_OPEN:

				fmt.Println("State: DOOR_OPEN")

				switch data.Event.(def.NewEvent).EventType {
				case def.BUTTONPUSH_EVENT:

					var doorOpen bool
					doorOpen, direction = stopAndOpenDoors(data.Map.(def.ElevMap), timeoutChan)

					if doorOpen {
						msg := def.ConstructChannelMessage(nil, def.NewEvent{def.DOOR_EVENT, def.DOOR_OPEN})
						outDataChan <- msg
						break

					}
				}
			}
		case timeoutData := <-timeoutChan:

			if timeoutData {

				fmt.Println("Door timeout")

				msg := def.ConstructChannelMessage(nil, def.NewEvent{def.DOOR_EVENT, def.DOOR_CLOSED})

				outDataChan <- msg

				var startedMoving bool

				startedMoving, direction = takeOrder(elevatorMap.GetMap())

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
