package darkeden

import (
	"encoding/binary"
	"errors"
	"github.com/tiancaiamao/ouster/packet"
	// "log"
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

func readMove(buf []byte, code uint8) (packet.Packet, error) {
	dir, err := encryptDir(buf[1])
	if err != nil {
		return nil, err
	}
	ret := CGMovePacket{
		Dir: dir,
		X:   buf[2] ^ 53,
		Y:   buf[0] ^ 53,
	}
	return ret, nil
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
	ret.X = buf[0] ^ 53
	ret.Y = buf[1] ^ 53
	dir, err := encryptDir(buf[2])
	if err != nil {
		return nil, err
	}
	ret.Dir = dir
	ret.ObjectID = binary.LittleEndian.Uint32(buf[3:]) ^ 53
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
	ret.TargetObjectID = binary.LittleEndian.Uint32(buf[:]) ^ 53
	ret.SkillType = binary.LittleEndian.Uint16(buf[4:]) ^ 53
	ret.CEffectID = binary.LittleEndian.Uint16(buf[6:]) ^ 53
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
	ret.CEffectID = binary.LittleEndian.Uint16(buf) ^ 53
	ret.SkillType = binary.LittleEndian.Uint16(buf[2:]) ^ 53
	return ret, nil
}

const (
	SKILL_RAPID_GLIDING uint16 = 203
	SKILL_METEOR_STRIKE uint16 = 180
	SKILL_INVISIBILITY  uint16 = 100
	SKILL_PARALYZE      uint16 = 89
	SKILL_BLOOD_SPEAR   uint16 = 97
)

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
	ret.CEffectID = binary.LittleEndian.Uint16(buf) ^ 53
	ret.X = buf[2] ^ 53
	ret.Y = buf[3] ^ 53
	ret.SkillType = binary.LittleEndian.Uint16(buf[4:]) ^ 53
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
