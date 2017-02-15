package elevatorMap

import ()

type elevatorInfo struct {
	firstFloorUp    int
	secondFloorUp   int
	secondFloorDown int
	thirdFloorUp    int
	thirdFloorDown  int
	fourthFloorDown int
	elevatorPos     [3]int
	elevatorDir     [3]int
}

var mapArray = [3]elevatorInfo{}

func InitMap(newEventCh chan) {
	readBackup()
}

