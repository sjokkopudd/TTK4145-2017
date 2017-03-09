package fsm

import (
	"time"
)

var timeoutFlag bool

var timerActive bool

var t time.Time

func timer(timeoutChan chan bool) {
	for {
		if timeoutFlag {
			t = time.Now()
			timeoutFlag = false
		}
		if time.Now().Second()-t.Second() > 1 && timerActive {
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
