package network

import (
	"def"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

//---------------------------------------------------------------
//------------------------- INTERFACE ---------------------------
//---------------------------------------------------------------

func StartNetworkCommunication(transmitChan chan def.ChannelMessage, receiveChan chan def.ChannelMessage, deadElev chan def.ChannelMessage) {

	fmt.Println("My IP: ", def.IPs[def.MY_ID])
	fmt.Println("Trying to setup nettwork connection")

	ackChan := make(chan ackInfo, 100)

	receivedPackages = make(map[string]bool)

	go reciveUdpPacket(receiveChan, ackChan)
	go transmitUdpPacket(transmitChan, ackChan, deadElev)
}

//-------------------------------------------------------------
//------------------------- MODULE ----------------------------
//-------------------------------------------------------------

var receivedPackages map[string]bool

const (
	MAP = 1
	ACK = 2
)

type ackInfo struct {
	IP    string
	Value bool
}

type udpPacket struct {
	Type        int
	SenderIP    string
	PacketID    string
	Data        interface{}
	AckReceived bool
}

func constructUdpPacket(m interface{}) udpPacket {

	id, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	newPacket := udpPacket{
		Type:     MAP,
		SenderIP: def.IPs[def.MY_ID],
		PacketID: string(id),
		Data:     m,
	}

	return newPacket
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
}

func (p udpPacket) sendAck() {
	newPacket := udpPacket{
		Type:     ACK,
		SenderIP: def.IPs[def.MY_ID],
		PacketID: p.PacketID,
		Data:     true,
	}

	newPacket.sendAsJSON(p.SenderIP)
}

func transmitUdpPacket(transmitChan chan def.ChannelMessage, ackChan chan ackInfo, deadElev chan def.ChannelMessage) {
	for {
		select {
		case msg := <-transmitChan:
			localMap := msg.Map.(def.ElevMap)
			for e := 0; e < def.ELEVATORS; e++ {
				if e != def.MY_ID {

					packet := constructUdpPacket(localMap)

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

						fmt.Println("No acknowledge recieved. ", def.IPs[e], " is dead.")

						//elevatorIsDead := def.NewEvent{def.ELEVATOR_DEAD, e}
						//deadElev <- elevatorIsDead
					}
				}
			}
		}

		time.Sleep(2 * time.Millisecond)
	}
}

func reciveUdpPacket(receiveChan chan def.ChannelMessage, ackChan chan ackInfo) {

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

				if !receivedPackages[receivedPacket.PacketID] {

					receivedPackages[receivedPacket.PacketID] = true

					msg := def.ConstructChannelMessage(m, nil)

					receiveChan <- msg

				} else {
					fmt.Println("RECEIVED AN OLD MAP - I THREW IT ON THE GROUND!!!")
				}

				receivedPacket.sendAck()

			case ACK:

				var v bool

				err = json.Unmarshal(data, &v)
				if err != nil {
					log.Fatal(err)
				}

				a := ackInfo{
					IP:    receivedPacket.SenderIP,
					Value: v,
				}

				ackChan <- a

			}
		}

		time.Sleep(2 * time.Millisecond)
	}
}
