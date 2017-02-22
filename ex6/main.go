package main

import (
	"encoding/json"
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

	destination_addr, err := net.ResolveUDPAddr("udp", "129.241.187.143:20005")
	if err != nil {
		log.Fatal(err)
	}

	send_conn, err := net.DialUDP("udp", nil, destination_addr)
	if err != nil {
		log.Fatal(err)
	}

	for i := num + 1; ; i++ {
		jsonBuf, _ := json.Marshal(i)
		send_conn.Write(jsonBuf)
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
			json.Unmarshal(buffer[0:length], &num)
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
