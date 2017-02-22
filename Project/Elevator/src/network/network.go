package network

import (
	"def"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func StartNetworkCommunication(transmitChan chan def.ElevMap, receiveChan chan def.ElevMap) {

	fmt.Println("Trying to setup nettwork connection")
	go reciveMap(receiveChan)
	for {
		select {
		case mapArray := <-transmitChan:
			go transmitMap(mapArray)
		}
		time.Sleep(100 * time.Millisecond)
	}

}

func transmitMap(mapArray def.ElevMap) {

	for i := 0; i < def.ELEVATORS; i++ {

		if def.IPs[i] != def.MY_IP {

			destination_addr, err := net.ResolveUDPAddr("udp", def.IPs[i]+def.MAP_PORT)
			if err != nil {
				log.Fatal(err)
			}

			send_conn, err := net.DialUDP("udp", nil, destination_addr)
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(500 * time.Millisecond)

			for j := 0; j < 5; j++ {

				json_buffer, _ := json.Marshal(mapArray)

				if len(json_buffer) > 0 {
					send_conn.Write(json_buffer)
					var m def.ElevMap
					err = json.Unmarshal(json_buffer, &m)
					if err != nil {
						log.Fatal(err)
					}

				}

				if receiveAcknowledge(def.IPs[i]) {
					break
				}
			}
			//elevator i is dead
		}

	}

}

func receiveAcknowledge(senderIP string) bool {
	localAddress, err := net.ResolveUDPAddr("udp", def.ACK_PORT)
	if err != nil {
		log.Fatal(err)
	}
	reciveConnection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		log.Fatal(err)
	}

	defer reciveConnection.Close()

	receiveBuffer := make([]byte, 1024)

	time.Sleep(50 * time.Millisecond)

	reciveConnection.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	n, _, err := reciveConnection.ReadFromUDP(receiveBuffer)
	fmt.Println("here")
	if n > 0 {
		var ackMsg def.Ack
		err = json.Unmarshal(receiveBuffer[0:n], &ackMsg)
		fmt.Println(ackMsg.IP)
		if err != nil {
			log.Fatal(err)
		}
		if ackMsg.Msg == "Ack" && ackMsg.IP == senderIP {
			fmt.Println("Acknowledge received from " + ackMsg.IP)
			return true
		}
	}
	return false
}

func reciveMap(receiveChan chan def.ElevMap) {

	localAddress, err := net.ResolveUDPAddr("udp", def.MAP_PORT)
	if err != nil {
		log.Fatal(err)
	}
	receiveConnection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		log.Fatal(err)
	}

	defer receiveConnection.Close()

	receiveBuffer := make([]byte, 1024)

	time.Sleep(500 * time.Millisecond)

	for {
		receiveConnection.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, senderIP, err := receiveConnection.ReadFromUDP(receiveBuffer)

		if n > 0 {

			var receivedMap def.ElevMap
			err = json.Unmarshal(receiveBuffer[0:n], &receivedMap)
			if err != nil {
				log.Fatal(err)
			}
			receiveChan <- receivedMap
			sendAcknowledge(senderIP.IP.String())

		}
	}
}

func sendAcknowledge(ip string) {
	destinationAddress, err := net.ResolveUDPAddr("udp", ip+def.ACK_PORT)
	if err != nil {
		log.Fatal(err)
	}

	transmitConnection, err := net.DialUDP("udp", nil, destinationAddress)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	ackMsg := def.Ack{"Ack", def.MY_IP}
	transmitBuffer, _ := json.Marshal(ackMsg)

	if len(transmitBuffer) > 0 {
		transmitConnection.Write([]byte(transmitBuffer))
		fmt.Println("Sending ack")
	}

}

func transmit(s chan string) {

	destination_addr, err := net.ResolveUDPAddr("udp", "129.241.187.141:20005")
	if err != nil {
		log.Fatal(err)
	}

	send_conn, err := net.DialUDP("udp", nil, destination_addr)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	for {
		s_buffer := <-s

		if len(s_buffer) > 0 {
			send_conn.Write([]byte(s_buffer))
			fmt.Println("Message send: ", s_buffer)
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func recive(r chan string) {

	local_addr, err := net.ResolveUDPAddr("udp", ":20005")
	if err != nil {
		log.Fatal(err)
	}
	recive_conn, err := net.ListenUDP("udp", local_addr)
	if err != nil {
		log.Fatal(err)
	}

	defer recive_conn.Close()

	r_buffer := make([]byte, 1024)

	time.Sleep(500 * time.Millisecond)

	for {
		recive_conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		n, _, err := recive_conn.ReadFromUDP(r_buffer)

		if n > 0 {
			r <- (string(r_buffer[0:n]))
			if err != nil {
				log.Fatal(err)
			}

		}
	}
}
