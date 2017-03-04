package def

const (
	ELEVATORS = 1
	FLOORS    = 4
	UP        = 0
	DOWN      = 1
	PANEL     = 2
	BUTTONS   = 3
	MY_IP     = "129.241.187.157"
	MAP_PORT  = ":20005"
	ACK_PORT  = ":30005"
	NEWFLOOR = 0
	BUTTONPUSH = 1
	CLOSEDOOR = 2
)

var IPs = [ELEVATORS]string{MY_IP}

type NewHardwareEvent struct {
	Type   int
	Pos    int
	Floor  int
	Button int
	Door   int
}

type ElevatorInfo struct {
	IP      string
	Buttons [FLOORS][3]int
	Dir     int
	Pos     int
	Door    int
}

type Ack struct {
	Msg string
	IP  string
}

type ElevMap map[string]*ElevatorInfo
