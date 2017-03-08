package fsm

import (
	"def"
	"hardware"
	"time"
)

var doorTimer float32

func stopAndOpenDoors(m def.ElevMap) bool {
	floor := m[def.MY_ID].Pos

	switch m[def.MY_ID].dir {

	case def.UP:
		if m[def.MY_ID].Buttons[floor][def.UP_BUTTON] == 1 || m[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] == 1  {
			doorTimer = time.Now().Second()
			hardware.SetMotorDir(def.STILL)
			return true
		} else if !isOrderAbove(m) {
			doorTimer = time.Now().Second() - 1
			hardware.SetMotorDir(def.STILL)
			return true
		}
		return false

	case def.DOWN:
		if m[def.MY_ID].Buttons[floor][def.DOWN_BUTTON] == 1  || m[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] == 1  {
			doorTimer = time.Now().Second()
			hardware.SetMotorDir(def.STILL)
			return true
		} else if !isOrderBelow(m) {
			doorTimer = time.Now().Second() - 1
			hardware.SetMotorDir(def.STILL)
			return true
		}
		return false

	case def.STILL:
		if m[def.MY_ID].Buttons[floor][def.UP_BUTTON] == 1  || m[def.MY_ID].Buttons[floor][def.DOWN_BUTTON] == 1  || m[def.MY_ID].Buttons[floor][def.PANEL_BUTTON] == 1  {
			doorTimer = time.Now().Second()
			return true
		}
		return false
	}

}

func isOrderAbove(currentMap def.ElevMap) bool {
	for f := currentMap[def.MY_ID].Pos + 1; f < def.FLOORS; f++ {
		for b := 0; b < def.BUTTONS; b++ {
			if currentMap[def.MY_ID].Buttons[f][b] == 1 {
				return true
			}
		}
	}
	return false
}

func isOrderBelow(currentMap def.ElevMap) bool {
	for f := 0; f < currentMap[def.MY_ID].Pos; f++ {
		for b := 0; b < def.BUTTONS; b++ {
			if currentMap[def.MY_ID].Buttons[f][b] == 1 {
				return true
			}
		}
	}
	return false
}

func doorTimeout() bool {
	return (time.Now().Second()-doorTimer > 1)
}


func takeOrder(m data.Map)bool,int{


	Idle?
		isOrderAbove(currentMap)
			if elevator above me but under floor where button pushed with direction upward or is idle
				dont take it
			else
				take it
		isOrderBelow(currentMap)
			if elevator below me but over floor where button pushed with direction upward or is idle
				dont take it
			else
				take it

}


func sortOrders(m def.ElevMap)[]int {

	switch m[def.MY_ID].Dir{
	case def.UP:

	case def.DOWN:

	case def.STILL:

	}

	var orders []int

	for f := m[def.MY_ID].Pos ; f < def.FLOORS ; f++{
		if m[def.MY_ID].Buttons[f][def.UP_BUTTON]  == 1 || m[def.MY_ID].Buttons[f][def.PANEL_BUTTON]  == 1  {
			orders = append(orders, f)
		}
	}

	for f := def.FLOORS; f > m[def.MY_ID].Pos ; f--{
		if m[def.MY_ID].Buttons[f][def.DOWN_BUTTON]  == 1{
			orders = append(orders, f)
		}
	}

	return orders
}

func sortOrdersBelow(m def.ElevMap)[]int {

	var orders []int

	for f := m[def.MY_ID].Pos ; f > -1 ; f--{
		if m[def.MY_ID].Buttons[f][def.UP_BUTTON]  == 1 || m[def.MY_ID].Buttons[f][def.PANEL_BUTTON]  == 1  {
			orders = append(orders, f)
		}
	}

	for f := 0; f < m[def.MY_ID].Pos ; f++{
		if m[def.MY_ID].Buttons[f][def.DOWN_BUTTON]  == 1{
			orders = append(orders, f)
		}
	}

	return orders
}

func chooseDirection(m def.ElevMap, ordes []int) bool {

	myPos := m[def.MY_ID].Pos

	for i, order = range(orders){

		if myPos < order { //we goin upp
			for e := 0 ; e < def.ELEVATORS ; e++{
				if e != def.MY_ID{
					if m[e].Pos > myPos && m[e].Pos < order && (m[e].Dir == def.UP || m[e].State == def.IDLE){
						return false
					} else if m[e].Pos == m[def.MY_ID].Pos  && (m[e].Dir == def.UP || (m[e].State == def.IDLE && e < def.MY_ID){
						return false	
					}
				}
			}
		} else { //we goin down
			for e := 0 ; e < def.ELEVATORS ; e++{
				if e != def.MY_ID{
					if m[e].Pos < myPos && m[e].Pos > order && (m[e].Dir == def.DOWN || m[e].State == def.IDLE){
						return false
					}  else if m[e].Pos == m[def.MY_ID].Pos  && (m[e].Dir == def.DOWN || (m[e].State == def.IDLE && e < def.MY_ID){
						return false	
					}
				}
			}
		}
	}
	return true
}