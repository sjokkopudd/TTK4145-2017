# Module interfaces

* udpNetwork
* taskHandler
* localElevatorDriver
* localElevatorPanel

## udpNetwork

`void startNetworkComunication()`
Starts UDP communication over the network.

`void bufferEvent(event)`
Puts an event into the network module's buffer.

## map

`void newEvent(event)`
Updates the local event-map with a new event.



## taskHandler

`void startElevator()`
Starts operation of the local elevator.

## localElevatorDriver

`void setDirection(dir)`
Sets the direction of the elevator (up, down, stop)

`bool readFloorSensor(floor)`
Returns if the elevator is on given floor.

## localElevatorButtons

`void startPollButtons()`
Starts a loop continuously polling the buttons.

`void setButtonLight(button)`
Sets the lights on a button.