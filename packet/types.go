package packet

import "errors"

type DirectionType uint8
const (
	UP = iota
	RU
	RIGHT
	RD
	DOWN
	LD
	LEFT
	LU
)
type MovePacket struct {
	Speed float32
	Direction DirectionType
}
type SkillPacket struct {
	Id int
	
}
type LoginPacket struct {
	Username string
	Password string
}

type CharactorInfoPacket struct {
	Name string
	Class string
	Level int
}

type LoginOkPacket struct {
}

type SelectCharactorPacket struct {
	Name string
}

func Parse(data []byte) (interface{}, error) {
	return nil, errors.New("not implement yet")
}
