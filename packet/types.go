package packet

import (
	"encoding/binary"
	"errors"
	// "log"
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

// data = id + struct of the packet
func Parse(data []byte) (interface{}, error) {
	if len(data) < 2 {
		return nil, errors.New("len packetId error")
	}

	packetId := binary.BigEndian.Uint16(data[:2])

	_, ok := PacketMap[packetId]
	if !ok {
		return nil, errors.New("unknown packetId")
	}

	//	reader := Reader(data[2:])
	//	ret, err := Unpack(reflect.New(tp).Interface(), reader)

	//	return ret, err
	return nil, nil
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
