package darkeden

import (
	"bytes"
	"encoding/binary"
	"github.com/tiancaiamao/ouster/packet"
	"io"
	"time"
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
func (moveOk GCMoveOKPacket) MarshalBinary() ([]byte, error) {
	return []byte{moveOk.X ^ 53, moveOk.Dir + 49, moveOk.Y ^ 53}, nil
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
func (move GCMovePacket) MarshalBinary() ([]byte, error) {
	ret := []byte{0, 0, 0, 0, move.X, move.Y, move.Dir}
	binary.LittleEndian.PutUint32(ret[:], move.ObjectID)
	return ret, nil
}

type NPCType struct{}
type InventoryInfo struct{}
type GearInfo struct{}
type ExtraInfo struct{}
type EffectInfo struct{}
type RideMotorcycleInfo struct{}
type Weather struct{}
type MonsterType struct{}
type NPCInfo struct{}
type BloodBibleSignInfo struct{}
type NicknameInfo struct{}

type GCUpdateInfoPacket struct {
	PCInfo             PCInfo
	InventoryInfo      InventoryInfo
	GearInfo           GearInfo
	ExtraInfo          ExtraInfo
	EffectInfo         EffectInfo
	hasMotorcycle      bool
	RideMotorcycleInfo RideMotorcycleInfo
	ZoneID_t           uint16
	ZoneX              uint8
	ZoneY              uint8
	GameTime           time.Time

	Weather      Weather
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
	GuildUnionID       uint
	GuildUnionUserType uint8
	BloodBibleSignInfo BloodBibleSignInfo
	PowerPoint         int
}

func (updateInfo *GCUpdateInfoPacket) Id() packet.PacketID {
	return PACKET_GC_UPDATE_INFO
}
func (updateInfo *GCUpdateInfoPacket) String() string {
	return "update info"
}
func (updateInfo *GCUpdateInfoPacket) MarshalBinary() ([]byte, error) {
	//154 1 60 1 0 0 0 86 117 48 0 0 4 183 232 191 241 150 0 0 0 164 1 0 76 29 0 0
	//20 0 20 0 20 0 20 0 20 0 20 0 20 0 20 0 20 0 216 1 216 1 50 204 41 0 0 125 0 0
	//0 0 0 0 0 26 1 0 0 13 15 39 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 4 0 0
	//0 0 100 0 0 0 0 6 118 48 0 0 30 0 0 0 232 3 0 0 0 0 1 0 0 0 0 1 0 0 0 0 0 119
	//48 0 0 44 0 0 2 16 1 136 19 0 0 0 0 3 0 0 0 0 1 0 0 0 1 0 120 48 0 0 34 5 0 0
	//1 0 0 0 0 0 255 255 255 255 0 8 0 0 0 0 1 121 48 0 0 32 0 0 2 53 43 232 3 0 0
	//0 0 4 0 0 0]

	// return []byte{86, 117, 48, 0, 0, 4, 183, 232, 191, 241, 150, 0, 0, 0, 164, 1, 0, 76, 29, 0, 0,
	// 		20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 216, 1, 216, 1, 50, 204, 41, 0, 0, 125, 0, 0,
	// 		0, 0, 0, 0, 0, 26, 1, 0, 0, 13, 15, 39, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0,
	// 		0, 0, 100, 0, 0, 0, 0, 6, 118, 48, 0, 0, 30, 0, 0, 0, 232, 3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 119,
	// 		48, 0, 0, 44, 0, 0, 2, 16, 1, 136, 19, 0, 0, 0, 0, 3, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 120, 48, 0, 0, 34, 5, 0, 0,
	// 		1, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 8, 0, 0, 0, 0, 1, 121, 48, 0, 0, 32, 0, 0, 2, 53, 43, 232, 3, 0, 0,
	// 		0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 122, 48, 0, 0, 32, 1, 0, 0, 232, 3, 0, 0, 0, 0, 2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 3, 123, 48, 0,
	// 		0, 44, 0, 0, 2, 58, 38, 32, 28, 0, 0, 0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 4, 0, 0, 2, 146, 1, 54, 66, 109, 0, 246, 224, 0, 21, 0, 145, 237, 190, 7, 3, 19, 16,
	// 		10, 40, 0, 0, 13, 2, 0, 5, 9, 0, 61, 0, 62, 0, 64, 0, 163, 0, 0, 0, 17, 0, 0, 0, 0, 24, 125, 0, 0, 0, 0, 0, 0, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 52, 1, 5, 0, 0, 0, 1, 0, 117, 48, 0, 0},
	// nil

	// return []byte{86, 117, 48, 0, 0, 4, 183, 232, 191, 241, 150, 0, 0, 0, 164, 1, 0, 76, 29, 0, 0,
	// 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 216, 1, 216, 1, 50, 204, 41, 0, 0, 125, 0, 0,
	// 0, 0, 0, 0, 0, 26, 1, 0, 0, 13, 15, 39, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0,
	// 0, 0, 100, 0, 0, 0, 0, 6, 118, 48, 0, 0, 30, 0, 0, 0, 232, 3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 119,
	// 48, 0, 0, 44, 0, 0, 2, 16, 1, 136, 19, 0, 0, 0, 0, 3, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 120, 48, 0, 0, 34, 5, 0, 0,
	// 1, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 8, 0, 0, 0, 0, 1, 121, 48, 0, 0, 32, 0, 0, 2, 53, 43, 232, 3, 0, 0,
	// 0, 0, 4, 0, 0, 0}, nil
	return []byte{86, 117, 48, 0, 0, 4, 183, 232, 191, 241, 150, 0, 0, 0, 164, 1, 0, 76, 29, 0, 0,
			20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 216, 1, 216, 1, 50, 204, 41, 0, 0, 125, 0, 0,
			0, 0, 0, 0, 0, 26, 1, 0, 0, 13, 15, 39, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0,
			0, 0, 100, 0, 0, 0, 0, 6, 118, 48, 0, 0, 30, 0, 0, 0, 232, 3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 119,
			48, 0, 0, 44, 0, 0, 2, 16, 1, 136, 19, 0, 0, 0, 0, 3, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 120, 48, 0, 0, 34, 5, 0, 0,
			1, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 8, 0, 0, 0, 0, 1, 121, 48, 0, 0, 32, 0, 0, 2, 53, 43, 232, 3, 0, 0,
			0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 122, 48, 0, 0, 32, 1, 0, 0, 232, 3, 0, 0, 0, 0, 2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 3, 123, 48, 0,
			0, 44, 0, 0, 2, 58, 38, 32, 28, 0, 0, 0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 4, 0, 0, 2, 146, 1, 54, 66, 109, 0, 246, 224, 0, 21, 0, 145, 237, 190, 7, 3, 19, 16,
			10, 40, 0, 0, 13, 2, 0, 5, 9, 0, 61, 0, 62, 0, 64, 0, 163, 0, 0, 0, 17, 0, 0, 0, 0, 24, 125, 0, 0, 0, 0, 0, 0, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		nil
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
func (pet *GCPetInfoPacket) MarshalBinary() ([]byte, error) {
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
func (setPosition GCSetPositionPacket) MarshalBinary() ([]byte, error) {
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
func (bat *GCAddBat) MarshalBinary() ([]byte, error) {

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
func (monster *GCAddMonsterFromBurrowing) MarshalBinary() ([]byte, error) {

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
func (monster *GCAddMonster) MarshalBinary() ([]byte, error) {
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
func (status GCStatusCurrentHP) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, status.ObjectID)
	binary.Write(buf, binary.LittleEndian, status.CurrentHP)
	return buf.Bytes(), nil
}

type GCAttackMeleeOK1 uint32

func (attackOk GCAttackMeleeOK1) Id() packet.PacketID {
	return PACKET_GC_ATTACK_MELEE_OK_1
}

func (attackOk GCAttackMeleeOK1) String() string {
	return "attack melee ok 1"
}
func (attackOk GCAttackMeleeOK1) MarshalBinary() ([]byte, error) {

	ret := []byte{0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(ret, uint32(attackOk))
	return ret, nil
}

type GCCannotUsePacket uint32

func (cannot GCCannotUsePacket) Id() packet.PacketID {
	return PACKET_GC_CANNOT_USE
}
func (cannot GCCannotUsePacket) String() string {
	return "cannot use"
}
func (cannot GCCannotUsePacket) MarshalBinary() ([]byte, error) {
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
func (bdo *GCBloodDrainOK1) MarshalBinary() ([]byte, error) {
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
func (modify *GCModifyInformationPacket) MarshalBinary() ([]byte, error) {
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
func (effect GCAddEffect) MarshalBinary() ([]byte, error) {
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
func (corpse *GCAddMonsterCorpse) MarshalBinary() ([]byte, error) {
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
func (died GCCreatureDiedPacket) MarshalBinary() ([]byte, error) {
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
func (obj GCDeleteObjectPacket) MarshalBinary() ([]byte, error) {
	ret := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(ret, uint32(obj))
	return ret, nil
}
