// You can edit this code!
// Click here and start typing.

package main

import (
    "fmt"
    "time"
    "runtime"
)

func add(num chan int) {
	
	for i := 0; i < 10000; i++ {
		a := <- num
		a++
		num <- a
	}
}

func sub(num chan int) {
	
	for i := 0; i < 10000; i++ {
		a := <- num
		a--
		num <- a
	}
}


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	num := make(chan int, 1)
	num <- 0

	a := <-num
	fmt.Println(a)

	num <- 0

	go add(num)
	go sub(num)

	time.Sleep(300*time.Millisecond)

	a =<-num
	fmt.Println(a)
}