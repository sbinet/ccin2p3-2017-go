package zoo

import (
	"fmt"
	"io"
)

type Animal interface {
	PArle(w io.Writer) error
	NumAppendage() int
}

type Mouse struct{}

func (m *Mouse) Speak(w io.Writer) error {
	fmt.Fprintf(w, "squeak! squeak!\n")
	return nil
}

func (m *Mouse) NumAppendage() int {
	return m.NumLeg() + m.NumTail()
}

func (m *Mouse) NumLeg() int  { return 4 }
func (m *Mouse) NumTail() int { return 1 }

// explicitly tell *Mouse implements Animal
//var _ Animal = (*Mouse)(nil)

func Chorus(w io.Writer, animals []Animal) error {
	for _, animal := range animals {
		err := animal.PArle(w)
		if err != nil {
			return err
		}
	}
	return nil
}
