package darkeden

import (
	"encoding/binary"
	"errors"
	"github.com/tiancaiamao/ouster/packet"
	"log"
)

type CGConnectPacket struct {
	Key uint32
	// Slayer or Vampire?
	PCType     uint8
	PCName     string
	MacAddress [4]byte
}

func (connect *CGConnectPacket) Id() packet.PacketID {
	return PACKET_CG_CONNECT
}
func (connect *CGConnectPacket) String() string {
	return "connect"
}
func readConnect(buf []byte, code uint8) (packet.Packet, error) {
	// [ 0 0 0 240 1 4 183 232 191 241 0 80 86 192 0 8]
	ret := new(CGConnectPacket)
	ret.Key = binary.LittleEndian.Uint32(buf[:4])
	ret.PCType = buf[4]
	length := buf[5]
	ret.PCName = string(buf[6 : 6+length])
	copy(ret.MacAddress[:], buf[6+length:])
	return ret, nil
}

type CGReadyPacket struct{}

func (ready CGReadyPacket) Id() packet.PacketID {
	return PACKET_CG_READY
}
func (ready CGReadyPacket) String() string {
	return "ready"
}

type CGMovePacket struct {
	Dir uint8
	X   uint8
	Y   uint8
}

func (move CGMovePacket) Id() packet.PacketID {
	return PACKET_CG_MOVE
}
func (move CGMovePacket) String() string {
	return "move"
}

const (
	dirLEFT      = 0
	dirRIGHT     = 4
	dirUP        = 6
	dirDOWN      = 2
	dirLEFTUP    = 7
	dirRIGHTUP   = 5
	dirLEFTDOWN  = 1
	dirRIGHTDOWN = 3
)

func encryptDir(dir byte) (uint8, error) {
	var ret uint8
	switch dir {
	case 53:
		dir = dirLEFT
	case 49:
		dir = dirRIGHT
	case 51:
		dir = dirUP
	case 55:
		dir = dirDOWN
	case 50:
		dir = dirLEFTUP
	case 48:
		dir = dirRIGHTUP
	case 52:
		dir = dirLEFTDOWN
	case 54:
		dir = dirRIGHTDOWN
	default:
		return ret, errors.New("unknow dir")
	}
	return ret, nil
}

func SHUFFLE_STATEMENT_3(code uint8, A func(), B func(), C func()) {
	switch code % 3 {
	case 0:
		A()
		B()
		C()
	case 1:
		B()
		C()
		A()
	case 2:
		C()
		A()
		B()
	}
	return
}

func SHUFFLE_STATEMENT_2(code uint8, A func(), B func()) {
	switch code % 2 {
	case 0:
		A()
		B()
	case 1:
		B()
		A()
	}
	return
}

func SHUFFLE_STATEMENT_4(code uint8, A func(), B func(), C func(), D func()) {
	switch code % 4 {
	case 0:
		A()
		B()
		C()
		D()
	case 1:
		B()
		C()
		D()
		A()
	case 2:
		C()
		D()
		A()
		B()
	case 3:
		D()
		A()
		C()
		B()
	}
	return
}

func readMove(buf []byte, code uint8) (packet.Packet, error) {
	var ret CGMovePacket
	var err error
	offset := 0
	A := func() {
		ret.X = buf[offset] ^ code
		offset++
	}
	B := func() {
		ret.Y = buf[offset] ^ code
		offset++
	}
	C := func() {
		ret.Dir = buf[offset] ^ code
		offset++
	}
	// encryption...fuck
	SHUFFLE_STATEMENT_3(code, A, B, C)
	if ret.Dir >= 8 {
		err = errors.New("Dir out of range")
	}
	return ret, err
}

type CGVerifyTimePacket struct{}

func (verifyTime CGVerifyTimePacket) Id() packet.PacketID {
	return PACKET_CG_VERIFY_TIME
}
func (verifyTime CGVerifyTimePacket) String() string {
	return "verify time"
}

type CGAttackPacket struct {
	ObjectID uint32
	X        uint8
	Y        uint8
	Dir      uint8
}

func (attack CGAttackPacket) Id() packet.PacketID {
	return PACKET_CG_ATTACK
}
func (attack CGAttackPacket) String() string {
	return "attack"
}
func readAttack(buf []byte, code uint8) (packet.Packet, error) {
	// [188 251 55 82 48 0 0]
	var ret CGAttackPacket
	offset := 0
	A := func() {
		ret.ObjectID = binary.LittleEndian.Uint32(buf[offset:]) ^ uint32(code)
		offset += 4
	}
	B := func() {
		ret.X = buf[offset] ^ code
		offset++
	}
	C := func() {
		ret.Y = buf[offset] ^ code
		offset++
	}
	D := func() {
		ret.Dir = buf[offset] ^ code
		offset++
	}
	SHUFFLE_STATEMENT_4(code, A, B, C, D)
	return ret, nil
}

type CGBloodDrainPacket struct {
	ObjectID uint32
}

func (bloodDrain CGBloodDrainPacket) Id() packet.PacketID {
	return PACKET_CG_BLOOD_DRAIN
}
func (bloodDrain CGBloodDrainPacket) String() string {
	return "blood drain"
}
func readBloodDrain(buf []byte, code uint8) (packet.Packet, error) {
	id := binary.LittleEndian.Uint32(buf)
	return CGBloodDrainPacket{id}, nil
}

type CGLearnSkillPacket struct {
	SkillType       uint16
	SkillDomainType uint8
}

func (learnSkill CGLearnSkillPacket) Id() packet.PacketID {
	return PACKET_CG_LEARN_SKILL
}

func (learnSkill CGLearnSkillPacket) String() string {
	return "learn skill"
}

func readLearnSkill(buf []byte, code uint8) (packet.Packet, error) {
	skillType := binary.LittleEndian.Uint16(buf)
	return CGLearnSkillPacket{
		SkillType:       skillType,
		SkillDomainType: uint8(buf[2]),
	}, nil
}

type CGSkillToObjectPacket struct {
	SkillType      uint16
	CEffectID      uint16
	TargetObjectID uint32
}

func (skill CGSkillToObjectPacket) Id() packet.PacketID {
	return PACKET_CG_SKILL_TO_OBJECT
}

func (skill CGSkillToObjectPacket) String() string {
	return "skill to object"
}

func readSkillToObject(buf []byte, code uint8) (packet.Packet, error) {
	// encrypt!!!
	var ret CGSkillToObjectPacket
	offset := 0
	A := func() {
		ret.SkillType = binary.LittleEndian.Uint16(buf[offset:]) ^ uint16(code)
		offset += 2
	}
	B := func() {
		ret.CEffectID = binary.LittleEndian.Uint16(buf[offset:]) ^ uint16(code)
		offset += 2
	}
	C := func() {
		ret.TargetObjectID = binary.LittleEndian.Uint32(buf[offset:]) ^ uint32(code)
		offset += 4
	}
	SHUFFLE_STATEMENT_3(code, A, B, C)
	return ret, nil
}

type CGSkillToSelfPacket struct {
	SkillType uint16
	CEffectID uint16
}

func (skill CGSkillToSelfPacket) Id() packet.PacketID {
	return PACKET_CG_SKILL_TO_SELF
}

func (skill CGSkillToSelfPacket) String() string {
	return "skill to self"
}

func readSkillToSelf(buf []byte, code uint8) (packet.Packet, error) {
	// encrypt!!!
	var ret CGSkillToSelfPacket
	offset := 0
	A := func() {
		ret.SkillType = binary.LittleEndian.Uint16(buf[offset:]) ^ uint16(code)
		offset += 2
	}
	B := func() {
		ret.CEffectID = binary.LittleEndian.Uint16(buf[offset:]) ^ uint16(code)
		offset += 2
	}
	SHUFFLE_STATEMENT_2(code, A, B)
	return ret, nil
}

type CGSkillToTilePacket struct {
	SkillType uint16
	CEffectID uint16
	X         uint8
	Y         uint8
}

func (skill CGSkillToTilePacket) Id() packet.PacketID {
	return PACKET_CG_SKILL_TO_TILE
}

func (skill CGSkillToTilePacket) String() string {
	return "skill to tile"
}

func readSkillToTile(buf []byte, code uint8) (packet.Packet, error) {
	// encrypt!!!
	var ret CGSkillToTilePacket
	offset := 0
	A := func() {
		ret.SkillType = binary.LittleEndian.Uint16(buf[offset:]) ^ uint16(code)
		offset += 2
	}
	B := func() {
		ret.CEffectID = binary.LittleEndian.Uint16(buf[offset:]) ^ uint16(code)
		offset += 2
	}
	C := func() {
		ret.X = buf[offset] ^ code
		offset++
	}
	D := func() {
		ret.Y = buf[offset] ^ code
		offset++
	}
	SHUFFLE_STATEMENT_4(code, A, B, C, D)
	return ret, nil
}

type CGSayPacket struct {
	Color   uint32
	Message string
}

func (say *CGSayPacket) Id() packet.PacketID {
	return PACKET_CG_SAY
}
func (say *CGSayPacket) String() string {
	return "say"
}
func readSay(buf []byte, code uint8) (packet.Packet, error) {
	ret := new(CGSayPacket)
	ret.Color = binary.LittleEndian.Uint32(buf)
	sz := buf[2]
	ret.Message = string(buf[3 : 3+sz])
	return ret, nil
}

type CGLogoutPacket struct{}

func (_ CGLogoutPacket) Id() packet.PacketID {
	return PACKET_CG_LOGOUT
}
func (_ CGLogoutPacket) String() string {
	return "logout"
}
