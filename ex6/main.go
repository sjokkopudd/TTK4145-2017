package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

func primary(num int) {
	newBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err := newBackup.Run()
	if err != nil {
		log.Fatal(err)
	}

	destination_addr, err := net.ResolveUDPAddr("udp", "129.241.187.255:20005")
	if err != nil {
		log.Fatal(err)
	}

	send_conn, err := net.DialUDP("udp", nil, destination_addr)
	if err != nil {
		log.Fatal(err)
	}

	msg := make([]byte, 16)

	for i := num + 1; ; i++ {
		msg[0] = byte(i)
		msg[1] = 0
		send_conn.Write(msg)
		fmt.Println(i)
		time.Sleep(1000 * time.Millisecond)
	}

}

func backup() int {

	num := 0

	addr, err := net.ResolveUDPAddr("udp", ":20005")
	if err != nil {
		log.Fatal(err)
	}

	listenCon, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer listenCon.Close()

	buffer := make([]byte, 16)

	for {
		listenCon.SetReadDeadline(time.Now().Add(1500 * time.Millisecond))
		length, _, err := listenCon.ReadFromUDP(buffer[:])

		if length > 0 {
			num = int(buffer[0])
			fmt.Println(num)
			if err != nil {
				fmt.Println("error: ")
				log.Fatal(err)
			}
		} else {
			fmt.Println("No signal found. Creating primary.")
			return num
		}
	}
}

func main() {

	num := backup()
	primary(num)

}
