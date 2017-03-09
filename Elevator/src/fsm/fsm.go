package fsm

import (
	"def"
	"elevatorMap"
	"hardware"
	"math"
)

const (
	IDLE      = 0
	MOVING    = 1
	DOOR_OPEN = 2
)

const (
	UP    = 1
	STILL = 0
	DOWN  = -1
)

var currentDir int
var state int

func Fsm(inDataChan chan def.ChannelMessage, outDataChan chan def.ChannelMessage) {

	timeoutChan := make(chan bool, 1)
	go timer(timeoutChan)

	for {
		select {
		case data := <-inDataChan:

			switch data.Event.(def.NewEvent).EventType {
			case def.BUTTON_PUSH:
				button := data.Event.(def.NewEvent).Data.([]int)
				onRequestButtonPressed(button[0], button[1], outDataChan)

			case def.FLOOR_ARRIVAL:
				onFloorArrival(data.Event.(def.NewEvent).Data.(int), outDataChan)

			}

		case timeout := <-timeoutChan:

			if timeout {
				onDoorTimeout(outDataChan)
			}
		}
	}
}

func onRequestButtonPressed(f int, b int, outDataChan chan def.ChannelMessage) {

	localMap := elevatorMap.GetMap()

	switch state {
	case IDLE:
		if localMap[def.MY_ID].Pos == f {
			localMap[def.MY_ID].Door = 1
			msg := def.ConstructChannelMessage(localMap, nil)
			outDataChan <- msg
			timerStart()
		} else {
			currentDir = chooseDirection(localMap)
			hardware.SetMotorDirection(currentDircurrentDir)
			localMap[def.MY_ID].Door = currentDir
			state = MOVING

			msg := def.ConstructChannelMessage(localMap, nil)
			outDataChan <- msg

		}

	case MOVING:

	case DOOR_OPEN:
		if localMap[def.MY_ID].Pos == f {
			timer.Start(1)
		}
	}
}

func onFloorArrival(f int, outDataChan chan def.ChannelMessage) {

	localMap := elevator.GetMap()

	switch state {
	case MOVING:
		if shouldStop(localMap) {
			hardware.SetMotorDirection(0)
			localMap[def.MY_ID].Door = 1

			timer.Start(1)

			elevatorMap.ClearRequests(localMap)

			state = DOOR_OPEN

			msg := def.ConstructChannelMessage(localMap, nil)
			outDataChan <- msg
		}
	}
}

func onDoorTimeout(outDataChan chan def.ChannelMessage) {
	localMap := elevator.GetMap()

	switch state {
	case DOOR_OPEN:
		hardware.SetDoorLigth(0)
		localMap[def.MY_ID].Door = 0

		currentDir = chooseDirection(localMap)
		hardware.SetMotorDirection(currentDir)
		localMap[def.MY_ID].Door = currentDir

		if currentDir == STILL {
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
		for f := m[def.MY_ID].Pos; f < def.FLOORS; f++ {
			if orderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return UP
				}
			}
		}
		for f := m[def.MY_ID].Pos; f > -1; f-- {
			if orderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return DOWN
				}
			}
		}
		return STILL

	case DOWN:
		for f := m[def.MY_ID].Pos; f > -1; f-- {
			if orderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return DOWN
				}
			}
		}
		for f := m[def.MY_ID].Pos; f < def.FLOORS; f++ {
			if orderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return UP
				}
			}
		}
		return STILL

	case STILL:
		for f := m[def.MY_ID].Pos; f > -1; f-- {
			if orderOnFloor(m, f) {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return DOWN
				}
			}
		}
		for f := m[def.MY_ID].Pos; f < def.FLOORS; f++ {
			if orderOnFloor(m, f) {
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

func orderOnFloor(m def.ElevMap, f int) {
	return m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.UP_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.DOWN_BUTTON] == 1
}

func iAmClosest(m def.ElevMap, f int) bool {
	result := true

	// Ordren er over oss
	if m[def.MY_ID].Pos < f {

		myDistance := math.Abs(m[def.MY_ID].Pos - f)

		for e := 0; e < def.ELEVATORS; e++ {

			if e != def.MY_ID {

				eDistance := math.Abs(m[e].Pos - f)

				if eDistance < myDistance { // Om en annen heis er nærmere order

					if m[e].Pos < f && (m[e].dir == UP || m[e].dir == IDLE) { // Om denne heisen er under order og på vei opp eller idle
						result = false
					} else if m[e].Pos > f && (m[e].dir == DOWN || m[e].dir == IDLE) { // Om denne heisen er over order og på vei ned eller idle
						result = false
					} else if m[e].Pos == f { // Om denne heisen er på samme floor som order
						result = false

					}

				} else if eDistance == myDistance && (m[e].dir == UP || m[e].dir == IDLE) { // Om en annen heis er like nærme order og skal opp eller idle
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

				eDistance := math.Abs(m[e].Pos - f)

				if eDistance < myDistance { // Om en annen heis er nærmere order

					if m[e].Pos < f && (m[e].dir == UP || m[e].dir == IDLE) { // Om denne heisen er under order og på vei opp eller idle
						result = false
					} else if m[e].Pos > f && (m[e].dir == DOWN || m[e].dir == IDLE) { // Om denne heisen er over order og på vei ned eller idle
						result = false
					} else if m[e].Pos == f { // Om denne heisen er på samme floor som order
						result = false

					}

				} else if eDistance == myDistance && (m[e].dir == DOWN || m[e].dir == IDLE) { // Om en annen heis er like nærme order og skal ned eller idle
					if m[e].ID < m[def.MY_ID].ID { // Den med lavest ID tar ordren
						result = false
					}
				}

			}
		}
	}
	return result
}
