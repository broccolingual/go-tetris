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
	QUEUE_MAX    = 5
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

var MINOS = [8][4][2]int{
	{{0, 0}, {0, 0}, {0, 0}, {0, 0}},    // VOID
	{{0, 0}, {-1, 0}, {1, 0}, {2, 0}},   // I
	{{0, 0}, {0, -1}, {1, 0}, {1, -1}},  // O
	{{0, 0}, {-1, 0}, {0, -1}, {1, -1}}, // S
	{{0, 0}, {-1, -1}, {0, -1}, {1, 0}}, // Z
	{{0, 0}, {-1, -1}, {-1, 0}, {1, 0}}, // J
	{{0, 0}, {-1, 0}, {1, 0}, {1, -1}},  // L
	{{0, 0}, {-1, 0}, {0, -1}, {1, 0}},  // T
}

var FIELDS [FIELD_HEIGHT][FIELD_WIDTH]int // y, x

type Target struct {
	Type int
	Data [4][2]int
	X    int
	Y    int
}

type Queue struct {
	Elements []Target
}

func (q Queue) IsEmpty() bool {
	return len(q.Elements) == 0
}

func (q Queue) IsFull() bool {
	return len(q.Elements) == QUEUE_MAX
}

func (q *Queue) Enqueue(t *Target) bool {
	if q.IsFull() {
		return false
	}
	q.Elements = append(q.Elements, *t)
	return true
}

func (q *Queue) Dequeue() *Target {
	if q.IsEmpty() {
		return nil
	}
	dequeueData := q.Elements[0]
	q.Elements = q.Elements[1:]
	return &dequeueData
}

func (t *Target) InitTargetMino(x int, y int) {
	i := rand.Intn(7) + 1
	t.Type = i
	t.Data = MINOS[i]
	t.X = x
	t.Y = y
}

func (t Target) SetMino2Field() {
	for _, cood := range t.Data {
		FIELDS[t.Y+cood[1]][t.X+cood[0]] = t.Type
	}
}

func (t Target) CanMove(dx int, dy int) bool {
	for _, cood := range t.Data {
		nx := t.X + cood[0] + dx
		ny := t.Y + cood[1] + dy
		if ny < 0 || FIELD_HEIGHT <= ny || nx < 0 || FIELD_WIDTH <= nx || FIELDS[ny][nx] > 10 {
			return false
		}
	}
	return true
}

func (t Target) CanRotateRight() bool {
	for _, cood := range t.Data {
		nx := t.X - cood[1]
		ny := t.Y + cood[0]
		if ny < 0 || FIELD_HEIGHT <= ny || nx < 0 || FIELD_WIDTH <= nx || FIELDS[ny][nx] > 10 {
			return false
		}
	}
	return true
}

func (t *Target) RotateRight() {
	for i, cood := range t.Data {
		t.Data[i][0] = -cood[1]
		t.Data[i][1] = cood[0]
	}
}

func (t Target) IsTouching() bool {
	for _, cood := range t.Data {
		nx := t.X + cood[0]
		ny := t.Y + cood[1] + 1
		if ny == FIELD_HEIGHT || FIELDS[ny][nx] > 10 {
			return true
		}
	}
	return false
}

func (t Target) Fix() {
	for _, cood := range t.Data {
		FIELDS[t.Y+cood[1]][t.X+cood[0]] = t.Type + 10
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
	for dy := 0; dy < FIELD_HEIGHT; dy++ {
		for dx := 0; dx < FIELD_WIDTH; dx++ {
			col := FIELDS[dy][dx]
			if col > 10 {
				col -= 10
			}
			DrawString(x+dx*2, y+dy, "  ", coldef, C_PAIRS[col])
		}
	}
}

func DrawNextMinos(x int, y int, tQueue Queue) {
	for i, t := range tQueue.Elements {
		for _, cood := range t.Data {
			DrawString(x+cood[0]*2, y+cood[1]+i*4, "  ", coldef, C_PAIRS[t.Type])
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()
	// termbox.SetOutputMode(termbox.Output256)

	var tQueue Queue

	w, h := termbox.Size()
	var x int = w / 2
	var y int = h / 2
	var cx int = x - FIELD_WIDTH/2
	var cy int = y - FIELD_HEIGHT/2

TargetSetLoop:
	for {
		t := new(Target)
		t.InitTargetMino(4, 1)
		if !tQueue.Enqueue(t) {
			break TargetSetLoop
		}
	}

	t := tQueue.Dequeue()
	nt := new(Target)
	nt.InitTargetMino(4, 1)
	tQueue.Enqueue(nt)
	t.SetMino2Field()
	termbox.Clear(coldef, coldef)
	DrawField(cx, cy)
	DrawNextMinos(cx+FIELD_WIDTH+14, cy+2, tQueue)
	termbox.Flush()

MainLoop:
	for {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyEsc:
				break MainLoop
			case termbox.KeyArrowDown:
				if t.CanMove(0, 1) {
					t.Y++
				}
			case termbox.KeyArrowUp:
				if t.CanRotateRight() {
					t.RotateRight()
				}
			case termbox.KeyArrowRight:
				if t.CanMove(1, 0) {
					t.X++
				}
			case termbox.KeyArrowLeft:
				if t.CanMove(-1, 0) {
					t.X--
				}
			case termbox.KeySpace:

			}
		}
		if t.IsTouching() {
			t.Fix()
			t = tQueue.Dequeue()
			nt := new(Target)
			nt.InitTargetMino(4, 1)
			tQueue.Enqueue(nt)
		}
		ClearField()
		t.SetMino2Field()
		termbox.Clear(coldef, coldef)
		DrawField(cx, cy)
		DrawNextMinos(cx+FIELD_WIDTH+14, cy+2, tQueue)
		termbox.Flush()
	}
}
