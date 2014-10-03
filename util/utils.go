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

func GetPercentValue(value, percent int) int {
    return value * percent / 100
}

func ComputeDirection(originX, originY, destX, destY int) Dir_t {
    const (
        BASIS_DIRECTION_HIGH = 2.0
        BASIS_DIRECTION_LOW  = 0.5
    )

    stepX := destX - originX
    stepY := destY - originY

    var k float64
    if stepX == 0 {
        k = 0
    } else {
        k = float64(stepY) / float64(stepX)
    }

    if stepY == 0 {
        switch {
        case (stepX == 0):
            return DOWN
        case stepX > 0:
            return RIGHT
        default:
            return LEFT
        }
    } else if stepY < 0 {
        switch {
        case (stepX == 0):
            return UP
        case (stepX > 0):
            switch {
            case (k < -BASIS_DIRECTION_HIGH):
                return UP
            case (k < -BASIS_DIRECTION_LOW):
                return RIGHTUP
            default:
                return RIGHT
            }
        default:
            switch {
            case (k > BASIS_DIRECTION_HIGH):
                return UP
            case (k > BASIS_DIRECTION_LOW):
                return LEFTUP
            default:
                return LEFT
            }
        }
    } else {
        switch {
        case (stepX == 0):
            return DOWN
        case (stepX > 0):
            switch {
            case (k > BASIS_DIRECTION_HIGH):
                return DOWN
            case (k > BASIS_DIRECTION_LOW):
                return RIGHTDOWN
            default:
                return RIGHT
            }
        default:
            switch {
            case (k < -BASIS_DIRECTION_HIGH):
                return DOWN
            case (k < -BASIS_DIRECTION_LOW):
                return LEFTDOWN
            default:
                return LEFT
            }
        }
    }
    return DIR_NONE
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
