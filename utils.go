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

type Error struct {
	e string
	file string
	line int
}

func (e *Error) Error() string {
	return e.e + "\nat " + e.file + ":"+ strconv.Itoa(e.line)
}

func NewError(str string) *Error {
	err := &Error{
		e:str,
	}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		err.file = file
		err.line = line
	}

	return err
}
