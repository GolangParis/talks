package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)

// walk START OMIT

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	recWalk(t, ch)
	// closing the channel so range can finish
	close(ch)
}

// recWalk walks recursively through the tree and push values to the channel
// at each recursion
func recWalk(t *tree.Tree, ch chan int) {
	if t != nil {
		// send the left part of the tree to be iterated over first
		recWalk(t.Left, ch)
		// push the value to the channel
		ch <- t.Value
		// send the right part of the tree to be iterated over last
		recWalk(t.Right, ch)
	}
}

// walk END OMIT

// same START OMIT

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	// iterate over the first channel
	for i := range ch1 {
		// if the value of the second channel doesn't match
		if i != <-ch2 {
			return false
		}
	}
	return true
}

// same END OMIT

// main START OMIT
func main() {
	ch := make(chan int)
	go Walk(tree.New(1), ch)
	for v := range ch {
		fmt.Println(v)
	}
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
// main END OMIT
