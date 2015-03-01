package main

import (
	"fmt"
	"time"
)

func main() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick: // HLxxx
			fmt.Println("tick.")
		case <-boom: // HLxxx
			fmt.Println("BOOM!")
			return
		default:     // HLxxx
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

