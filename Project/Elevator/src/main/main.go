package main

import (
	//"bufio"
	//"fmt"
	"elevatorMap"
	//"network"
	//"os"
)

func main() {
	/*
		r := make(chan string)

		s := make(chan string)

		go network.StartNetworkCommunication(r, s)

		for {
			reader := bufio.NewReader(os.Stdin)

			s_msg, _ := reader.ReadString('\n')
			s <- s_msg

			msg := <-r
			fmt.Println("Recived message: ", msg)
		}

	*/
	//hardware.InitHardware()

	mapArray := elevatorMap.ReadBackup()
	elevatorMap.WriteBackup(mapArray)
}
