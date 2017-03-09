package fsm

/*
import (
	"def"
	"hardware"
	"time"
)

var doorTimer float32

func doorTimeout(m def.ElevMap, timeoutChan chan bool, offset float32) {
	doorTimer = float32(time.Now().Second()) - offset

	time.Sleep(1000 * time.Millisecond)

	if float32(time.Now().Second())-doorTimer >= 1 {

		msg := true
		timeoutChan <- msg
	}
}

func stopAndOpenDoors(m def.ElevMap, timeoutChan chan bool) (bool, int) {
	floor := m[def.MY_ID].Pos

	switch m[def.MY_ID].Dir {

	case def.UP:
		if m[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] == 1 {
			go doorTimeout(m, timeoutChan, 0)
			hardware.SetMotorDir(def.STILL)
			return true, m[def.MY_ID].Dir
		} else if m[def.MY_ID].Buttons[floor][def.UP_BUTTON] == 1 {
			go doorTimeout(m, timeoutChan, 0)
			hardware.SetMotorDir(def.STILL)
			return true, def.UP
		} else if !isOrderAbove(m) {
			go doorTimeout(m, timeoutChan, 1)
			hardware.SetMotorDir(def.STILL)
			return true, def.DOWN
		}
		return false, m[def.MY_ID].Dir

	case def.DOWN:
		if m[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] == 1 {
			go doorTimeout(m, timeoutChan, 0)
			hardware.SetMotorDir(def.STILL)
			return true, m[def.MY_ID].Dir
		} else if m[def.MY_ID].Buttons[floor][def.DOWN_BUTTON] == 1 {
			go doorTimeout(m, timeoutChan, 0)
			hardware.SetMotorDir(def.STILL)
			return true, def.DOWN
		} else if !isOrderBelow(m) {
			go doorTimeout(m, timeoutChan, 1)
			hardware.SetMotorDir(def.STILL)
			return true, def.UP
		}
		return false, m[def.MY_ID].Dir

	case def.STILL:
		if m[def.MY_ID].Buttons[floor][def.UP_BUTTON] == 1 || m[def.MY_ID].Buttons[floor][def.DOWN_BUTTON] == 1 || m[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] == 1 {
			go doorTimeout(m, timeoutChan, 0)
			return true, def.STILL
		}
		return false, m[def.MY_ID].Dir

	}
	return false, m[def.MY_ID].Dir
}

func isOrderAbove(m def.ElevMap) bool {
	for f := m[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
		isOrder, _ := isOrderHere(m, f)
		if isOrder {
			return true
		}
	}
	return false
}

func isOrderBelow(m def.ElevMap) bool {
	for f := 0; f < m[def.MY_ID].Pos; f++ {
		isOrder, _ := isOrderHere(m, f)
		if isOrder {
			return true
		}
	}
	return false
}

func shouldTake(m def.ElevMap, f int, d int) bool {
	floor := m[def.MY_ID].Pos
	isOrder, button := isOrderHere(m, f)
	if isOrder {
		if button == def.PANEL_BUTTON {
			return true
		}
		for e := 0; e < def.MY_ID; e++ {
			if (m[e].Pos < f) || (m[e].Pos >= floor) {
				if m[e].Dir == def.STILL || m[e].Dir == d {
					return false
				}
			}

		}
		return true
	}
	return false
}

func takeOrder(m def.ElevMap) (bool, int) {

	floor := m[def.MY_ID].Pos

	switch direction {
	case def.UP:

		for f := def.FLOORS - 1; f > floor; f-- {
			if shouldTake(m, f, def.UP) {
				hardware.SetMotorDir(def.UP)
				return true, def.UP
			}
		}

		for f := 0; f < floor; f++ {
			if shouldTake(m, f, def.DOWN) {
				hardware.SetMotorDir(def.DOWN)
				return true, def.DOWN
			}
		}

		hardware.SetMotorDir(def.STILL)
		return false, def.STILL
	case def.DOWN:

		for f := 0; f < floor; f++ {
			if shouldTake(m, f, def.DOWN) {
				hardware.SetMotorDir(def.DOWN)
				return true, def.DOWN
			}
		}

		for f := def.FLOORS - 1; f > floor; f-- {
			if shouldTake(m, f, def.UP) {
				hardware.SetMotorDir(def.UP)
				return true, def.UP
			}
		}

		hardware.SetMotorDir(def.STILL)
		return false, def.STILL

	case def.STILL:
		for r := 1; r < def.FLOORS; r++ {
			f := findClosestOrder(m, r)
			if f > -1 && (f < floor) {
				if shouldTake(m, f, def.DOWN) {
					hardware.SetMotorDir(def.DOWN)
					return true, def.DOWN
				}
			} else if f > floor {
				if shouldTake(m, f, def.UP) {
					hardware.SetMotorDir(def.UP)
					return true, def.UP
				}
			}
		}
		hardware.SetMotorDir(def.STILL)
		return false, def.STILL
	}
	hardware.SetMotorDir(def.STILL)
	return false, def.STILL
}

func isOrderHere(m def.ElevMap, f int) (bool, int) {
	for b := 0; b < def.BUTTONS; b++ {
		if m[def.MY_ID].Buttons[f][b] == 1 {
			return true, b
		}
	}
	return false, -1
}

func findClosestOrder(m def.ElevMap, r int) int {
	floor := m[def.MY_ID].Pos

	if (floor + r) < def.FLOORS {
		isOrder, _ := isOrderHere(m, floor+r)
		if isOrder {
			return floor + r
		}
	}
	if (floor - r) > -1 {
		isOrder, _ := isOrderHere(m, floor-r)
		if isOrder {
			return floor - r
		}
	}

	return -1
}*/
