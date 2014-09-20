package util

import (
    "runtime"
    "strconv"
)

type Rect struct {
    X   float32
    Y   float32
    W   float32
    H   float32
}

type FPoint struct {
    X   float32
    Y   float32
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

func getPercentValue(value, percent int) int {
    return value * percent / 100
}

func abs(x int) int {
    if x > 0 {
        return x
    }
    return -x
}

func max(x, y int) int {
    if x > y {
        return x
    }
    return y
}

func min(x, y int) int {
    if x < y {
        return x
    }
    return y
}
