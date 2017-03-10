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
	IDLE_TIMEOUT = 5

	UP    = 1
	STILL = 0
	DOWN  = -1
)

var currentDir int
var state int

func InitFsm(inDataChan chan def.ChannelMessage, outDataChan chan def.ChannelMessage) {

	timer := time.NewTimer(DOOR_TIMEOUT * time.Second)
	timer.Stop()
	idleTimer := time.NewTimer(IDLE_TIMEOUT * time.Second)
	idleTimer.Stop()

	for {

		select {
		case data := <-inDataChan:

			switch data.Event.(def.NewEvent).EventType {
			case def.BUTTON_PUSH:
				button := data.Event.(def.NewEvent).Data.([]int)
				onRequestButtonPressed(button[0], button[1], outDataChan, timer)
				idleTimer.Reset(IDLE_TIMEOUT * time.Second)

			case def.FLOOR_ARRIVAL:
				onFloorArrival(data.Event.(def.NewEvent).Data.(int), outDataChan, timer)
				idleTimer.Reset(IDLE_TIMEOUT * time.Second)

			}

		case <-timer.C:

			onDoorTimeout(outDataChan, idleTimer)
			idleTimer.Reset(IDLE_TIMEOUT * time.Second)

		case <-idleTimer.C:
			forceOrder(outDataChan, idleTimer)

		default:
			//fmt.Println("STATE: ", state)
		}

	}
}

func forceOrder(outDataChan chan def.ChannelMessage, idleTimer *time.Timer) {

	localMap := elevatorMap.GetMap()

	if isOrderAbove(localMap) || isOrderBelow(localMap) {
		dir := chooseDirection(localMap)
		if dir == def.STILL {
			if isOrderAbove(localMap) {
				dir = def.UP
			} else {
				dir = def.DOWN
			}
		}
		hardware.SetMotorDir(dir)
		localMap[def.MY_ID].Dir = dir
		msg := def.ConstructChannelMessage(localMap, nil)
		outDataChan <- msg

		state = MOVING

	} else {
		idleTimer.Reset(IDLE_TIMEOUT * time.Second)
		state = IDLE
	}
}

func onRequestButtonPressed(f int, b int, outDataChan chan def.ChannelMessage, timer *time.Timer) {
	localMap := elevatorMap.GetMap()
	switch state {
	case IDLE:

		if localMap[def.MY_ID].Pos == f {

			localMap[def.MY_ID].Door = f

			localMap = clearRequests(localMap, localMap[def.MY_ID].Pos)

			hardware.SetDoorLight(1)
			msg := def.ConstructChannelMessage(localMap, nil)
			outDataChan <- msg
			timer.Reset(DOOR_TIMEOUT * time.Second)
			state = DOOR_OPEN
		} else {
			localMap[def.MY_ID].Buttons[f][b] = 1

			currentDir = chooseDirection(localMap)
			hardware.SetMotorDir(currentDir)

			localMap[def.MY_ID].Dir = currentDir

			if currentDir != def.STILL {
				state = MOVING
			}

			msg := def.ConstructChannelMessage(localMap, nil)
			outDataChan <- msg

		}

	case MOVING:
		localMap[def.MY_ID].Buttons[f][b] = 1
		msg := def.ConstructChannelMessage(localMap, nil)
		outDataChan <- msg

	case DOOR_OPEN:
		localMap := elevatorMap.GetMap()

		if localMap[def.MY_ID].Pos == f {

			localMap[def.MY_ID].Door = f

			localMap = clearRequests(localMap, localMap[def.MY_ID].Pos)

			msg := def.ConstructChannelMessage(localMap, nil)

			outDataChan <- msg

			timer.Reset(DOOR_TIMEOUT * time.Second)
		} else {
			localMap[def.MY_ID].Buttons[f][b] = 1
			msg := def.ConstructChannelMessage(localMap, nil)
			outDataChan <- msg
		}
	}
}

func onFloorArrival(f int, outDataChan chan def.ChannelMessage, timer *time.Timer) {

	localMap := elevatorMap.GetMap()
	localMap[def.MY_ID].Pos = f

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
			outDataChan <- msg
		} else {
			msg := def.ConstructChannelMessage(localMap, nil)
			outDataChan <- msg
		}
	case IDLE:
		msg := def.ConstructChannelMessage(localMap, nil)
		outDataChan <- msg
	}
}

func onDoorTimeout(outDataChan chan def.ChannelMessage, idleTimer *time.Timer) {

	switch state {
	case DOOR_OPEN:

		localMap := elevatorMap.GetMap()

		localMap[def.MY_ID].Door = -1
		hardware.SetDoorLight(0)

		currentDir = chooseDirection(localMap)
		hardware.SetMotorDir(currentDir)
		localMap[def.MY_ID].Dir = currentDir

		if currentDir == STILL {
			idleTimer.Reset(IDLE_TIMEOUT * time.Second)
			state = IDLE
		} else {
			state = MOVING
		}

		msg := def.ConstructChannelMessage(localMap, nil)
		outDataChan <- msg
	}
}

func chooseDirection(m def.ElevMap) int {

	switch currentDir {
	case UP:
		for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return UP
				}
			}
		}
		for f := m[def.MY_ID].Pos - 1; f > -1; f-- {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return DOWN
				}
			}
		}
		return STILL

	case DOWN:
		for f := m[def.MY_ID].Pos - 1; f > -1; f-- {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return DOWN
				}
			}
		}
		for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return UP
				}
			}
		}
		return STILL

	case STILL:
		for f := m[def.MY_ID].Pos - 1; f > -1; f-- {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return DOWN
				}
			}
		}
		for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
			if validOrderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return UP
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
		if m[e].Buttons[f][def.UP_BUTTON] != 1 && m[e].Buttons[f][def.DOWN_BUTTON] != 1 {
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

			if e != def.MY_ID {

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

			if e != def.MY_ID {

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
		}

	case DOWN:
		if m[def.MY_ID].Buttons[f][def.DOWN_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 {
			return true
		} else if !isOrderBelow(m) && m[def.MY_ID].Buttons[f][def.UP_BUTTON] == 1 {
			return true
		} else if !isOrderBelow(m) {
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
