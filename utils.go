package ouster

import (
	"math/rand"
	"runtime"
	"strconv"
)

type Rect struct {
	X float32
	Y float32
	W float32
	H float32
}

type Point struct {
	X int
	Y int
}

type FPoint struct {
	X float32
	Y float32
}

type Creature interface {
	Strength() int
	Agility() int
	Intelligence() int

	Damage() int
	Dodge() int
	ToHit() int
}

func Distance(p1, p2 FPoint) float32 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return dx*dx + dy*dy
}

type Error struct {
	e    string
	file string
	line int
}

func (e *Error) Error() string {
	return e.e + "\nat " + e.file + ":" + strconv.Itoa(e.line)
}

func NewError(str string) *Error {
	err := &Error{
		e: str,
	}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		err.file = file
		err.line = line
	}

	return err
}

// if tohit == dodge, the default formula is 0.85
// if tohit < dodge, then tohit / dodge should be primary factor, also take other factor into consideration
// if tohit > dodge, then the differential should be important, also dodge.
func HitTest(tohit int, dodge int) bool {
	var prob float32
	if tohit < dodge {
		prob = 0.85*float32(tohit)/float32(dodge) - 0.15*float32(dodge-tohit)/float32(tohit)
	} else {
		prob = 0.85 + 0.15*float32(tohit-dodge)/float32(dodge)
	}

	return rand.Float32 < prob
}
