package fsm

import (
	"time"
)

/*
var timeoutFlag bool

var timerActive bool

var t time.Time
*/
var doorTimer float32

func doorTimeout(timeoutChan chan bool, duration int) {
	doorTimer = float32(time.Now().Second())

	time.Sleep(time.Duration(duration) * time.Second)

	if float32(time.Now().Second())-doorTimer >= 1 {

		msg := true
		timeoutChan <- msg
	}
}

/*
func timer(timeoutChan chan bool) {
	for {
		if timeoutFlag {
			t = time.Now()
			timeoutFlag = false
		}
		if time.Now().Second()-t.Second() > 2 && timerActive {
			timerActive = false
			timeoutChan <- true
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func timerStart() {
	timeoutFlag = true
	timerActive = true
}
*/
