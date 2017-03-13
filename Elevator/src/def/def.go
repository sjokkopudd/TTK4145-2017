package def

const (
	ELEVATORS    = 2
	FLOORS       = 4
	UP_BUTTON    = 0
	DOWN_BUTTON  = 1
	PANEL_BUTTON = 2
	BUTTONS      = 3

	MY_ID       = 0
	BACKUP_IP   = "127.0.0.1:30000"
	BACKUP_PORT = ":30000"

	FLOOR_ARRIVAL = 1
	BUTTON_PUSH   = 2
	DOOR_TIMEOUT  = 3
	ELEVATOR_DEAD = 4

	UP    = 1
	STILL = 0
	DOWN  = -1

	//Simulator constants
	SIM_SERV_ADDR   = "127.0.0.1:15657"
	USING_SIMULATOR = false
)

type NewEvent struct {
	EventType int
	Data      interface{}
}

type ElevatorInfo struct {
	ID        int
	Buttons   [FLOORS][BUTTONS]int
	Direction int
	Position  int
	Door      int
	IsAlive   int
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
