package ouster

import "runtime"

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
//	return e.err + " in function " + e.fun + "at " + e.file + ":"+ e.line
	return "hello"
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
