// +build race

package racy

import "testing"

func TestPonies(t *testing.T) {
	o := Ponies()
	t.Logf("ponies=%v\n", o)
}
