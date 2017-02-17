// package racy shows how to use the race detector.
package racy

// Ponies runs racy goroutines.
func Ponies() int {
	x := 42
	y := -42
	go func() { x++ }()
	go func() {
		if x > 42 {
			x += 4
		}
	}()
	go func() { y-- }()
	return x + y
}
