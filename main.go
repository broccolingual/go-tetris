package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

const (
	FIELD_HEIGHT = 20
	FIELD_WIDTH  = 10
)

const (
	VOID = iota
	I
	O
	S
	Z
	J
	L
	T
)

var C_PAIRS = map[int]termbox.Attribute{
	VOID: termbox.ColorWhite,
	I:    termbox.ColorCyan,
	O:    termbox.ColorYellow,
	S:    termbox.ColorGreen,
	Z:    termbox.ColorRed,
	J:    termbox.ColorBlue,
	L:    termbox.ColorBlack,
	T:    termbox.ColorMagenta,
}

var MINOS = [8][3][2]int8{
	{{0, 0}, {0, 0}, {0, 0}},    // VOID
	{{-1, 0}, {1, 0}, {2, 0}},   // I
	{{0, -1}, {1, 0}, {1, -1}},  // O
	{{-1, 0}, {0, -1}, {1, -1}}, // S
	{{-1, -1}, {0, -1}, {1, 0}}, // Z
	{{-1, -1}, {-1, 0}, {1, 0}}, // J
	{{-1, 0}, {1, 0}, {1, -1}},  // L
	{{-1, 0}, {0, -1}, {1, 0}},  // T
}

var FIELDS [FIELD_HEIGHT][FIELD_WIDTH]int // y, x

type Target struct {
	Type int
	X    int
	Y    int
}

func (t *Target) InitTargetMino(x int, y int) {
	t.Type = rand.Intn(7) + 1
	t.X = x
	t.Y = y
}

func (t Target) SetMino2Field() {
	FIELDS[t.Y][t.X] = t.Type
	for _, cood := range MINOS[t.Type] {
		FIELDS[t.Y+int(cood[1])][t.X+int(cood[0])] = t.Type
	}
}

func (t Target) CanMove(dx int, dy int) bool {
	for _, cood := range MINOS[t.Type] {
		nx := t.X + int(cood[0]) + dx
		ny := t.Y + int(cood[1]) + dy
		if ny < 0 || FIELD_HEIGHT <= ny || nx < 0 || FIELD_WIDTH <= nx || FIELDS[ny][nx] > 10 {
			return false
		}
	}
	return true
}

func (t Target) isTouching() bool {
	for _, cood := range MINOS[t.Type] {
		nx := t.X + int(cood[0])
		ny := t.Y + int(cood[1]) + 1
		if ny == FIELD_HEIGHT || FIELDS[ny][nx] > 10 {
			return true
		}
	}
	return false
}

func (t Target) Fix() {
	FIELDS[t.Y][t.X] = t.Type + 10
	for _, cood := range MINOS[t.Type] {
		FIELDS[t.Y+int(cood[1])][t.X+int(cood[0])] = t.Type + 10
	}
}

func ClearField() {
	for dy := 0; dy < FIELD_HEIGHT; dy++ {
		for dx := 0; dx < FIELD_WIDTH; dx++ {
			if FIELDS[dy][dx] < 10 {
				FIELDS[dy][dx] = 0
			}
		}
	}
}

func DrawString(x int, y int, msg string, fg termbox.Attribute, bg termbox.Attribute) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func DrawField(x int, y int) {
	termbox.Clear(coldef, coldef)
	for dy := 0; dy < FIELD_HEIGHT; dy++ {
		for dx := 0; dx < FIELD_WIDTH; dx++ {
			col := FIELDS[dy][dx]
			if col > 10 {
				col -= 10
			}
			DrawString(x+dx*2, y+dy, "  ", coldef, C_PAIRS[col])
		}
	}
	termbox.Flush()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()
	// termbox.SetOutputMode(termbox.Output256)

	w, h := termbox.Size()
	var x int = w / 2
	var y int = h / 2
	var cx int = x - FIELD_WIDTH/2
	var cy int = y - FIELD_HEIGHT/2
	t := new(Target)
	t.InitTargetMino(4, 1)
	t.SetMino2Field()
	DrawField(cx, cy)

mainloop:
	for {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowDown:
				if t.CanMove(0, 1) {
					t.Y++
				}
			case termbox.KeyArrowUp:
				if t.CanMove(0, -1) {
					t.Y--
				}
			case termbox.KeyArrowRight:
				if t.CanMove(1, 0) {
					t.X++
				}
			case termbox.KeyArrowLeft:
				if t.CanMove(-1, 0) {
					t.X--
				}
			}
		}
		if t.isTouching() {
			t.Fix()
			t.InitTargetMino(4, 1)
		}
		ClearField()
		t.SetMino2Field()
		DrawField(cx, cy)
	}
}
