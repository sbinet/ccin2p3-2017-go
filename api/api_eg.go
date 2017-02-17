// +build ignore

// package P is a placeholder to run eg, like so:
//  $> eg -t ./api_eg.go .
package P

import (
	"github.com/sbinet/ccin2p3-2017-go/api"
)

func before(x1, y1, x2, y2 float64) (float64, float64) { return api.MidPoint(x1, y1, x2, y2) }
func after(x1, y1, x2, y2 float64) (float64, float64) {
	return api.MidPointXY(api.Point{x1, y1}, api.Point{x2, y2}).XY()
}
