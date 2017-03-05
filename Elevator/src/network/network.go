package network

import (
	"def"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

//---------------------------------------------------------------
//------------------------- INTERFACE ---------------------------
//---------------------------------------------------------------

func StartNetworkCommunication(transmitChan chan def.ElevMap, receiveChan chan def.ElevMap, deadElev chan def.NewEvent) {

	fmt.Println("My IP: ", def.IPs[def.MY_ID])
	fmt.Println("Trying to setup nettwork connection")

	ackChan := make(chan ackInfo)

	go reciveUdpPacket(receiveChan, ackChan)
	go transmitUdpPacket(transmitChan, ackChan, deadElev)
}

//-------------------------------------------------------------
//------------------------- MODULE ----------------------------
//-------------------------------------------------------------

const (
	MAP = 1
	ACK = 2
)

type ackInfo struct {
	IP    string
	Value bool
}

type udpPacket struct {
	Type int
	IP   string
	Data interface{}
}

func constructUdpPacket(t int, m interface{}) udpPacket {

	switch t {
	case MAP:
		newPacket := udpPacket{
			Type: t,
			IP:   def.IPs[def.MY_ID],
			Data: m,
		}
		return newPacket
	case ACK:
		newPacket := udpPacket{
			Type: t,
			IP:   def.IPs[def.MY_ID],
			Data: true,
		}
		return newPacket
	default:
		newPacket := udpPacket{
			Type: t,
			IP:   def.IPs[def.MY_ID],
			Data: false,
		}
		return newPacket
	}
}

func (p udpPacket) sendAsJSON(r string) {

	json_buffer, _ := json.Marshal(p)

	destination_addr, err := net.ResolveUDPAddr("udp", r)
	if err != nil {
		log.Fatal(err)
	}

	send_conn, err := net.DialUDP("udp", nil, destination_addr)
	if err != nil {
		log.Fatal(err)
	}

	defer send_conn.Close()

	send_conn.Write(json_buffer)

	switch p.Type {
	case MAP:
		fmt.Println("Map sent \n")

	case ACK:
		fmt.Println("Ack sent \n")
	}
}

func transmitUdpPacket(transmitChan chan def.ElevMap, ackChan chan ackInfo, deadElev chan def.NewEvent) {
	for {
		select {
		case mapArray := <-transmitChan:
			for e := 0; e < def.ELEVATORS; e++ {
				if e != def.MY_ID {

					packet := constructUdpPacket(MAP, mapArray)

					var ackRecived ackInfo

				WAIT_FOR_ACK:
					for a := 0; a < 5; a++ {

						packet.sendAsJSON(def.IPs[e])

						time.Sleep(100 * time.Millisecond)

						select {
						case ackRecived = <-ackChan:
							if ackRecived.IP == def.IPs[e] {
								break WAIT_FOR_ACK
							}
						default:

						}
					}

					if !ackRecived.Value {

						fmt.Println("elev dead")

						//elevatorIsDead := def.NewEvent{def.ELEVATOR_DEAD, e}
						//deadElev <- elevatorIsDead
					}
				}
			}
		}
	}
}

func reciveUdpPacket(receiveChan chan def.ElevMap, ackChan chan ackInfo) {

	localAddress, err := net.ResolveUDPAddr("udp", def.PORT)
	if err != nil {
		log.Fatal(err)
	}
	receiveConnection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		log.Fatal(err)
	}

	receiveBuffer := make([]byte, 1024)

	for {

		receiveConnection.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

		n, _, err := receiveConnection.ReadFromUDP(receiveBuffer)

		if n > 0 {

			var data json.RawMessage
			receivedPacket := udpPacket{
				Data: &data,
			}

			err = json.Unmarshal(receiveBuffer[0:n], &receivedPacket)
			if err != nil {
				log.Fatal(err)
			}

			switch receivedPacket.Type {
			case MAP:

				var m def.ElevMap

				err = json.Unmarshal(data, &m)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Recived map:\n", m, "\n")

				//receiveChan <- m

				sendAcknowledge(receivedPacket.IP)

			case ACK:

				var v bool

				err = json.Unmarshal(data, &v)
				if err != nil {
					log.Fatal(err)
				}

				a := ackInfo{
					IP:    receivedPacket.IP,
					Value: v,
				}

				fmt.Println("Recived ack:\n", a, "\n")

				ackChan <- a

			}
		}
	}
}

func sendAcknowledge(ip string) {
	packet := constructUdpPacket(ACK, nil)
	packet.sendAsJSON(ip)

}
