package main

import "fmt"

// main START OMIT
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x: // HLxxx
			x, y = y, x+y
		case <-quit: // HLxxx
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0 // HLxxx
	}()
	fibonacci(c, quit)
}
// main END OMIT
