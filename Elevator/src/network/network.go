package network

import (
	"def"
	"elevatorMap"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

func StartNetworkCommunication(msgChan_toNetwork chan def.ChannelMessage, msgChan_fromNetwork chan def.ChannelMessage, deadElevatorChan chan def.ChannelMessage) {

	fmt.Println("My IP: ", IPs[def.MY_ID])
	fmt.Println("Trying to setup nettwork connection")

	ackChan := make(chan ackInfo, 100)

	contactDeadElevCounter = 0

	go reciveUdpPacket(msgChan_fromNetwork, ackChan)
	go transmitUdpPacket(msgChan_toNetwork, ackChan, deadElevatorChan)
}

var contactDeadElevCounter int

const (
	MAP       = 1
	ACK       = 2
	IP_ELEV_1 = "129.241.187.140:20517"
	IP_ELEV_2 = "129.241.187.150:20518"
	IP_ELEV_3 = "129.241.187.154:20519"
	PORT      = ":20517"
)

var IPs = [def.ELEVATORS]string{IP_ELEV_1, IP_ELEV_2, IP_ELEV_3}

type ackInfo struct {
	IP    string
	Value bool
}

type udpPacket struct {
	Type     int
	SenderIP string
	Data     interface{}
}

func constructUdpPacket(localMap interface{}) udpPacket {

	newPacket := udpPacket{
		Type:     MAP,
		SenderIP: IPs[def.MY_ID],
		Data:     localMap,
	}

	return newPacket
}

func (packet udpPacket) sendAsJSON(recipientIP string) bool {

	jsonBuffer, _ := json.Marshal(packet)

	destinationAddr, err := net.ResolveUDPAddr("udp", recipientIP)
	if err != nil {
		return true
	}

	sendConn, err := net.DialUDP("udp", nil, destinationAddr)
	if err != nil {
		return true
	}

	defer sendConn.Close()

	sendConn.Write(jsonBuffer)

	return false
}

func (packet udpPacket) sendAck() {
	newPacket := udpPacket{
		Type:     ACK,
		SenderIP: IPs[def.MY_ID],
		Data:     true,
	}

	newPacket.sendAsJSON(packet.SenderIP)
}

func transmitUdpPacket(msgChan_toNetwork chan def.ChannelMessage, ackChan chan ackInfo, deadElevatorChan chan def.ChannelMessage) {
	for {
		select {
		case msg := <-msgChan_toNetwork:
			localMap := msg.Map.(elevatorMap.ElevMap)
			for e := 0; e < def.ELEVATORS; e++ {
				if e != def.MY_ID && (localMap[e].IsAlive == 1 || contactDeadElevCounter > 5) {

					packet := constructUdpPacket(localMap)

					var ackRecived ackInfo
					var noConnection bool

				WAIT_FOR_ACK:
					for a := 0; a < 5; a++ {

						noConnection = packet.sendAsJSON(IPs[e])

						time.Sleep(200 * time.Millisecond)

						select {
						case ackRecived = <-ackChan:
							if ackRecived.IP == IPs[e] {
								break WAIT_FOR_ACK
							}
						default:

						}
					}

					if !ackRecived.Value || noConnection {

						fmt.Println("No acknowledge recieved. ", IPs[e], " is dead.")

						elevatorIsDead := def.NewEvent{def.ELEVATOR_DEAD, e}
						msg := def.ConstructChannelMessage(nil, elevatorIsDead)
						deadElevatorChan <- msg
					}
				}
			}
			if contactDeadElevCounter > 5 {
				contactDeadElevCounter = 0
			} else {
				contactDeadElevCounter++
			}
		}

		time.Sleep(2 * time.Millisecond)
	}
}

func reciveUdpPacket(msgChan_fromNetwork chan def.ChannelMessage, ackChan chan ackInfo) {

	localAddress, err := net.ResolveUDPAddr("udp", PORT)
	if err != nil {
		log.Fatal(err)
	}
	receiveConnection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		newBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "make run")
		err1 := newBackup.Run()
		if err1 != nil {
			fmt.Println("Unable to spawn backup; you're on your own.")
			log.Fatal(err)
		}
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

				var receivedMap elevatorMap.ElevMap

				err = json.Unmarshal(data, &receivedMap)
				if err != nil {
					log.Fatal(err)
				}

				msg := def.ConstructChannelMessage(receivedMap, nil)

				msgChan_fromNetwork <- msg

				receivedPacket.sendAck()

			case ACK:

				var val bool

				err = json.Unmarshal(data, &val)
				if err != nil {
					log.Fatal(err)
				}

				ack := ackInfo{
					IP:    receivedPacket.SenderIP,
					Value: val,
				}

				ackChan <- ack

			}
		}

		time.Sleep(2 * time.Millisecond)
	}
}
