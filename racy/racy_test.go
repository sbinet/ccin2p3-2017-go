// +build !race

package racy

import "testing"

func TestDeadlock(t *testing.T) {
	x := make(chan int)
	y := make(chan int)

	go func() {
		y <- <-x
	}()
	// x <- 1 // dead-lock
	oops := <-y // dead lock
	t.Logf("oops: %v\n", oops)
}
