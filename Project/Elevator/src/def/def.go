package def

const Elevators = 2
const Floors = 4
const MyIP = "129.241.187.150"
const MapPort = ":20005"
const AcknowledegePort = ":30005"

var IPs = [Elevators]string{MyIP, "129.241.187.149"}

type NewHardwareEvent struct {
	Pos    int
	Floor  int
	Button int
}

type ElevatorInfo struct {
	ID      int
	Buttons [Floors][3]int
	Dir     int
	Pos     int
}

type Ack struct {
	Msg string
	IP  string
}

type ElevMap map[string]*ElevatorInfo
