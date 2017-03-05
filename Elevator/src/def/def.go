package def

const (
	//Physical constants
	ELEVATORS    = 2
	FLOORS       = 4
	UP_BUTTON    = 0
	DOWN_BUTTON  = 1
	PANEL_BUTTON = 2
	BUTTONS      = 3

	//Identification constants
	MY_ID  = 0
	ELEV_1 = "127.0.0.1:20005"
	ELEV_2 = "127.0.0.1:20010"
	ELEV_3 = "127.0.0.1:20015"
	PORT   = ":20005"

	//Event types
	NEWFLOOR      = 0
	BUTTONPUSH    = 1
	DOOR          = 2
	OTHERELEVATOR = 3
	ELEVATORDEAD  = 4
	NEWDIR        = 5

	//Directions and door cases
	UP         = -1
	IDLE       = 0
	DOWN       = 1
	DOOR_OPEN  = 1
	DOOR_CLOSE = 0
)

var IPs = [ELEVATORS]string{ELEV_1, ELEV_2}

type NewEvent struct {
	EventType int
	Data      interface{}
}

type ElevatorInfo struct {
	ID      int
	Buttons [FLOORS][3]int
	Dir     int
	Pos     int
	Door    int
	IsAlive int
}

type Ack struct {
	Msg string
	IP  string
}

type ElevMap [ELEVATORS]ElevatorInfo

func NewCleanElevMap() ElevMap {

	var newMap ElevMap

	for e := 0; e < ELEVATORS; e++ {
		newMap[e].ID = e
		for f := 0; f < FLOORS; f++ {
			for b := 0; b < BUTTONS; b++ {
				newMap[e].Buttons[f][b] = 0
			}
		}
		newMap[e].Dir = 0
		newMap[e].Pos = 0
		newMap[e].Door = 0
		newMap[e].IsAlive = 1
	}
	return newMap
}
