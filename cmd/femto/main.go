// Copyright 2015 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Program femto is a very simple text editor.
// It can display text from a file on disk and give the byte position of
// a word.
// Point the mouse at the character and the byte offset of that character
// on disk will be displayed.
//
// Hit 'Q' or 'Escape' to exit.
// Hit 'C' to clear the screen.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"unicode/utf8"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"

	"github.com/mattn/go-runewidth"
)

func main() {
	flag.Usage = func() {
		fmt.Printf(`Program femto is a very simple text editor.

It can display text from a file on disk and give the byte position of
a word.
Point the mouse at the character and the byte offset of that character
on disk will be displayed.

Hit 'Q' or 'Escape' to exit.
Hit 'C' to clear the screen.

Usage of femto:

	femto [options] <path-to-file>

Example:

	$> femto ./zoo/p1.go
`)
	}

	run()
}

func run() {

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(2)
	}

	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	defer s.Fini()

	defStyle = tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.Clear()

	posfmt := "Mouse:  %d, %d  "
	buffmt := "Buffer: %d"
	white := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	gray := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).Background(tcell.ColorGray)
	black := tcell.StyleDefault.
		Foreground(tcell.ColorBlack).Background(tcell.ColorRed)

	mx, my := -1, -1
	ox, oy := -1, -1
	bx, by := -1, -1
	w, h := s.Size()

	buf, err := Open(flag.Arg(0), s)
	if err != nil {
		log.Fatal(err)
	}

	buf.display(s, white)
loop:
	for {
		drawBox(s, 1, 1, 42, buf.header, gray, ' ')
		emitStr(s, 2, 2, gray, "Press Q or ESC to exit.")
		emitStr(s, 2, 3, gray, "Press C to clear screen.")
		emitStr(s, 2, 4, gray, fmt.Sprintf(posfmt, mx, my))
		emitStr(s, 2, 5, gray, fmt.Sprintf(buffmt, buf.loc(mx, my)))

		s.Show()
		ev := s.PollEvent()
		st := tcell.StyleDefault.Background(tcell.ColorRed)
		w, h = s.Size()
		buf.update(s, white)

		// always clear any old selection box
		if ox >= 0 && oy >= 0 && bx >= 0 {
			drawSelect(s, ox, oy, bx, by, false)
		}

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			s.SetContent(w-1, h-1, 'R', nil, st)
		case *tcell.EventKey:
			s.SetContent(w-2, h-2, ev.Rune(), nil, st)
			s.SetContent(w-1, h-1, 'K', nil, st)
			switch ev.Key() {
			case tcell.KeyCtrlL:
				s.Sync()
			case tcell.KeyEscape:
				break loop

			default:
				switch ev.Rune() {
				case 'C', 'c':
					s.Clear()
					buf.display(s, white)
				case 'Q', 'q':
					break loop
				}
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			// Only buttons, not wheel events
			button &= tcell.ButtonMask(0xff)
			if button != tcell.ButtonNone && ox < 0 {
				ox, oy = x, y
			}
			switch ev.Buttons() {
			case tcell.ButtonNone:
				if ox >= 0 {
					highlightBox(s, ox, oy, x, y, black, buf)
					ox, oy = -1, -1
					bx, by = -1, -1
				}
			}
			if button != tcell.ButtonNone {
				bx, by = x, y
			}
			s.SetContent(w-1, h-1, 'M', nil, st)
			mx, my = x, y
		}

		if ox >= 0 && bx >= 0 {
			drawSelect(s, ox, oy, bx, by, true)
		}
	}
}

// File represents a file on disk.
// File exposes facilities to match on-disk byte-indices with tokens (words).
type File struct {
	arr    *LineArray
	w, h   int
	header int
	footer int
}

// Open opens a file on disk to be diplayed later on.
func Open(fname string, s tcell.Screen) (*File, error) {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	w, h := s.Size()

	return &File{
		arr:    NewLineArray(buf),
		w:      w,
		h:      h,
		header: 6,
		footer: 0,
	}, nil
}

func (f *File) pos(x, y int) (ix int, iy int, ok bool) {
	if x < 1 || y < 0 || y <= f.header {
		return -1, -1, ok
	}

	iy = y - f.header - 2
	if iy >= len(f.arr.lines) {
		return -1, -1, ok
	}
	ix = x - 1
	if ix >= len(f.arr.lines[iy]) {
		return -1, -1, ok
	}
	ok = true
	return ix, iy, ok
}

func (f *File) txt(x, y int) rune {
	x, y, ok := f.pos(x, y)
	if !ok {
		return 0
	}
	return rune(f.arr.lines[y][x])
}

func (f *File) loc(x, y int) int {
	c := 0
	ix, iy, ok := f.pos(x, y)
	if !ok {
		return -1
	}
	for _, line := range f.arr.lines[:iy] {
		c += len(line) + 1
	}
	c += ix
	return c

}

func (f *File) display(s tcell.Screen, sty tcell.Style) {
	max := f.h - f.footer
	if max > len(f.arr.lines) {
		max = len(f.arr.lines)
	}

	xmax := 0
	for _, line := range f.arr.lines {
		if xlen := len(line); xlen > xmax {
			xmax = xlen
		}
	}

	drawBox(s, 0, f.header+1, xmax+1, f.header+1+len(f.arr.lines), sty, ' ')

	for row := 0; row < max; row++ {
		emitStr(s, 1, row+f.header+2, sty, string(f.arr.lines[row]))
	}
}

func (f *File) update(s tcell.Screen, sty tcell.Style) {
	w, h := s.Size()
	if w == f.w && h == f.h {
		return
	}

	f.w = w
	f.h = h
	f.display(s, sty)
}

func runeToByteIndex(n int, txt []byte) int {
	if n == 0 {
		return 0
	}

	count := 0
	i := 0
	for len(txt) > 0 {
		_, size := utf8.DecodeRune(txt)

		txt = txt[size:]
		count += size
		i++

		if i == n {
			break
		}
	}
	return count
}

// A LineArray simply stores and array of lines and makes it easy to insert
// and delete in it
type LineArray struct {
	lines [][]byte
}

// NewLineArray returns a new line array from an array of bytes
func NewLineArray(text []byte) *LineArray {
	la := new(LineArray)
	// Split the bytes into lines
	split := bytes.Split(text, []byte("\n"))
	la.lines = make([][]byte, len(split))
	for i := range split {
		la.lines[i] = make([]byte, len(split[i]))
		copy(la.lines[i], split[i])
	}

	return la
}

// Substr returns the string representation between two locations
func (la *LineArray) Substr(start, end Loc) string {
	startX := runeToByteIndex(start.X, la.lines[start.Y])
	endX := runeToByteIndex(end.X, la.lines[end.Y])
	if start.Y == end.Y {
		return string(la.lines[start.Y][startX:endX])
	}
	var str string
	str += string(la.lines[start.Y][startX:]) + "\n"
	for i := start.Y + 1; i <= end.Y-1; i++ {
		str += string(la.lines[i]) + "\n"
	}
	str += string(la.lines[end.Y][:endX])
	return str
}

// Loc stores a location
type Loc struct {
	X, Y int
}

var defStyle tcell.Style

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func highlightBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, f *File) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			val := f.txt(col, row)
			s.SetContent(col, row, val, nil, style)
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, r, nil, style)
		}
	}
}

func drawSelect(s tcell.Screen, x1, y1, x2, y2 int, sel bool) {

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			mainc, combc, style, width := s.GetContent(col, row)
			if style == tcell.StyleDefault {
				style = defStyle
			}
			style = style.Reverse(sel)
			s.SetContent(col, row, mainc, combc, style)
			col += width - 1
		}
	}
}
