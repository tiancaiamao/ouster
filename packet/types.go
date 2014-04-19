package packet

import (
	//	"encoding/binary"
	//	 "log"
	"github.com/tiancaiamao/ouster"
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

type SkillTargetEffectPacket struct {
	Skill int
	From  int
	To    int
	Hurt  int
	Succ  bool
}

type SkillRegionEffectPacket struct {
	Skill int
	From  int
	To    ouster.FPoint
}

type SkillPacket struct {
	Id     int
	Target int
	Region ouster.FPoint
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

type PlayerInfoPacket map[string]interface{}

type SelectCharactorPacket struct {
	Which int
}

type CMovePacket struct {
	X float32
	Y float32
}

type PosSyncPacket struct {
	Cur ouster.FPoint
	To  ouster.FPoint
}

type SMovePacket struct {
	Id  uint32
	Cur ouster.FPoint
	To  ouster.FPoint
}

var PacketMap map[uint16]reflect.Type

const (
	_             = iota
	PLogin uint16 = iota
	PCharactorInfo
	PSelectCharactor
	PLoginOk
	PPlayerInfo
	PCMove
	PSMove
	PPosSync
	PSkillTargetEffect
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
	PacketMap[PPlayerInfo] = reflect.TypeOf(PlayerInfoPacket{})
	PacketMap[PCMove] = reflect.TypeOf(CMovePacket{})
	PacketMap[PSMove] = reflect.TypeOf(SMovePacket{})
	PacketMap[PPosSync] = reflect.TypeOf(PosSyncPacket{})
	PacketMap[PSkillTargetEffect] = reflect.TypeOf(SkillTargetEffectPacket{})
}
