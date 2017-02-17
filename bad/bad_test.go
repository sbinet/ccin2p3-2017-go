package bad_test

import (
	"testing"

	"github.com/sbinet/ccin2p3-2017-go/bad"
)

func TestSum(t *testing.T) {
	o := bad.Sum(10, 20)
	if got, want := o, 10+20; got != want {
		t.Errorf("bad.Sum: got=%v. want=%v\n", got, want)
	}

	for _, test := range []struct {
		want int
		a, b int
	}{
		{
			a:    10,
			b:    20,
			want: 30,
		},
		{
			a:    42,
			b:    -42,
			want: 0,
		},
	} {
		got := bad.Sum(test.a, test.b)
		if got != test.want {
			t.Errorf("bad.Sum(%d, %d): got=%v. want=%v\n",
				test.a, test.b, got, test.want,
			)
		}
	}
}

func TestMyMult(t *testing.T) {
	if got, want := bad.Mult(2, 10), 20; got != want {
		t.Errorf("bad.Mult: got=%v. want=%v\n", got, want)
	}
}
