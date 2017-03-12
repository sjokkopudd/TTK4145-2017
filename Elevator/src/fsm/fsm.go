package fsm

import (
	"def"
	"elevatorMap"
	"fmt"
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

	UP    = 1
	STILL = 0
	DOWN  = -1
)

var currentDir int
var state int
var watchdog *time.Timer

func InitFsm(msgChan_buttonEvent chan def.ChannelMessage, msgChan_fromHardware_floors chan def.ChannelMessage, msgChan_fromFsm chan def.ChannelMessage, msgChan_deadElevator chan def.ChannelMessage) {

	timer := time.NewTimer(DOOR_TIMEOUT * time.Second)
	timer.Stop()
	watchdog = time.NewTimer(IDLE_TIMEOUT * time.Second)

	for {

		select {
		case floorMsg := <-msgChan_fromHardware_floors:
			fmt.Println("Floorevent")
			switch floorMsg.Event.(def.NewEvent).EventType {

			case def.FLOOR_ARRIVAL:
				onFloorArrival(floorMsg.Event.(def.NewEvent).Data.(int), msgChan_fromFsm, timer)
				watchdog.Reset(IDLE_TIMEOUT * time.Second)
			}

		case buttonMsg := <-msgChan_buttonEvent:
			fmt.Println("Buttonevent")
			switch buttonMsg.Event.(def.NewEvent).EventType {

			case def.BUTTON_PUSH:
				button := buttonMsg.Event.(def.NewEvent).Data.([]int)
				onRequestButtonPressed(button[0], button[1], msgChan_fromFsm, timer)
				watchdog.Reset(IDLE_TIMEOUT * time.Second)
			}

		case deadElevatorMsg := <-msgChan_deadElevator:
			fmt.Println("Deadevent")
			switch deadElevatorMsg.Event.(def.NewEvent).EventType {
			case def.ELEVATOR_DEAD:
				deadElev := deadElevatorMsg.Event.(def.NewEvent).Data.(int)
				onDeadElevator(deadElev, msgChan_fromFsm)
				watchdog.Reset(IDLE_TIMEOUT * time.Second)

			}

		case <-timer.C:
			fmt.Println("Door")
			onDoorTimeout(msgChan_fromFsm)
			watchdog.Reset(IDLE_TIMEOUT * time.Second)

		case <-watchdog.C:
			fmt.Println("Deadevent")
			forceOrder(msgChan_fromFsm)
			watchdog.Reset(IDLE_TIMEOUT * time.Second)

		default:
			//fmt.Println("STATE: ", state)
		}

	}
}

func onDeadElevator(deadElev int, msgChan_fromFsm chan def.ChannelMessage) {

	m := elevatorMap.GetMap()
	m[deadElev].IsAlive = 0

	switch state {
	case IDLE:
		dir := chooseDirection(m)
		hardware.SetMotorDir(dir)
		m[def.MY_ID].Dir = dir
		if dir != def.STILL {
			state = MOVING
		}
	}

	msg := def.ConstructChannelMessage(m, nil)
	msgChan_fromFsm <- msg
}

func forceOrder(msgChan_fromFsm chan def.ChannelMessage) {

	hardware.GoToNearestFloor()

	localMap := elevatorMap.GetMap()

	switch state {
	case IDLE:
		currentDir = forceChooseDirection(localMap)
		hardware.SetMotorDir(currentDir)

		localMap[def.MY_ID].Dir = currentDir

		if currentDir != def.STILL {
			state = MOVING
		}

		msg := def.ConstructChannelMessage(localMap, nil)
		msgChan_fromFsm <- msg

	}
}

func onRequestButtonPressed(f int, b int, msgChan_fromFsm chan def.ChannelMessage, timer *time.Timer) {
	localMap := elevatorMap.GetMap()
	switch state {
	case IDLE:

		if localMap[def.MY_ID].Pos == f {

			localMap[def.MY_ID].Door = f

			localMap = clearRequests(localMap, localMap[def.MY_ID].Pos)

			hardware.SetDoorLight(1)
			msg := def.ConstructChannelMessage(localMap, nil)
			msgChan_fromFsm <- msg
			timer.Reset(DOOR_TIMEOUT * time.Second)
			state = DOOR_OPEN
		} else {
			if localMap[def.MY_ID].Buttons[f][b] != 1 {

				localMap[def.MY_ID].Buttons[f][b] = 1

				currentDir = chooseDirection(localMap)
				hardware.SetMotorDir(currentDir)

				localMap[def.MY_ID].Dir = currentDir

				if currentDir != def.STILL {
					state = MOVING
				}

				msg := def.ConstructChannelMessage(localMap, nil)
				msgChan_fromFsm <- msg
			}
		}

	case MOVING:
		if localMap[def.MY_ID].Buttons[f][b] != 1 {
			localMap[def.MY_ID].Buttons[f][b] = 1
			msg := def.ConstructChannelMessage(localMap, nil)
			msgChan_fromFsm <- msg
		}

	case DOOR_OPEN:

		if localMap[def.MY_ID].Pos == f {

			localMap[def.MY_ID].Door = f

			localMap = clearRequests(localMap, localMap[def.MY_ID].Pos)

			msg := def.ConstructChannelMessage(localMap, nil)

			msgChan_fromFsm <- msg

			timer.Reset(DOOR_TIMEOUT * time.Second)
		} else {
			if localMap[def.MY_ID].Buttons[f][b] != 1 {
				localMap[def.MY_ID].Buttons[f][b] = 1
				msg := def.ConstructChannelMessage(localMap, nil)
				msgChan_fromFsm <- msg
			}
		}
	}
}

func onFloorArrival(f int, msgChan_fromFsm chan def.ChannelMessage, timer *time.Timer) {

	localMap := elevatorMap.GetMap()
	localMap[def.MY_ID].Pos = f

	fmt.Println("New floor: ", f)

	switch state {
	case MOVING:
		if shouldStop(localMap) {
			hardware.SetMotorDir(0)

			localMap[def.MY_ID].Door = localMap[def.MY_ID].Pos

			localMap = clearRequests(localMap, localMap[def.MY_ID].Pos)

			hardware.SetDoorLight(1)

			timer.Reset(DOOR_TIMEOUT * time.Second)

			state = DOOR_OPEN

			msg := def.ConstructChannelMessage(localMap, nil)
			msgChan_fromFsm <- msg
		} else {
			msg := def.ConstructChannelMessage(localMap, nil)
			msgChan_fromFsm <- msg
		}
	case IDLE:
		msg := def.ConstructChannelMessage(localMap, nil)
		msgChan_fromFsm <- msg
	}
}

func onDoorTimeout(msgChan_fromFsm chan def.ChannelMessage) {

	switch state {
	case DOOR_OPEN:

		localMap := elevatorMap.GetMap()

		localMap[def.MY_ID].Door = -1
		hardware.SetDoorLight(0)

		currentDir = chooseDirection(localMap)
		hardware.SetMotorDir(currentDir)
		localMap[def.MY_ID].Dir = currentDir

		if currentDir == STILL {
			state = IDLE
		} else {
			state = MOVING
		}

		msg := def.ConstructChannelMessage(localMap, nil)
		msgChan_fromFsm <- msg
	}
}

func forceChooseDirection(m def.ElevMap) int {

	switch currentDir {
	case UP:
		for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Pos != 3 {
					return UP
				}
			}
		}
		for f := m[def.MY_ID].Pos - 1; f > -1; f-- {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Pos != 0 {
					return DOWN
				}
			}
		}
		return STILL

	case DOWN:
		for f := m[def.MY_ID].Pos - 1; f > -1; f-- {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Pos != 0 {
					return DOWN
				}
			}
		}
		for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Pos != 3 {
					return UP
				}
			}
		}
		return STILL

	case STILL:
		for f := m[def.MY_ID].Pos - 1; f > -1; f-- {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Pos != 0 {
					return DOWN
				}
			}
		}
		for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Pos != 3 {
					return UP
				}
			}
		}
		return STILL

	default:
		return STILL

	}

}

func chooseDirection(m def.ElevMap) int {

	switch currentDir {
	case UP:
		for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					if m[def.MY_ID].Pos != 3 {
						return UP
					}
				}
			}
		}
		for f := m[def.MY_ID].Pos - 1; f > -1; f-- {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					if m[def.MY_ID].Pos != 0 {
						return DOWN
					}
				}
			}
		}
		return STILL

	case DOWN:
		for f := m[def.MY_ID].Pos - 1; f > -1; f-- {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					if m[def.MY_ID].Pos != 0 {
						return DOWN
					}
				}
			}
		}
		for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					if m[def.MY_ID].Pos != 3 {
						return UP
					}
				}
			}
		}
		return STILL

	case STILL:
		for f := m[def.MY_ID].Pos - 1; f > -1; f-- {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					if m[def.MY_ID].Pos != 0 {
						return DOWN
					}
				}
			}
		}
		for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					if m[def.MY_ID].Pos != 3 {
						return UP
					}
				}
			}
		}
		return STILL

	default:
		return STILL

	}

}

func orderOnFloor(m def.ElevMap, f int) bool {
	return m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.UP_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.DOWN_BUTTON] == 1
}

func validOrderOnFloor(m def.ElevMap, f int) bool {

	if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 {
		return true
	}
	for e := 0; e < def.ELEVATORS; e++ {
		if m[e].Buttons[f][def.UP_BUTTON] != 1 && m[e].Buttons[f][def.DOWN_BUTTON] != 1 && m[e].IsAlive == 1 {
			return false
		}
	}
	return true
}

func iAmClosest(m def.ElevMap, f int) bool {
	result := true
	myDistance := int(math.Abs(float64(m[def.MY_ID].Pos - f)))

	// Ordren er over oss
	if m[def.MY_ID].Pos < f {

		for e := 0; e < def.ELEVATORS; e++ {

			if e != def.MY_ID && m[e].IsAlive == 1 {

				eDistance := int(math.Abs(float64(m[e].Pos - f)))

				if eDistance < myDistance { // Om en annen heis er nærmere order

					if m[e].Pos < f && (m[e].Dir == UP || m[e].Dir == IDLE) { // Om denne heisen er under order og på vei opp eller idle
						result = false
					} else if m[e].Pos > f && (m[e].Dir == DOWN || m[e].Dir == IDLE) { // Om denne heisen er over order og på vei ned eller idle
						result = false
					} else if m[e].Pos == f && m[e].Dir == IDLE { // Om denne heisen er på samme floor som order
						result = false

					}

				} else if eDistance == myDistance && (m[e].Dir == UP || m[e].Dir == IDLE) { // Om en annen heis er like nærme order og skal opp eller idle
					if m[e].ID < m[def.MY_ID].ID { // Den med lavest ID tar ordren
						result = false
					}
				}
			}
		}
		// Ordren er under oss
	} else if m[def.MY_ID].Pos > f {
		for e := 0; e < def.ELEVATORS; e++ {

			if e != def.MY_ID && m[e].IsAlive == 1 {

				eDistance := int(math.Abs(float64(m[e].Pos - f)))

				if eDistance < myDistance { // Om en annen heis er nærmere order

					if m[e].Pos < f && (m[e].Dir == UP || m[e].Dir == IDLE) { // Om denne heisen er under order og på vei opp eller idle
						result = false
					} else if m[e].Pos > f && (m[e].Dir == DOWN || m[e].Dir == IDLE) { // Om denne heisen er over order og på vei ned eller idle
						result = false
					} else if m[e].Pos == f && m[e].Dir == IDLE { // Om denne heisen er på samme floor som order
						result = false

					}

				} else if eDistance == myDistance && (m[e].Dir == DOWN || m[e].Dir == IDLE) { // Om en annen heis er like nærme order og skal ned eller idle
					if m[e].ID < m[def.MY_ID].ID { // Den med lavest ID tar ordren
						result = false
					}
				}

			}
		}
	}
	return result
}

func shouldStop(m def.ElevMap) bool {

	f := m[def.MY_ID].Pos

	switch currentDir {

	case UP:
		if m[def.MY_ID].Buttons[f][def.UP_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 {
			return true
		} else if !isOrderAbove(m) && m[def.MY_ID].Buttons[f][def.DOWN_BUTTON] == 1 {
			return true
		} else if !isOrderAbove(m) {
			return true
		} else if f == 3 {
			return true
		}

	case DOWN:
		if m[def.MY_ID].Buttons[f][def.DOWN_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 {
			return true
		} else if !isOrderBelow(m) && m[def.MY_ID].Buttons[f][def.UP_BUTTON] == 1 {
			return true
		} else if !isOrderBelow(m) {
			return true
		} else if f == 0 {
			return true
		}
	}

	return false

}

func isOrderAbove(m def.ElevMap) bool {
	for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
		if orderOnFloor(m, f) {
			return true
		}
	}
	return false
}

func isOrderBelow(m def.ElevMap) bool {
	for f := 0; f < m[def.MY_ID].Pos; f++ {
		if orderOnFloor(m, f) {
			return true
		}
	}
	return false
}

func clearRequests(m def.ElevMap, f int) def.ElevMap {
	for e := 0; e < def.ELEVATORS; e++ {
		m[e].Buttons[f][def.UP_BUTTON] = 0
		m[e].Buttons[f][def.DOWN_BUTTON] = 0
	}
	m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] = 0

	return m
}
