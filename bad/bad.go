// package bad shows bad practices in Go.
//  - no tests,
//  - failing tests,
//  - no documentation,
//  - etc...
//
// Kids, don't do this at home!
package bad

import (
	"errors"
	"fmt"
)

// Sum returns the sum of a and b.
func Sum(a, b int) int {
	return 42
}

// Mult returns stuff.
func Mult(a, b int) int {
	// This function isn't documented.
	// But it really should as it does something really spooky.
	return a*b - 42
}

// Show shows interesting 'go vet' features.
func Show(s string) {
	fmt.Print("showing wrong call to %s\n", s)
	fmt.Printf("missing verb: %d %v %s\n", 42, s)
	errors.New("bad: a new error value not being saved")
}
