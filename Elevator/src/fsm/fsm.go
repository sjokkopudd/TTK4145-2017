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

const (
	UP    = 1
	STILL = 0
	DOWN  = -1
)

var currentDir int

func Fsm(inDataChan chan def.ChannelMessage) {

	for {
		select {
		case data := <-inDataChan:

			switch data.Event.(def.NewEvent).EventType {
			case def.BUTTONPUSH_EVENT:
				onRequestButtonPressed(data.Event.(def.NewEvent).Data[0], data.Event.(def.NewEvent).Data[1])

			case def.NEWFLOOR_EVENT:
				onFloorArrival(data.Event.(def.NewEvent).Data)

			case def.TIMEOUT_EVENT:
				onDoorTimeOut()
			}
		}
	}
}

func onRequestButtonPressed(f int, b int) {

	localMap := elevatorMap.GetMap()

	switch state {
	case IDLE:
		if localMap[def.MY_ID].Pos == f {
			hardware.SetDoorLigth(1)
			timer.Start(1)
		} else {
			dir = chooseDirection(localMap)
			hardware.SetMotorDirection(dir)
			state = MOVING
		}

	case MOVING:

	case DOOR_OPEN:
		if localMap[def.MY_ID].Pos == f {
			timer.Start(1)
		}
	}
}

func onFloorArrival(f int) {

	localMap := elevator.GetMap()

	switch state {
	case MOVING:
		if shouldStop(localMap) {
			hardware.SetMotorDirection(0)
			hardware.SetDoorLigth(1)
			timer.Start(1)

			elevatorMap.ClearRequests(localMap)

			state = DOOR_OPEN
		}
	}
}

func onDoorTimeout() {
	localMap := elevator.GetMap()

	switch state {
	case DOOR_OPEN:
		hardware.SetDoorLigth(0)

		dir = chooseDirection(localMap)
		hardware.SetMotorDirection(dir)

		if dir == 0 {
			state = IDLE
		} else {
			state = MOVING
		}
	}
}

func chooseDirection(m def.ElevMap) int {

	switch currentDir {
	case UP:
		for f := m[def.MY_ID].Pos; f < def.FLOORS; f++ {
			if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.UP_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.DOWN_BUTTON] == 1 {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return UP
				}
			}
		}

	case DOWN:
		for f := m[def.MY_ID].Pos; f > -1; f-- {
			if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.UP_BUTTON] == 1 || m[def.MY_ID].Buttons[f][def.DOWN_BUTTON] == 1 {
				if m[def.MY_ID].Buttons[f][def.PANEL_BUTTON] == 1 || iAmClosest(m, f) {
					return DOWN
				}
			}
		}

	default:
		return STILL
	}

}

func iAmClosest(m def.ElevMap, f int) bool {
	result := true
	for e := 0; e < def.ELEVATORS; e++ {
		if e != def.MY_ID {

		}
	}
}
