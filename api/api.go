// package api exposes a bad api that need to be refactored with 'eg'.
package api

import "fmt"

// MidPoint returns the mid-point of the segment (x1,y1)-(x2,y2).
func MidPoint(x1, y1, x2, y2 float64) (float64, float64) {
	return (x1 + x2) / 2, (y1 + y2) / 2
}

// MidPointXY is a better api.
func MidPointXY(p1, p2 Point) Point {
	return Point{
		X: 0.5 * (p1.X + p2.X),
		Y: 0.5 * (p1.Y + p2.Y),
	}
}

// Point represents a 2D point
type Point struct {
	X float64
	Y float64
}

func (p Point) XY() (float64, float64) {
	return p.X, p.Y
}

// useMidPoint is a simple function that uses the old MidPoint function
func useMidPoint() {
	x1, y1 := 10.0, 10.0
	x2, y2 := 20.0, 20.0

	x3, y3 := MidPoint(x1, y1, x2, y2)
	fmt.Printf("x3,y3 = (%v, %v)\n", x3, y3)
}
