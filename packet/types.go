package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	Name string
}

// id | struct of the packet
func Parse(data []byte) (interface{}, error) {
	if len(data) < 2 {
		return nil, errors.New("len packetId error")
	}
	var packetId uint16
	err := binary.Read(bytes.NewReader(data[:2]), binary.BigEndian, &packetId)
	if err != nil {
		return nil, errors.New("read packetId error")
	}

	tp, ok := packetMap[packetId]
	if !ok {
		return nil, errors.New("unknown packetId")
	}
	reader := Reader(data[2:])
	ret, err := Unpack(reflect.New(tp).Interface(), reader)
	
	return ret, err
}

var packetMap map[uint16]reflect.Type
const (
	_ = iota
	PCharactorInfo uint16 = iota
)

func init() {
	packetMap = make(map[uint16]reflect.Type)
	packetMap[PCharactorInfo] = reflect.TypeOf(CharactorInfoPacket{})
}
