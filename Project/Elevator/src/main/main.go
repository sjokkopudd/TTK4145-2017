package main

import (
	//"bufio"
	//"fmt"
	"elevatorMap"
	//"network"
	//"os"
)

func main() {
	newEventCh := make(chan int)
	elevatorMap.InitMap(newEventCh)
}
