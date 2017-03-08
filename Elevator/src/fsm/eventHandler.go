package fsm

import (
	"def"
	"hardware"
	"time"
)

var doorTimer float32

func stopAndOpenDoors(m def.ElevMap) bool, int {
	floor := m[def.MY_ID].Pos

	switch m[def.MY_ID].dir {

	case def.UP:
		if m[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] == 1  {
			doorTimer = time.Now().Second()
			hardware.SetMotorDir(def.STILL)
			return true, m[def.MY_ID].dir
		}else if  m[def.MY_ID].Buttons[floor][def.UP_BUTTON] == 1{
			doorTimer = time.Now().Second()
			hardware.SetMotorDir(def.STILL)
			return true, def.UP
		}else if !isOrderAbove(m) {
			doorTimer = time.Now().Second() - 1
			hardware.SetMotorDir(def.STILL)
			return true, def.DOWN
		}
		return false, m[def.MY_ID].dir

	case def.DOWN:
		if  m[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] == 1  {
			doorTimer = time.Now().Second()
			hardware.SetMotorDir(def.STILL)
			return true, m[def.MY_ID].dir
		}else if  m[def.MY_ID].Buttons[floor][def.DOWN_BUTTON] == 1{
			doorTimer = time.Now().Second()
			hardware.SetMotorDir(def.STILL)
			return true, def.DOWN
		}else if !isOrderBelow(m) {
			doorTimer = time.Now().Second() - 1
			hardware.SetMotorDir(def.STILL)
			return true, def.UP
		}
		return false, m[def.MY_ID].dir

	case def.STILL:
		if m[def.MY_ID].Buttons[floor][def.UP_BUTTON] == 1  || m[def.MY_ID].Buttons[floor][def.DOWN_BUTTON] == 1  || m[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] == 1  {
			doorTimer = time.Now().Second()
			return true, def.STILL
		}
		return false, m[def.MY_ID].dir

}

func isOrderAbove(currentMap def.ElevMap) bool {
	for f := currentMap[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
		isOrder,_ := isOrderHere(m,f)
		if isOrder{
			return true
		}
	}
	return false
}

func isOrderBelow(m def.ElevMap) bool {
	for f := 0; f < m[def.MY_ID].Pos; f++ {
		isOrder,_ := isOrderHere(m,f)
		if isOrder{
			return true
		}
	}
	return false
}

func doorTimeout() bool {
	return (time.Now().Second()-doorTimer > 1)
}

func shouldTake(m data.Map, f int, d int) bool {
	floor := m[def.MY_ID].Pos
	isOrder, button := isOrderHere(m,f)
	if isOrder{
		if b == def.PANEL_BUTTON{
				return true
		}
		for e := 0; e < def.MY_ID; e++{
			if (m[e].Pos < f) || (m[e].Pos >= floor){
				if m[e].Dir == STILL || m[e].Dir == d{
					return false
				}
			}
				 
		}
	}
	true
}


func takeOrder(m data.Map)bool,int{

	floor := m[def.MY_ID].Pos

	switch direction{
	case def.UP:

		for f := def.FLOORS; f > floor ; f--{
			if shouldTake(m, f, def.UP){
				return true, def.UP
			}
		}

		for f := 0; f < floor; f++{
			if shouldTake(m, f, def.DOWN){
				return true, def.DOWN
			}
		}

		return false, def.STILL
	case def.DOWN:

		for f := 0; f < floor; f++{
			if shouldTake(m, f, def.DOWN){
				return true, def.DOWN
			}
		}

		for f := def.FLOORS; f > floor ; f--{
			if shouldTake(m, f, def.UP){
				return true, def.UP
			}
		}

		return false, def.STILL

	case def.STILL:
		f := findClosestOrder(m)
		if f == -1{
			return false, def.STILL
		}else if f < floor {
			if shouldTake(m, f, def.DOWN){
				return true, def.DOWN
			}
		}else if f > floor{
			if shouldTake(m, f, def.UP){
				return true, def.UP
			}
		}

		return false, def.STILL
	}
	return false, def.STILL
}

func isOrderHere(m def.ElevMap, f int) (bool, int){
	for b := 0; b < def.BUTTONS; b++ {
		if currentMap[def.MY_ID].Buttons[f][b] == 1 {
				return true, b
		}
	}
	return false, -1
}

func findClosestOrder(m def.ElevMap) (int, int){
	floor := m[def.MY_ID].Pos

	for i := 0; i < def.FLOORS; i ++{
		if floor + i < def.FLOORS{
			isOrder,_ := isOrderHere(m,floor+i)
			if isOrder{
				return floor + i
			}
		}
		if floor - i > -1{
			isOrder,_ := isOrderHere(m,floor+i)
			if isOrder{
				return floor - i
			}
		}

	}

	return -1
}