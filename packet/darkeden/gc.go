package darkeden

import (
	"bytes"
	"encoding/binary"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
	"io"
)

type GCMoveOKPacket struct {
	Dir uint8
	X   uint8
	Y   uint8
}

func (moveOk GCMoveOKPacket) Id() packet.PacketID {
	return PACKET_GC_MOVE_OK
}
func (moveOk GCMoveOKPacket) String() string {
	return "move ok"
}
func (moveOk GCMoveOKPacket) MarshalBinary(code uint8) ([]byte, error) {
	ret := make([]byte, 3)
	offset := 0
	A := func() {
		ret[offset] = moveOk.X ^ code
		offset++
	}
	B := func() {
		ret[offset] = moveOk.Y ^ code
		offset++
	}
	C := func() {
		ret[offset] = moveOk.Dir ^ code
		offset++
	}
	SHUFFLE_STATEMENT_3(code, A, B, C)
	return ret, nil
}

type GCMoveErrorPacket struct {
	X uint8
	Y uint8
}

func (moveError GCMoveErrorPacket) Id() packet.PacketID {
	return PACKET_GC_MOVE_ERROR
}
func (moveError GCMoveErrorPacket) String() string {
	return "move error"
}
func (moveError GCMoveErrorPacket) MarshalBinary(code uint8) ([]byte, error) {
	ret := make([]byte, 2)
	offset := 0
	A := func() {
		ret[offset] = moveError.X ^ code
		offset++
	}
	B := func() {
		ret[offset] = moveError.Y ^ code
		offset++
	}
	SHUFFLE_STATEMENT_2(code, A, B)
	return ret, nil
}

type GCMovePacket struct {
	ObjectID uint32
	X        uint8
	Y        uint8
	Dir      uint8
}

func (move GCMovePacket) Id() packet.PacketID {
	return PACKET_GC_MOVE
}
func (move GCMovePacket) String() string {
	return "move"
}
func (move GCMovePacket) MarshalBinary(code uint8) ([]byte, error) {
	ret := []byte{0, 0, 0, 0, move.X, move.Y, move.Dir}
	binary.LittleEndian.PutUint32(ret[:], move.ObjectID)
	return ret, nil
}

type NPCType uint16
type InventoryInfo struct{}

func (info InventoryInfo) Dump(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, uint8(0))
}

type GearInfo struct{}

func (info GearInfo) Dump(writer io.Writer) {
	// TODO
	binary.Write(writer, binary.LittleEndian, uint8(0))
}

type ExtraInfo struct{}

func (info ExtraInfo) Dump(writer io.Writer) {
	// TODO
	binary.Write(writer, binary.LittleEndian, uint8(0))
}

type EffectInfo struct{}

func (info EffectInfo) Dump(writer io.Writer) {
	// TODO
	binary.Write(writer, binary.LittleEndian, uint8(0))
}

type RideMotorcycleInfo struct{}

func (info RideMotorcycleInfo) Dump(writer io.Writer) {
	// TODO
	binary.Write(writer, binary.LittleEndian, uint8(0))
}

type Weather struct{}
type MonsterType uint16
type NPCInfo struct{}

func (info NPCInfo) Dump(writer io.Writer) {
	// TODO
}

type BloodBibleSignInfo struct{}

func (info BloodBibleSignInfo) Dump(writer io.Writer) {
	// TODO
	binary.Write(writer, binary.LittleEndian, uint8(0))
	binary.Write(writer, binary.LittleEndian, uint8(0))
	return
}

type NicknameInfo struct {
	NicknameID    uint16
	NicknameType  uint8
	Nickname      string
	NicknameIndex uint16
}

func (info NicknameInfo) Dump(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, info.NicknameID)
	binary.Write(writer, binary.LittleEndian, uint8(0))
	return
}

type GameTimeType struct {
	Year  uint16
	Month uint8
	Day   uint8

	Hour   uint8
	Minute uint8
	Second uint8
}

func (time GameTimeType) Dump(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, time.Year)
	binary.Write(writer, binary.LittleEndian, time.Month)
	binary.Write(writer, binary.LittleEndian, time.Day)
	binary.Write(writer, binary.LittleEndian, time.Hour)
	binary.Write(writer, binary.LittleEndian, time.Minute)
	binary.Write(writer, binary.LittleEndian, time.Second)
	return
}

type GCUpdateInfoPacket struct {
	PCType             byte
	PCInfo             data.PCInfo
	InventoryInfo      InventoryInfo
	GearInfo           GearInfo
	ExtraInfo          ExtraInfo
	EffectInfo         EffectInfo
	hasMotorcycle      bool
	RideMotorcycleInfo RideMotorcycleInfo

	ZoneID   uint16
	ZoneX    uint8
	ZoneY    uint8
	GameTime GameTimeType

	Weather      uint8
	WeatherLevel uint8

	DarkLevel  uint8
	LightLevel uint8

	NPCNum   uint8
	NPCTypes [256]NPCType

	MonsterNum   uint8
	MonsterTypes [256]MonsterType

	NPCInfos []NPCInfo

	ServerStat   uint8
	Premium      uint8
	SMSCharge    uint32
	NicknameInfo NicknameInfo

	NonPK              bool
	GuildUnionID       uint32
	GuildUnionUserType uint8
	BloodBibleSignInfo BloodBibleSignInfo
	PowerPoint         int
}

func (info *GCUpdateInfoPacket) Id() packet.PacketID {
	return PACKET_GC_UPDATE_INFO
}
func (info *GCUpdateInfoPacket) String() string {
	return "update info"
}
func (info *GCUpdateInfoPacket) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteByte(info.PCType)
	info.PCInfo.Dump(buf)

	info.InventoryInfo.Dump(buf)
	info.GearInfo.Dump(buf)
	info.ExtraInfo.Dump(buf)
	info.EffectInfo.Dump(buf)
	if info.hasMotorcycle {
		buf.WriteByte(1)
		info.RideMotorcycleInfo.Dump(buf)
	} else {
		buf.WriteByte(0)
	}

	// write zone info
	binary.Write(buf, binary.LittleEndian, info.ZoneID)
	binary.Write(buf, binary.LittleEndian, info.ZoneX)
	binary.Write(buf, binary.LittleEndian, info.ZoneY)

	// info.GameTime.Dump(buf)
	// 	binary.Write(buf, binary.LittleEndian, info.Weather)
	// 	binary.Write(buf, binary.LittleEndian, info.WeatherLevel)
	// 	binary.Write(buf, binary.LittleEndian, info.DarkLevel)
	// 	binary.Write(buf, binary.LittleEndian, info.LightLevel)
	//
	// 	binary.Write(buf, binary.LittleEndian, info.NPCNum)
	// 	for i:=0; i<int(info.NPCNum); i++ {
	// 		binary.Write(buf, binary.LittleEndian, info.NPCTypes[i])
	// 	}
	//
	// 	binary.Write(buf, binary.LittleEndian, info.MonsterNum)
	// 	for i:=0; i<int(info.MonsterNum); i++ {
	// 		binary.Write(buf, binary.LittleEndian, info.MonsterTypes[i])
	// 	}
	//
	// 	binary.Write(buf, binary.LittleEndian, uint8(len(info.NPCInfos)))
	// 	for i:=0; i<len(info.NPCInfos); i++ {
	// 		info.NPCInfos[i].Dump(buf)
	// 	}
	//
	// 	binary.Write(buf, binary.LittleEndian, info.ServerStat)
	// 	binary.Write(buf, binary.LittleEndian, info.Premium)
	// 	binary.Write(buf, binary.LittleEndian, info.SMSCharge)
	//
	// 	info.NicknameInfo.Dump(buf)
	//
	// 	if info.NonPK {
	// 		binary.Write(buf, binary.LittleEndian, uint8(1))
	// 	} else {
	// 		binary.Write(buf, binary.LittleEndian, uint8(0))
	// 	}
	//
	// 	binary.Write(buf, binary.LittleEndian, info.GuildUnionID)
	// 	binary.Write(buf, binary.LittleEndian, info.GuildUnionUserType)
	//
	// 	info.BloodBibleSignInfo.Dump(buf)
	//
	// 	binary.Write(buf, binary.LittleEndian, info.PowerPoint)
	//
	// 	return buf.Bytes(), nil

	buf.Write([]byte{190, 7, 3, 19, 16,
		10, 40, 0, 0, 13, 2, 0, 5, 9, 0, 61, 0, 62, 0, 64, 0, 163, 0, 0, 0, 17, 0, 0, 0, 0, 24, 125, 0, 0, 0, 0, 0, 0, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0})
	return buf.Bytes(), nil

	// return []byte{86, 117, 48, 0, 0, 4, 183, 232, 191, 241, 150, 0, 0, 0, 164, 1, 0, 76, 29, 0, 0,
	// 		20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 216, 1, 216, 1, 50, 204, 41, 0, 0, 125, 0, 0,
	// 		0, 0, 0, 0, 0, 26, 1, 0, 0, 13, 15, 39, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0,
	// 		0, 0, 100, 0, 0, 0, 0, 6, 118, 48, 0, 0, 30, 0, 0, 0, 232, 3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 119,
	// 		48, 0, 0, 44, 0, 0, 2, 16, 1, 136, 19, 0, 0, 0, 0, 3, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 120, 48, 0, 0, 34, 5, 0, 0,
	// 		1, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 8, 0, 0, 0, 0, 1, 121, 48, 0, 0, 32, 0, 0, 2, 53, 43, 232, 3, 0, 0,
	// 		0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 122, 48, 0, 0, 32, 1, 0, 0, 232, 3, 0, 0, 0, 0, 2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 3, 123, 48, 0,
	// 		0, 44, 0, 0, 2, 58, 38, 32, 28, 0, 0, 0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 4, 0, 0, 2, 146, 1, 54, 66, 109, 0, 246, 224, 0, 21, 0, 145, 237, 190, 7, 3, 19, 16,
	// 		10, 40, 0, 0, 13, 2, 0, 5, 9, 0, 61, 0, 62, 0, 64, 0, 163, 0, 0, 0, 17, 0, 0, 0, 0, 24, 125, 0, 0, 0, 0, 0, 0, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	nil
}

type GCPetInfoPacket struct {
	ObjectID   uint32
	PetInfo    []struct{}
	SummonInfo uint8
}

func (pet *GCPetInfoPacket) Id() packet.PacketID {
	return PACKET_GC_PET_INFO
}
func (pet *GCPetInfoPacket) String() string {
	return "pet info"
}
func (pet *GCPetInfoPacket) MarshalBinary(code uint8) ([]byte, error) {
	return []byte{0, 117, 48, 0, 0}, nil
}

type GCSetPositionPacket struct {
	X   uint8
	Y   uint8
	Dir uint8
}

func (setPosition GCSetPositionPacket) Id() packet.PacketID {
	return PACKET_GC_SET_POSITION
}
func (setPosition GCSetPositionPacket) String() string {
	return "set position"
}
func (setPosition GCSetPositionPacket) MarshalBinary(code uint8) ([]byte, error) {
	return []byte{setPosition.X, setPosition.Y, setPosition.Dir}, nil
}

type GCAddBat struct {
	ObjectID    uint32
	MonsterName string
	ItemType    uint16
	X           uint8
	Y           uint8
	Dir         uint8
	CurrentHP   uint16
	MaxHP       uint16
	GuildID     uint16
	Color       uint16
}

func (bat *GCAddBat) Id() packet.PacketID {
	return PACKET_GC_ADD_BAT
}
func (bat *GCAddBat) String() string {
	return "add bat"
}
func (bat *GCAddBat) MarshalBinary(code uint8) ([]byte, error) {

	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, bat.ObjectID)
	binary.Write(buf, binary.LittleEndian, uint8(len(bat.MonsterName)))
	io.WriteString(buf, bat.MonsterName)
	binary.Write(buf, binary.LittleEndian, bat.ItemType)
	binary.Write(buf, binary.LittleEndian, bat.X)
	binary.Write(buf, binary.LittleEndian, bat.Y)
	binary.Write(buf, binary.LittleEndian, bat.Dir)
	binary.Write(buf, binary.LittleEndian, bat.CurrentHP)
	binary.Write(buf, binary.LittleEndian, bat.MaxHP)
	binary.Write(buf, binary.LittleEndian, bat.GuildID)
	binary.Write(buf, binary.LittleEndian, bat.Color)
	return buf.Bytes(), nil
}

type GCAddMonsterFromBurrowing struct {
	ObjectID    uint32
	MonsterType uint16
	MonsterName string
	MainColor   uint16
	SubColor    uint16
	X           uint8
	Y           uint8
	Dir         uint8
	EffectInfo  []EffectInfo
	CurrentHP   uint16
	MaxHP       uint16
}

func (monster *GCAddMonsterFromBurrowing) Id() packet.PacketID {
	return PACKET_GC_ADD_MONSTER_FROM_BURROWING
}
func (monster *GCAddMonsterFromBurrowing) String() string {
	return "add monster from burrowing"
}
func (monster *GCAddMonsterFromBurrowing) MarshalBinary(code uint8) ([]byte, error) {

	return []byte{62, 48, 0, 0, 213, 0, 8, 185, 197, 181, 194, 203, 185, 182, 161, 53, 0, 0, 0, 137, 238, 0, 0, 54, 1, 54, 1}, nil
}

type GCAddMonster struct {
	ObjectID    uint32
	MonsterType uint16
	MonsterName string
	MainColor   uint16
	SubColor    uint16
	X           uint8
	Y           uint8
	Dir         uint8
	EffectInfo  []EffectInfo
	CurrentHP   uint16
	MaxHP       uint16
	FromFlag    byte
}

func (monster *GCAddMonster) Id() packet.PacketID {
	return PACKET_GC_ADD_MONSTER
}
func (monster *GCAddMonster) String() string {
	return "add monster"
}
func (monster *GCAddMonster) MarshalBinary(code uint8) ([]byte, error) {
	//[218 47 0 0 223 0 6 196 218 185 254 203 185 7 0 174 0 102 79 5 0 133 0 133 0 0]
	//[166 47 0 0 72 0 4 192 188 197 181 5 137 133 0 164 214 6 0 156 0 156 0 0]
	//[24 47 0 0 8 0 10 203 185 196 170 191 203 206 172 198 230 53 48 48 58 137 192 6 0 156 0 156 0 0]
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, monster.ObjectID)
	binary.Write(buf, binary.LittleEndian, monster.MonsterType)
	binary.Write(buf, binary.LittleEndian, uint8(len(monster.MonsterName)))
	io.WriteString(buf, monster.MonsterName)
	binary.Write(buf, binary.LittleEndian, monster.MainColor)
	binary.Write(buf, binary.LittleEndian, monster.SubColor)
	binary.Write(buf, binary.LittleEndian, monster.X)
	binary.Write(buf, binary.LittleEndian, monster.Y)
	binary.Write(buf, binary.LittleEndian, monster.Dir)
	binary.Write(buf, binary.LittleEndian, uint8(0))
	binary.Write(buf, binary.LittleEndian, monster.CurrentHP)
	binary.Write(buf, binary.LittleEndian, monster.MaxHP)
	binary.Write(buf, binary.LittleEndian, monster.FromFlag)
	return buf.Bytes(), nil
}

type GCStatusCurrentHP struct {
	ObjectID  uint32
	CurrentHP uint16
}

func (status GCStatusCurrentHP) Id() packet.PacketID {
	return PACKET_GC_STATUS_CURRENT_HP
}
func (status GCStatusCurrentHP) String() string {
	return "status current HP"
}
func (status GCStatusCurrentHP) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, status.ObjectID)
	binary.Write(buf, binary.LittleEndian, status.CurrentHP)
	return buf.Bytes(), nil
}

type GCAttackMeleeOK1 struct {
	ObjectID uint32
	ModifyInfo
}

func (attackOk GCAttackMeleeOK1) Id() packet.PacketID {
	return PACKET_GC_ATTACK_MELEE_OK_1
}

func (attackOk GCAttackMeleeOK1) String() string {
	return "attack melee ok 1"
}
func (attackOk GCAttackMeleeOK1) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, attackOk.ObjectID)
	attackOk.Dump(buf)
	return buf.Bytes(), nil
}

type GCAttackMeleeOK2 struct {
	ObjectID uint32
	ModifyInfo
}

func (attackOk GCAttackMeleeOK2) Id() packet.PacketID {
	return PACKET_GC_ATTACK_MELEE_OK_2
}

func (attackOk GCAttackMeleeOK2) String() string {
	return "attack melee ok 2"
}
func (attackOk GCAttackMeleeOK2) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, attackOk.ObjectID)
	attackOk.Dump(buf)
	return buf.Bytes(), nil
}

type GCCannotUsePacket uint32

func (cannot GCCannotUsePacket) Id() packet.PacketID {
	return PACKET_GC_CANNOT_USE
}
func (cannot GCCannotUsePacket) String() string {
	return "cannot use"
}
func (cannot GCCannotUsePacket) MarshalBinary(code uint8) ([]byte, error) {
	ret := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(ret, uint32(cannot))
	return ret, nil
}

type ModifyType byte

const (
	MODIFY_BASIC_STR ModifyType = iota
	MODIFY_CURRENT_STR
	MODIFY_MAX_STR
	MODIFY_STR_EXP
	MODIFY_BASIC_DEX
	MODIFY_CURRENT_DEX
	MODIFY_MAX_DEX
	MODIFY_DEX_EXP
	MODIFY_BASIC_INT
	MODIFY_CURRENT_INT
	MODIFY_MAX_INT
	MODIFY_INT_EXP
	MODIFY_CURRENT_HP
	MODIFY_MAX_HP
	MODIFY_CURRENT_MP
	MODIFY_MAX_MP
	MODIFY_MIN_DAMAGE
	MODIFY_MAX_DAMAGE
	MODIFY_DEFENSE
	MODIFY_PROTECTION
	MODIFY_TOHIT
	MODIFY_VISION
	MODIFY_FAME
	MODIFY_GOLD
	MODIFY_SWORD_DOMAIN_LEVEL
	MODIFY_SWORD_DOMAIN_EXP
	MODIFY_SWORD_DOMAIN_GOAL_EXP
	MODIFY_BLADE_DOMAIN_LEVEL
	MODIFY_BLADE_DOMAIN_EXP
	MODIFY_BLADE_DOMAIN_GOAL_EXP
	MODIFY_HEAL_DOMAIN_LEVEL
	MODIFY_HEAL_DOMAIN_EXP
	MODIFY_HEAL_DOMAIN_GOAL_EXP
	MODIFY_ENCHANT_DOMAIN_LEVEL
	MODIFY_ENCHANT_DOMAIN_EXP
	MODIFY_ENCHANT_DOMAIN_GOAL_EXP
	MODIFY_GUN_DOMAIN_LEVEL
	MODIFY_GUN_DOMAIN_EXP
	MODIFY_GUN_DOMAIN_GOAL_EXP
	MODIFY_ETC_DOMAIN_LEVEL
	MODIFY_ETC_DOMAIN_EXP
	MODIFY_ETC_DOMAIN_GOAL_EXP
	MODIFY_SKILL_LEVEL
	MODIFY_LEVEL
	MODIFY_EFFECT_STAT
	MODIFY_DURATION
	MODIFY_BULLET
	MODIFY_BONUS_POINT
	MODIFY_DURABILITY
	MODIFY_NOTORIETY
	MODIFY_VAMP_GOAL_EXP
	MODIFY_SILVER_DAMAGE
	MODIFY_ATTACK_SPEED
	MODIFY_ALIGNMENT
	MODIFY_SILVER_DURABILITY
	MODIFY_REGEN_RATE
	MODIFY_GUILDID
	MODIFY_RANK
	MODIFY_RANK_EXP
	MODIFY_OUSTERS_GOAL_EXP
	MODIFY_SKILL_BONUS_POINT
	MODIFY_ELEMENTAL_FIRE
	MODIFY_ELEMENTAL_WATER
	MODIFY_ELEMENTAL_EARTH
	MODIFY_ELEMENTAL_WIND
	MODIFY_SKILL_EXP
	MODIFY_PET_HP
	MODIFY_PET_EXP
	MODIFY_LAST_TARGET
	MODIFY_UNIONID
	MODIFY_UNIONGRADE
	MODIFY_ADVANCEMENT_CLASS_LEVEL
	MODIFY_ADVANCEMENT_CLASS_GOAL_EXP
	MODIFY_MAX
)

var ModifyType2String []string = []string{
	"BASIC_STR",
	"CURRENT_STR",
	"MAX_STR",
	"STR_EXP",
	"BASIC_DEX",
	"CURRENT_DEX",
	"MAX_DEX",
	"DEX_EXP",
	"BASIC_INT",
	"CURRENT_INT",
	"MAX_INT",
	"INT_EXP",
	"CURRENT_HP",
	"MAX_HP",
	"CURRENT_MP",
	"MAX_MP",
	"MIN_DAMAGE",
	"MAX_DAMAGE",
	"DEFENSE",
	"PROTECTION",
	"TOHIT",
	"VISION",
	"FAME",
	"GOLD",
	"SWORD_DOMAIN_LEVEL",
	"SWORD_DOMAIN_EXP",
	"SWORD_DOMAIN_GOAL_EXP",
	"BLADE_DOMAIN_LEVEL",
	"BLADE_DOMAIN_EXP",
	"BLADE_DOMAIN_GOAL_EXP",
	"HEAL_DOMAIN_LEVEL",
	"HEAL_DOMAIN_EXP",
	"HEAL_DOMAIN_GOAL_EXP",
	"ENCHANT_DOMAIN_LEVEL",
	"ENCHANT_DOMAIN_EXP",
	"ENCHANT_DOMAIN_GOAL_EXP",
	"GUN_DOMAIN_LEVEL",
	"GUN_DOMAIN_EXP",
	"GUN_DOMAIN_GOAL_EXP",
	"ETC_DOMAIN_LEVEL",
	"ETC_DOMAIN_EXP",
	"ETC_DOMAIN_GOAL_EXP",
	"SKILL_LEVEL",
	"LEVEL",
	"EFFECT_STAT",
	"DURATION",
	"BULLET",
	"BONUS_POINT",
	"DURABILITY",
	"NOTORIETY",
	"VAMP_EXP",
	"SILVER_DAMAGE",
	"ATTACK_SPEED",
	"ALIGNMENT",
	"SILVER_DURABILITY",
	"REGEN_RATE",
	"GUILDID",
	"RANK",
	"RANK_EXP",
	"MODIFY_OUSTERS_EXP",
	"MODIFY_SKILL_BONUS_POINT",
	"MODIFY_ELEMENTAL_FIRE",
	"MODIFY_ELEMENTAL_WATER",
	"MODIFY_ELEMENTAL_EARTH",
	"MODIFY_ELEMENTAL_WIND",
	"MODIFY_SKILL_EXP",
	"MODIFY_PET_HP",
	"MODIFY_PET_EXP",
	"MODIFY_LAST_TARGET",
	"MODIFY_UNIONID",
	"MODIFY_UNIONGRADE",
	"MODIFY_ADVANCEMENT_CLASS_LEVEL",
	"MODIFY_ADVANCEMENT_CLASS_GOAL_EXP",
	"MAX",
}

type ModifyInfo struct {
	Short map[ModifyType]uint16
	Long  map[ModifyType]uint32
}

func (modify *ModifyInfo) Dump(writer io.Writer) {
	szShort := uint8(len(modify.Short))
	binary.Write(writer, binary.LittleEndian, szShort)
	for k, v := range modify.Short {
		binary.Write(writer, binary.LittleEndian, k)
		binary.Write(writer, binary.LittleEndian, v)
	}

	szLong := uint8(len(modify.Long))
	binary.Write(writer, binary.LittleEndian, szLong)
	for k, v := range modify.Long {
		binary.Write(writer, binary.LittleEndian, k)
		binary.Write(writer, binary.LittleEndian, v)
	}
}

type GCBloodDrainOK1 struct {
	Modify   ModifyInfo
	ObjectID uint32
}

func (bdo *GCBloodDrainOK1) Id() packet.PacketID {
	return PACKET_GC_BLOOD_DRAIN_OK_1
}
func (bdo *GCBloodDrainOK1) String() string {
	return "blood drain ok 1"
}
func (bdo *GCBloodDrainOK1) MarshalBinary(code uint8) ([]byte, error) {
	// 237, 53, 0, 0, 2, 51, 0, 0, 12, 216, 1, 0}
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, bdo.ObjectID)
	bdo.Modify.Dump(buf)
	return buf.Bytes(), nil
}

type GCModifyInformationPacket ModifyInfo

func (modify *GCModifyInformationPacket) Id() packet.PacketID {
	return PACKET_GC_MODIFY_INFORMATION
}
func (modify *GCModifyInformationPacket) String() string {
	return "modify information"
}
func (modify *GCModifyInformationPacket) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	raw := (*ModifyInfo)(modify)
	raw.Dump(buf)
	return buf.Bytes(), nil
}

type GCAddEffect struct {
	ObjectID uint32
	EffectID uint16
	Duration uint16
}

func (effect GCAddEffect) Id() packet.PacketID {
	return PACKET_GC_ADD_EFFECT
}
func (effect GCAddEffect) String() string {
	return "add effect"
}
func (effect GCAddEffect) MarshalBinary(code uint8) ([]byte, error) {
	ret := make([]byte, 8)
	binary.LittleEndian.PutUint32(ret, effect.ObjectID)
	binary.LittleEndian.PutUint16(ret[4:], effect.EffectID)
	binary.LittleEndian.PutUint16(ret[6:], effect.Duration)
	return ret, nil
}

type GCAddMonsterCorpse struct {
	ObjectID    uint32
	MonsterType uint16
	MonsterName string

	X       uint8
	Y       uint8
	Dir     uint8
	HasHead bool

	TreasureCount uint8
	LastKiller    uint32
}

func (corpse *GCAddMonsterCorpse) Id() packet.PacketID {
	return PACKET_GC_ADD_MONSTER_CORPSE
}
func (corpse *GCAddMonsterCorpse) String() string {
	return "add monster corpse"
}
func (corpse *GCAddMonsterCorpse) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, corpse.ObjectID)
	binary.Write(buf, binary.LittleEndian, corpse.MonsterType)
	binary.Write(buf, binary.LittleEndian, uint8(len(corpse.MonsterName)))
	if len(corpse.MonsterName) > 0 {
		io.WriteString(buf, corpse.MonsterName)
	}
	binary.Write(buf, binary.LittleEndian, corpse.X)
	binary.Write(buf, binary.LittleEndian, corpse.Y)
	binary.Write(buf, binary.LittleEndian, corpse.Dir)
	if corpse.HasHead {
		binary.Write(buf, binary.LittleEndian, uint8(1))
	} else {
		binary.Write(buf, binary.LittleEndian, uint8(0))
	}

	binary.Write(buf, binary.LittleEndian, corpse.TreasureCount)
	binary.Write(buf, binary.LittleEndian, corpse.LastKiller)
	return buf.Bytes(), nil
}

type GCCreatureDiedPacket uint32

func (died GCCreatureDiedPacket) Id() packet.PacketID {
	return PACKET_GC_CREATURE_DIED
}
func (died GCCreatureDiedPacket) String() string {
	return "creature died"
}
func (died GCCreatureDiedPacket) MarshalBinary(code uint8) ([]byte, error) {
	ret := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(ret, uint32(died))
	return ret, nil
}

type GCDeleteObjectPacket uint32

func (obj GCDeleteObjectPacket) Id() packet.PacketID {
	return PACKET_GC_DELETE_OBJECT
}
func (obj GCDeleteObjectPacket) String() string {
	return "delete object"
}
func (obj GCDeleteObjectPacket) MarshalBinary(code uint8) ([]byte, error) {
	ret := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(ret, uint32(obj))
	return ret, nil
}

type GCAddEffectPacket struct {
	ObjectID uint32
	EffectID uint16
	Duration uint16
}

func (obj GCAddEffectPacket) Id() packet.PacketID {
	return PACKET_GC_ADD_EFFECT
}
func (obj GCAddEffectPacket) String() string {
	return "add effect"
}
func (obj GCAddEffectPacket) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, obj.ObjectID)
	binary.Write(buf, binary.LittleEndian, obj.EffectID)
	binary.Write(buf, binary.LittleEndian, obj.Duration)
	return buf.Bytes(), nil
}

type GCFastMovePacket struct {
	ObjectID  uint32
	FromX     uint8
	FromY     uint8
	ToX       uint8
	ToY       uint8
	SkillType uint16
}

func (fastMove *GCFastMovePacket) Id() packet.PacketID {
	return PACKET_GC_FAST_MOVE
}
func (fastMove *GCFastMovePacket) String() string {
	return "fast move"
}
func (fastMove *GCFastMovePacket) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, fastMove.ObjectID)
	binary.Write(buf, binary.LittleEndian, fastMove.FromX)
	binary.Write(buf, binary.LittleEndian, fastMove.FromY)
	binary.Write(buf, binary.LittleEndian, fastMove.ToX)
	binary.Write(buf, binary.LittleEndian, fastMove.ToY)
	binary.Write(buf, binary.LittleEndian, fastMove.SkillType)
	return buf.Bytes(), nil
}

type GCLearnSkillOK struct {
	SkillType       uint16
	SkillDomainType uint8
}

func (ok *GCLearnSkillOK) Id() packet.PacketID {
	return PACKET_GC_LEARN_SKILL_OK
}
func (ok *GCLearnSkillOK) String() string {
	return "learn skill ok"
}
func (ok *GCLearnSkillOK) MarshalBinary(code uint8) ([]byte, error) {
	ret := []byte{0, 0, ok.SkillDomainType}
	binary.LittleEndian.PutUint16(ret, ok.SkillType)
	return ret, nil
}

type GCRemoveEffect struct {
	ObjectID   uint32
	EffectList []uint16
}

func (remove GCRemoveEffect) Id() packet.PacketID {
	return PACKET_GC_REMOVE_EFFECT
}
func (remove GCRemoveEffect) String() string {
	return "remove effect"
}
func (remove GCRemoveEffect) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, remove.ObjectID)
	binary.Write(buf, binary.LittleEndian, uint8(len(remove.EffectList)))
	for _, v := range remove.EffectList {
		binary.Write(buf, binary.LittleEndian, v)
	}

	return buf.Bytes(), nil
}

type GCSkillFailed1 struct {
	SkillType uint16
	Grade     uint8
	ModifyInfo
}

func (failed *GCSkillFailed1) Id() packet.PacketID {
	return PACKET_GC_SKILL_FAILED_1
}
func (failed *GCSkillFailed1) String() string {
	return "skill failed1"
}
func (failed *GCSkillFailed1) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, failed.SkillType)
	binary.Write(buf, binary.LittleEndian, failed.Grade)
	failed.Dump(buf)
	return buf.Bytes(), nil
}

type GCSkillFailed2 struct {
	ObjectID       uint32
	TargetObjectID uint32
	SkillType      uint16
	Grade          uint8
}

func (failed GCSkillFailed2) Id() packet.PacketID {
	return PACKET_GC_SKILL_FAILED_2
}
func (failed GCSkillFailed2) String() string {
	return "skill failed1"
}
func (failed GCSkillFailed2) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, failed.ObjectID)
	binary.Write(buf, binary.LittleEndian, failed.TargetObjectID)
	binary.Write(buf, binary.LittleEndian, failed.SkillType)
	binary.Write(buf, binary.LittleEndian, failed.Grade)
	return buf.Bytes(), nil
}

type GCSkillToObjectOK1 struct {
	SkillType      uint16
	CEffectID      uint16
	TargetObjectID uint32
	Duration       uint16
	Grade          uint8
	ModifyInfo
}

func (ok *GCSkillToObjectOK1) Id() packet.PacketID {
	return PACKET_GC_SKILL_TO_OBJECT_OK_1
}
func (ok *GCSkillToObjectOK1) String() string {
	return "skill to object ok 1"
}
func (ok *GCSkillToObjectOK1) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, ok.SkillType)
	binary.Write(buf, binary.LittleEndian, ok.CEffectID)
	binary.Write(buf, binary.LittleEndian, ok.TargetObjectID)
	binary.Write(buf, binary.LittleEndian, ok.Duration)
	binary.Write(buf, binary.LittleEndian, ok.Grade)
	ok.Dump(buf)
	return buf.Bytes(), nil
}

type GCSkillToObjectOK4 struct {
	ObjectID  uint32
	SkillType uint16
	Duration  uint16
	Grade     uint8
}

func (ok GCSkillToObjectOK4) Id() packet.PacketID {
	return PACKET_GC_SKILL_TO_OBJECT_OK_4
}
func (ok GCSkillToObjectOK4) String() string {
	return "skill to object ok 4"
}
func (ok GCSkillToObjectOK4) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, ok.ObjectID)
	binary.Write(buf, binary.LittleEndian, ok.SkillType)
	binary.Write(buf, binary.LittleEndian, ok.Duration)
	binary.Write(buf, binary.LittleEndian, ok.Grade)
	return buf.Bytes(), nil
}

type GCSkillToSelfOK1 struct {
	SkillType uint16
	CEffectID uint16
	Duration  uint16
	Grade     uint8
	ModifyInfo
}

func (ok *GCSkillToSelfOK1) Id() packet.PacketID {
	return PACKET_GC_SKILL_TO_SELF_OK_1
}
func (ok *GCSkillToSelfOK1) String() string {
	return "skill to self ok 1"
}
func (ok *GCSkillToSelfOK1) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, ok.SkillType)
	binary.Write(buf, binary.LittleEndian, ok.CEffectID)
	binary.Write(buf, binary.LittleEndian, ok.Duration)
	binary.Write(buf, binary.LittleEndian, ok.Grade)
	ok.Dump(buf)
	return buf.Bytes(), nil
}

type GCSkillToTileOK1 struct {
	SkillType    uint16
	CEffectID    uint16
	Duration     uint16
	Range        uint8
	X            uint8
	Y            uint8
	CreatureList []uint32
	Grade        uint8
	ModifyInfo
}

func (ok *GCSkillToTileOK1) Id() packet.PacketID {
	return PACKET_GC_SKILL_TO_TILE_OK_1
}
func (ok *GCSkillToTileOK1) String() string {
	return "skill to tile ok 1"
}
func (ok *GCSkillToTileOK1) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, ok.SkillType)
	binary.Write(buf, binary.LittleEndian, ok.CEffectID)
	binary.Write(buf, binary.LittleEndian, ok.X)
	binary.Write(buf, binary.LittleEndian, ok.Y)
	binary.Write(buf, binary.LittleEndian, ok.Duration)
	binary.Write(buf, binary.LittleEndian, ok.Range)
	binary.Write(buf, binary.LittleEndian, ok.Grade)
	binary.Write(buf, binary.LittleEndian, uint8(len(ok.CreatureList)))
	for _, v := range ok.CreatureList {
		binary.Write(buf, binary.LittleEndian, v)
	}
	ok.Dump(buf)

	return buf.Bytes(), nil
}

const (
	SYSTEM_MESSAGE_NORMAL = iota
	SYSTEM_MESSAGE_OPERATOR
	SYSTEM_MESSAGE_MASTER_LAIR
	SYSTEM_MESSAGE_COMBAT
	SYSTEM_MESSAGE_INFO
	SYSTEM_MESSAGE_HOLY_LAND
	SYSTEM_MESSAGE_RANGER_SAY
	SYSTEM_MESSAGE_MAX
)

type GCSystemMessagePacket struct {
	Message string
	Color   uint32
	Type    uint8
}

func (msg *GCSystemMessagePacket) Id() packet.PacketID {
	return PACKET_GC_SYSTEM_MESSAGE
}
func (msg *GCSystemMessagePacket) String() string {
	return "system message"
}
func (msg *GCSystemMessagePacket) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, uint8(len(msg.Message)))
	io.WriteString(buf, msg.Message)
	binary.Write(buf, binary.LittleEndian, msg.Color)
	binary.Write(buf, binary.LittleEndian, msg.Type)
	return buf.Bytes(), nil
}

const (
	PC_SLAYER PCType = iota
	PC_VAMPIRE
	PC_OUSTER
)

type GCSkillInfoPacket struct {
	PCType          PCType
	PCSkillInfoList []SkillInfo
}

func (info *GCSkillInfoPacket) Id() packet.PacketID {
	return PACKET_GC_SKILL_INFO
}
func (info *GCSkillInfoPacket) String() string {
	return "skill info"
}
func (info *GCSkillInfoPacket) MarshalBinary(code uint8) ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, uint8(info.PCType))
	binary.Write(buf, binary.LittleEndian, uint8(len(info.PCSkillInfoList)))
	for _, v := range info.PCSkillInfoList {
		v.Dump(buf)
	}
	return buf.Bytes(), nil
}

type SkillInfo interface {
	Dump(io.Writer)
}
type VampireSkillInfo struct {
	LearnNewSkill           bool
	SubVampireSkillInfoList []SubVampireSkillInfo
}

func (info VampireSkillInfo) Dump(writer io.Writer) {
	if info.LearnNewSkill {
		binary.Write(writer, binary.LittleEndian, uint8(1))
	} else {
		binary.Write(writer, binary.LittleEndian, uint8(0))
	}

	binary.Write(writer, binary.LittleEndian, uint8(len(info.SubVampireSkillInfoList)))
	for _, v := range info.SubVampireSkillInfoList {
		v.Dump(writer)
	}
}

type SubVampireSkillInfo struct {
	SkillType   uint16
	Interval    uint32
	CastingTime uint32
}

func (info SubVampireSkillInfo) Dump(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, info.SkillType)
	binary.Write(writer, binary.LittleEndian, info.Interval)
	binary.Write(writer, binary.LittleEndian, info.CastingTime)
	return
}
