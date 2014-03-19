package packet

import (
//	"encoding/binary"
//	 "log"
	"reflect"
)

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
	Speed     float32
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
	Name  string
	Class string
	Level int
}

type LoginOkPacket struct {
}

type SelectCharactorPacket struct {
	Which int
}

var PacketMap map[uint16]reflect.Type

const (
	_             = iota
	PLogin uint16 = iota
	PCharactorInfo
	PSelectCharactor
	PLoginOk
	PTest
	PMax
)

func init() {
	mh.StructToArray = true
	PacketMap = make(map[uint16]reflect.Type)
	PacketMap[PCharactorInfo] = reflect.TypeOf(CharactorInfoPacket{})
	PacketMap[PLogin] = reflect.TypeOf(LoginPacket{})
	PacketMap[PSelectCharactor] = reflect.TypeOf(SelectCharactorPacket{})
	PacketMap[PLoginOk] = reflect.TypeOf(LoginOkPacket{})
}
