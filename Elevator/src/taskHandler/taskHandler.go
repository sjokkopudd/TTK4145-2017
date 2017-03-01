package taskHandler

import (
	"def"
)

func TaskHandler(eventChan chan def.NewHardwareEvent) {
	for {

		switch state {
		case IDLE:
			select {
			case newEvent <- eventChan:
				if newEvent.Pos != -1 {

				}
			}
		}
	}
}
