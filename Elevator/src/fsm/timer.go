package fsm

import (
	"time"
)

var timeoutFlag bool

var t time.Time

func timer(timeoutChan chan bool) {
	for {
		if timeoutFlag {
			t = time.Now()
			timeoutFlag = false
		}
		if time.Now().Second()-t.Second() > 1 {
			timeoutChan <- true
		}
	}
}

func timerStart() {
	timeoutFlag = true
}
