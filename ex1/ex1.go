// You can edit this code!
// Click here and start typing.
package main

import (
    "fmt"
    "time"
    "runtime"
)

var a int = 0

func add() {
	for i := 0; i < 10000; i++ {
		a++
	}
}

func sub() {
	for i := 0; i < 10000; i++ {
		a--
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	go add()
	go sub()

	time.Sleep(300*time.Millisecond)

	fmt.Println(a)
}