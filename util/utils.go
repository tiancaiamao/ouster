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

func computeFinalDamage(minDamage maxDamage, realDamage Damage_t, protection Protection_t, bCritical bool) {
    // 致命一击无视防御
    if bCritical {
        return realDamage
    }

    finalDamage := realDamage - (realDamage*(Protection/8))/100

    // avgDamage := (minDamage + maxDamage) / 2
    //
    //     	int      DamageRatio = 100;
    //
    //     	if (Protection < avgDamage)
    //     	{
    //     		DamageRatio = 100;
    //     	}
    //     	else if (Protection < getPercentValue(avgDamage, 150))
    //     	{
    //     		DamageRatio = 90;
    //     	}
    //     	else if (Protection < getPercentValue(avgDamage, 200))
    //     	{
    //     		DamageRatio = 80;
    //     	}
    //     	else if (Protection < getPercentValue(avgDamage, 250))
    //     	{
    //     		DamageRatio = 70;
    //     	}
    //     	else if (Protection < getPercentValue(avgDamage, 300))
    //     	{
    //     		DamageRatio = 60;
    //     	}
    //     	else
    //     	{
    //     		DamageRatio = 50;
    //     	}

    return max(1, getPercentValue(realDamage, DamageRatio))
}
