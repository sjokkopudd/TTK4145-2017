package def

const (
	//Physical constants
	ELEVATORS    = 3
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
	FLOOR_ARRIVAL = 1
	BUTTON_PUSH   = 2
	DOOR_TIMEOUT  = 3
	ELEVATOR_DEAD = 4

	//Directions and door cases
	UP    = 1
	STILL = 0
	DOWN  = -1

	DOOR_CLOSED = 0
	DOOR_OPEN   = 1

	//Simulator constants
	SIM_SERV_ADDR   = "127.0.0.1:15656"
	USING_SIMULATOR = false
)

var IPs = [ELEVATORS]string{ELEV_1, ELEV_2, ELEV_3}

type NewEvent struct {
	EventType int
	Data      interface{}
}

type ElevatorInfo struct {
	ID      int
	Buttons [FLOORS][BUTTONS]int
	Dir     int
	Pos     int
	Door    int
	IsAlive int
}

type ElevMap [ELEVATORS]ElevatorInfo

func NewCleanElevMap() *ElevMap {

	newMap := new(ElevMap)

	for e := 0; e < ELEVATORS; e++ {
		newMap[e].ID = e
		for f := 0; f < FLOORS; f++ {
			for b := 0; b < BUTTONS; b++ {
				newMap[e].Buttons[f][b] = 0
			}
		}
		newMap[e].Dir = STILL
		newMap[e].Pos = 0
		newMap[e].Door = -1
		newMap[e].IsAlive = 0
	}
	return newMap
}

type ChannelMessage struct {
	Map   interface{}
	Event interface{}
}

func ConstructChannelMessage(m interface{}, e interface{}) ChannelMessage {
	newChannelMessage := ChannelMessage{
		Map:   m,
		Event: e,
	}

	return newChannelMessage
}
