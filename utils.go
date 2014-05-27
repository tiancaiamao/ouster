package ouster

import (
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
	HP() int
	Damage() int
	Defense() int
	Dodge() int
	ToHit() int
}

func Distance2(p1, p2 FPoint) float32 {
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

func Encrypt(ZoneID uint16, ServerID uint16) uint8 {
	return uint8(((ZoneID >> 8) ^ ZoneID) ^ ((ServerID+1) << 4))
}