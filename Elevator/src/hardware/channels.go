package hardware

import (
	"def"
)

//Stole this file from Morten Fyhn (www.github.com/mortenfyhn)

// In port 4
const (
	PORT4           = 3
	OBSTRUCTION     = (0x300 + 23)
	STOP            = (0x300 + 22)
	BUTTON_COMMAND1 = (0x300 + 21)
	BUTTON_COMMAND2 = (0x300 + 20)
	BUTTON_COMMAND3 = (0x300 + 19)
	BUTTON_COMMAND4 = (0x300 + 18)
	BUTTON_UP1      = (0x300 + 17)
	BUTTON_UP2      = (0x300 + 16)
)

// In port 1
const (
	PORT1         = 2
	BUTTON_DOWN2  = (0x200 + 0)
	BUTTON_UP3    = (0x200 + 1)
	BUTTON_DOWN3  = (0x200 + 2)
	BUTTON_DOWN4  = (0x200 + 3)
	SENSOR_FLOOR1 = (0x200 + 4)
	SENSOR_FLOOR2 = (0x200 + 5)
	SENSOR_FLOOR3 = (0x200 + 6)
	SENSOR_FLOOR4 = (0x200 + 7)
)

// Out port 3
const (
	PORT3          = 3
	MOTORDIR       = (0x300 + 15)
	LIGHT_STOP     = (0x300 + 14)
	LIGHT_COMMAND1 = (0x300 + 13)
	LIGHT_COMMAND2 = (0x300 + 12)
	LIGHT_COMMAND3 = (0x300 + 11)
	LIGHT_COMMAND4 = (0x300 + 10)
	LIGHT_UP1      = (0x300 + 9)
	LIGHT_UP2      = (0x300 + 8)
)

// Out port 2
const (
	PORT2            = 3
	LIGHT_DOWN2      = (0x300 + 7)
	LIGHT_UP3        = (0x300 + 6)
	LIGHT_DOWN3      = (0x300 + 5)
	LIGHT_DOWN4      = (0x300 + 4)
	LIGHT_DOOR_OPEN  = (0x300 + 3)
	LIGHT_FLOOR_IND2 = (0x300 + 1)
	LIGHT_FLOOR_IND1 = (0x300 + 0)
)

// Out port 0
const (
	PORT0 = 1
	MOTOR = (0x100 + 0)
)

// Non-existing ports = (for alignment)
const (
	BUTTON_DOWN1 = -1
	BUTTON_UP4   = -1
	LIGHT_DOWN1  = -1
	LIGHT_UP4    = -1
)

var lightChannelMatrix = [def.FLOORS][3]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var buttonChannelMatrix = [def.FLOORS][3]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}
