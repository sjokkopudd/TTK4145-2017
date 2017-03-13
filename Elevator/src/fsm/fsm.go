package fsm

import (
	"def"
	"elevatorMap"
	"hardware"
	"math"
	"time"
)

const (
	IDLE      = 0
	MOVING    = 1
	DOOR_OPEN = 2

	DOOR_TIMEOUT = 2
	IDLE_TIMEOUT = 10
)

var currentDirection int
var state int

func Fsm(msgChan_buttonEvent chan def.ChannelMessage, msgChan_fromHardware_floors chan def.ChannelMessage, msgChan_fromFsm chan def.ChannelMessage, msgChan_deadElevator chan def.ChannelMessage) {

	doorTimer := time.NewTimer(DOOR_TIMEOUT * time.Second)
	doorTimer.Stop()
	idleTimeoutTimer := time.NewTimer(IDLE_TIMEOUT * time.Second)

	state = IDLE
	currentDirection = def.STILL

	for {

		select {

		case floorMsg := <-msgChan_fromHardware_floors:

			switch floorMsg.Event.(def.NewEvent).EventType {

			case def.FLOOR_ARRIVAL:

				onFloorArrival(floorMsg.Event.(def.NewEvent).Data.(int), msgChan_fromFsm, doorTimer)
				idleTimeoutTimer.Reset(IDLE_TIMEOUT * time.Second)

			}

		case buttonMsg := <-msgChan_buttonEvent:

			switch buttonMsg.Event.(def.NewEvent).EventType {

			case def.BUTTON_PUSH:

				button := buttonMsg.Event.(def.NewEvent).Data.([]int)
				onRequestButtonPressed(button[0], button[1], msgChan_fromFsm, doorTimer)
				idleTimeoutTimer.Reset(IDLE_TIMEOUT * time.Second)

			}

		case deadElevatorMsg := <-msgChan_deadElevator:

			switch deadElevatorMsg.Event.(def.NewEvent).EventType {

			case def.ELEVATOR_DEAD:

				deadElev := deadElevatorMsg.Event.(def.NewEvent).Data.(int)
				onDeadElevator(deadElev, msgChan_fromFsm)
				idleTimeoutTimer.Reset(IDLE_TIMEOUT * time.Second)

			}

		case <-doorTimer.C:

			onDoorTimeout(msgChan_fromFsm)
			idleTimeoutTimer.Reset(IDLE_TIMEOUT * time.Second)

		case <-idleTimeoutTimer.C:

			forceOrder(msgChan_fromFsm)
			idleTimeoutTimer.Reset(IDLE_TIMEOUT * time.Second)

		}

	}
}

func onDeadElevator(deadElev int, msgChan_fromFsm chan def.ChannelMessage) {

	currentMap := elevatorMap.GetLocalMap()
	currentMap[deadElev].IsAlive = 0

	switch state {

	case IDLE:

		currentDirection = chooseDirection(currentMap)
		hardware.SetMotorDir(currentDirection)
		currentMap[def.MY_ID].Direction = currentDirection

		if currentDirection != def.STILL {
			state = MOVING
		}

	}

	msg := def.ConstructChannelMessage(currentMap, nil)
	msgChan_fromFsm <- msg
}

func forceOrder(msgChan_fromFsm chan def.ChannelMessage) {

	hardware.GoToNearestFloor()

	currentMap := elevatorMap.GetLocalMap()

	switch state {

	case IDLE:

		currentDirection = forceChooseDirection(currentMap)
		hardware.SetMotorDir(currentDirection)
		currentMap[def.MY_ID].Direction = currentDirection

		if currentDirection != def.STILL {
			state = MOVING
		}

		msg := def.ConstructChannelMessage(currentMap, nil)
		msgChan_fromFsm <- msg

	}
}

func onRequestButtonPressed(floor int, button int, msgChan_fromFsm chan def.ChannelMessage, doorTimer *time.Timer) {
	currentMap := elevatorMap.GetLocalMap()

	switch state {

	case IDLE:

		if currentMap[def.MY_ID].Position == floor {

			currentMap[def.MY_ID].Door = floor
			currentMap = clearRequests(currentMap, currentMap[def.MY_ID].Position)
			hardware.SetDoorLight(1)

			msg := def.ConstructChannelMessage(currentMap, nil)
			msgChan_fromFsm <- msg

			doorTimer.Reset(DOOR_TIMEOUT * time.Second)
			state = DOOR_OPEN

		} else {

			currentMap[def.MY_ID].Buttons[floor][button] = 1

			currentDirection = chooseDirection(currentMap)
			hardware.SetMotorDir(currentDirection)
			currentMap[def.MY_ID].Direction = currentDirection

			if currentDirection != def.STILL {
				state = MOVING
			}

			msg := def.ConstructChannelMessage(currentMap, nil)
			msgChan_fromFsm <- msg

		}

	case MOVING:

		if currentMap[def.MY_ID].Buttons[floor][button] != 1 {
			currentMap[def.MY_ID].Buttons[floor][button] = 1

			msg := def.ConstructChannelMessage(currentMap, nil)
			msgChan_fromFsm <- msg
		}

	case DOOR_OPEN:

		if currentMap[def.MY_ID].Position == floor {

			currentMap[def.MY_ID].Door = floor
			currentMap = clearRequests(currentMap, currentMap[def.MY_ID].Position)

			msg := def.ConstructChannelMessage(currentMap, nil)
			msgChan_fromFsm <- msg

			doorTimer.Reset(DOOR_TIMEOUT * time.Second)
		} else {
			if currentMap[def.MY_ID].Buttons[floor][button] != 1 {
				currentMap[def.MY_ID].Buttons[floor][button] = 1

				msg := def.ConstructChannelMessage(currentMap, nil)
				msgChan_fromFsm <- msg
			}
		}
	}
}

func onFloorArrival(floor int, msgChan_fromFsm chan def.ChannelMessage, doorTimer *time.Timer) {

	currentMap := elevatorMap.GetLocalMap()
	currentMap[def.MY_ID].Position = floor

	switch state {
	case MOVING:

		if shouldStop(currentMap) {
			hardware.SetMotorDir(0)

			currentMap[def.MY_ID].Door = currentMap[def.MY_ID].Position
			currentMap = clearRequests(currentMap, currentMap[def.MY_ID].Position)
			hardware.SetDoorLight(1)
			doorTimer.Reset(DOOR_TIMEOUT * time.Second)
			state = DOOR_OPEN

			msg := def.ConstructChannelMessage(currentMap, nil)
			msgChan_fromFsm <- msg
		} else {
			msg := def.ConstructChannelMessage(currentMap, nil)
			msgChan_fromFsm <- msg
		}

	case IDLE:

		msg := def.ConstructChannelMessage(currentMap, nil)
		msgChan_fromFsm <- msg
	}
}

func onDoorTimeout(msgChan_fromFsm chan def.ChannelMessage) {

	switch state {

	case DOOR_OPEN:

		currentMap := elevatorMap.GetLocalMap()

		currentMap[def.MY_ID].Door = -1
		hardware.SetDoorLight(0)

		currentDirection = chooseDirection(currentMap)
		hardware.SetMotorDir(currentDirection)
		currentMap[def.MY_ID].Direction = currentDirection

		if currentDirection == def.STILL {
			state = IDLE
		} else {
			state = MOVING
		}

		msg := def.ConstructChannelMessage(currentMap, nil)
		msgChan_fromFsm <- msg
	}
}

func forceChooseDirection(currentMap elevatorMap.ElevMap) int {

	switch currentDirection {

	case def.UP:

		for f := currentMap[def.MY_ID].Position + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Position != 3 {
					return def.UP
				}
			}
		}
		for f := currentMap[def.MY_ID].Position - 1; f > -1; f-- {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Position != 0 {
					return def.DOWN
				}
			}
		}
		return def.STILL

	case def.DOWN:

		for f := currentMap[def.MY_ID].Position - 1; f > -1; f-- {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Position != 0 {
					return def.DOWN
				}
			}
		}
		for f := currentMap[def.MY_ID].Position + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Position != 3 {
					return def.UP
				}
			}
		}
		return def.STILL

	case def.STILL:

		for f := currentMap[def.MY_ID].Position - 1; f > -1; f-- {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Position != 0 {
					return def.DOWN
				}
			}
		}
		for f := currentMap[def.MY_ID].Position + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Position != 3 {
					return def.UP
				}
			}
		}
		return def.STILL

	default:

		return def.STILL
	}

}

func chooseDirection(currentMap elevatorMap.ElevMap) int {

	switch currentDirection {
	case def.UP:
		for f := currentMap[def.MY_ID].Position + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || isClosestElevator(currentMap, f) {
					if currentMap[def.MY_ID].Position != 3 {
						return def.UP
					}
				}
			}
		}
		for f := currentMap[def.MY_ID].Position - 1; f > -1; f-- {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || isClosestElevator(currentMap, f) {
					if currentMap[def.MY_ID].Position != 0 {
						return def.DOWN
					}
				}
			}
		}
		return def.STILL

	case def.DOWN:
		for f := currentMap[def.MY_ID].Position - 1; f > -1; f-- {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || isClosestElevator(currentMap, f) {
					if currentMap[def.MY_ID].Position != 0 {
						return def.DOWN
					}
				}
			}
		}
		for f := currentMap[def.MY_ID].Position + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || isClosestElevator(currentMap, f) {
					if currentMap[def.MY_ID].Position != 3 {
						return def.UP
					}
				}
			}
		}
		return def.STILL

	case def.STILL:
		for f := currentMap[def.MY_ID].Position - 1; f > -1; f-- {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || isClosestElevator(currentMap, f) {
					if currentMap[def.MY_ID].Position != 0 {
						return def.DOWN
					}
				}
			}
		}
		for f := currentMap[def.MY_ID].Position + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(currentMap, f) {
				if currentMap[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || isClosestElevator(currentMap, f) {
					if currentMap[def.MY_ID].Position != 3 {
						return def.UP
					}
				}
			}
		}
		return def.STILL

	default:
		return def.STILL

	}

}

func validOrderOnFloor(currentMap elevatorMap.ElevMap, floor int) bool {

	if currentMap[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] == 1 {
		return true
	}
	for e := 0; e < def.ELEVATORS; e++ {
		if currentMap[e].Buttons[floor][def.UP_BUTTON] != 1 && currentMap[e].Buttons[floor][def.DOWN_BUTTON] != 1 && currentMap[e].IsAlive == 1 {
			return false
		}
	}
	return true
}

func isClosestElevator(currentMap elevatorMap.ElevMap, floor int) bool {
	result := true
	myDistance := int(math.Abs(float64(currentMap[def.MY_ID].Position - floor)))

	if currentMap[def.MY_ID].Position < floor {

		for e := 0; e < def.ELEVATORS; e++ {

			if e != def.MY_ID && currentMap[e].IsAlive == 1 {

				eDistance := int(math.Abs(float64(currentMap[e].Position - floor)))

				if eDistance < myDistance {

					if currentMap[e].Position < floor && (currentMap[e].Direction == def.UP || currentMap[e].Direction == IDLE) {
						result = false
					} else if currentMap[e].Position > floor && (currentMap[e].Direction == def.DOWN || currentMap[e].Direction == IDLE) {
						result = false
					} else if currentMap[e].Position == floor && currentMap[e].Direction == IDLE {
						result = false
					}

				} else if eDistance == myDistance && (currentMap[e].Direction == def.UP || currentMap[e].Direction == IDLE) {
					if e < def.MY_ID {
						result = false
					}
				}
			}
		}
	} else if currentMap[def.MY_ID].Position > floor {
		for e := 0; e < def.ELEVATORS; e++ {

			if e != def.MY_ID && currentMap[e].IsAlive == 1 {

				eDistance := int(math.Abs(float64(currentMap[e].Position - floor)))

				if eDistance < myDistance {

					if currentMap[e].Position < floor && (currentMap[e].Direction == def.UP || currentMap[e].Direction == IDLE) {
						result = false
					} else if currentMap[e].Position > floor && (currentMap[e].Direction == def.DOWN || currentMap[e].Direction == IDLE) {
						result = false
					} else if currentMap[e].Position == floor && currentMap[e].Direction == IDLE {
						result = false

					}

				} else if eDistance == myDistance && (currentMap[e].Direction == def.DOWN || currentMap[e].Direction == IDLE) {
					if currentMap[e].ID < currentMap[def.MY_ID].ID {
						result = false
					}
				}

			}
		}
	}
	return result
}

func shouldStop(currentMap elevatorMap.ElevMap) bool {

	pos := currentMap[def.MY_ID].Position

	switch currentDirection {

	case def.UP:
		if currentMap[def.MY_ID].Buttons[pos][def.UP_BUTTON] == 1 || currentMap[def.MY_ID].Buttons[pos][def.PANEL_BUTTON] == 1 {
			return true
		} else if !isOrderAbove(currentMap) && currentMap[def.MY_ID].Buttons[pos][def.DOWN_BUTTON] == 1 {
			return true
		} else if !isOrderAbove(currentMap) {
			return true
		} else if pos == 3 {
			return true
		}

	case def.DOWN:
		if currentMap[def.MY_ID].Buttons[pos][def.DOWN_BUTTON] == 1 || currentMap[def.MY_ID].Buttons[pos][def.PANEL_BUTTON] == 1 {
			return true
		} else if !isOrderBelow(currentMap) && currentMap[def.MY_ID].Buttons[pos][def.UP_BUTTON] == 1 {
			return true
		} else if !isOrderBelow(currentMap) {
			return true
		} else if pos == 0 {
			return true
		}
	}

	return false

}

func isOrderAbove(currentMap elevatorMap.ElevMap) bool {
	for f := currentMap[def.MY_ID].Position + 1; f < def.FLOORS; f++ {
		if validOrderOnFloor(currentMap, f) {
			return true
		}
	}
	return false
}

func isOrderBelow(currentMap elevatorMap.ElevMap) bool {
	for f := 0; f < currentMap[def.MY_ID].Position; f++ {
		if validOrderOnFloor(currentMap, f) {
			return true
		}
	}
	return false
}

func clearRequests(currentMap elevatorMap.ElevMap, floor int) elevatorMap.ElevMap {
	for e := 0; e < def.ELEVATORS; e++ {
		currentMap[e].Buttons[floor][def.UP_BUTTON] = 0
		currentMap[e].Buttons[floor][def.DOWN_BUTTON] = 0
	}
	currentMap[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] = 0

	return currentMap
}
