package darkeden

import (
	"time"
)

type GCMoveOKPacket struct {
	Dir uint8
	X   uint8
	Y   uint8
}

func (moveOk GCMoveOKPacket) Id() PacketID {
	return PACKET_GC_MOVE_OK
}
func (moveOk GCMoveOKPacket) String() string {
	return "move ok"
}
func (moveOk GCMoveOKPacket) Bytes() []byte {
	return []byte{0, moveOk.X, moveOk.Dir, moveOk.Y}
}

type GCMovePacket struct {
	ObjectID uint32
	X        uint8
	Y        uint8
	Dir      uint8
}

func (move GCMovePacket) Id() PacketID {
	return PACKET_GC_MOVE
}
func (move GCMovePacket) String() string {
	return "move"
}
func (move GCMovePacket) Bytes() []byte {
	ret := []byte{48, 0, 0, 0, 0, move.X, move.Y, move.Dir}
	binary.LittleEndian.PutUint32(ret[1:])
	return ret
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

func (updateInfo *GCUpdateInfoPacket) Id() PacketID {
	return PACKET_GC_UPDATE_INFO
}
func (updateInfo *GCUpdateInfoPacket) String() string {
	return "update info"
}
func (updateInfo *GCUpdateInfoPacket) Bytes() []byte {
	//154 1 60 1 0 0 0 86 117 48 0 0 4 183 232 191 241 150 0 0 0 164 1 0 76 29 0 0
	//20 0 20 0 20 0 20 0 20 0 20 0 20 0 20 0 20 0 216 1 216 1 50 204 41 0 0 125 0 0
	//0 0 0 0 0 26 1 0 0 13 15 39 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 4 0 0
	//0 0 100 0 0 0 0 6 118 48 0 0 30 0 0 0 232 3 0 0 0 0 1 0 0 0 0 1 0 0 0 0 0 119
	//48 0 0 44 0 0 2 16 1 136 19 0 0 0 0 3 0 0 0 0 1 0 0 0 1 0 120 48 0 0 34 5 0 0
	//1 0 0 0 0 0 255 255 255 255 0 8 0 0 0 0 1 121 48 0 0 32 0 0 2 53 43 232 3 0 0
	//0 0 4 0 0 0]
	return []byte{0, 86, 117, 48, 0, 0, 4, 183, 232, 191, 241, 150, 0, 0, 0, 164, 1, 0, 76, 29, 0, 0,
		20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 216, 1, 216, 1, 50, 204, 41, 0, 0, 125, 0, 0,
		0, 0, 0, 0, 0, 26, 1, 0, 0, 13, 15, 39, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0,
		0, 0, 100, 0, 0, 0, 0, 6, 118, 48, 0, 0, 30, 0, 0, 0, 232, 3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 119,
		48, 0, 0, 44, 0, 0, 2, 16, 1, 136, 19, 0, 0, 0, 0, 3, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 120, 48, 0, 0, 34, 5, 0, 0,
		1, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 8, 0, 0, 0, 0, 1, 121, 48, 0, 0, 32, 0, 0, 2, 53, 43, 232, 3, 0, 0,
		0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 122, 48, 0, 0, 32, 1, 0, 0, 232, 3, 0, 0, 0, 0, 2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 3, 123, 48, 0,
		0, 44, 0, 0, 2, 58, 38, 32, 28, 0, 0, 0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 4, 0, 0, 2, 146, 1, 54, 66, 109, 0, 246, 224, 0, 21, 0, 145, 237, 190, 7, 3, 19, 16,
		10, 40, 0, 0, 13, 2, 0, 5, 9, 0, 61, 0, 62, 0, 64, 0, 163, 0, 0, 0, 17, 0, 0, 0, 0, 24, 125, 0, 0, 0, 0, 0, 0, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 52, 1, 5, 0, 0, 0, 1, 0, 117, 48, 0, 0}
}

type GCSetPositionPacket struct {
	X   uint8
	Y   uint8
	Dir uint8
}

func (setPosition GCSetPositionPacket) Id() PacketID {
	return PACKET_GC_SET_POSITION
}
func (setPosition GCSetPositionPacket) String() string {
	return "set position"
}
func (setPosition GCSetPositionPacket) Bytes() []byte {
	return []byte{setPosition.Dir, setPosition.X, setPosition.Y, 2}
}

type GCAddBat struct {
	ObjectID  uint32
	Name      string
	ItemType  uint16
	X         uint8
	Y         uint8
	Dir       uint8
	CurrentHP uint16
	MaxHP     uint16
	GuildID   uint16
	Color     uint16
}

func (addBat *GCAddBat) Id() PacketID {
	return PACKET_GC_ADD_BAT
}
func (addBat *GCAddBat) String() string {
	return "add bat"
}
func (addBat *GCAddBat) Bytes() []byte {
	// [7 80 48 0 0 8 194 179 181 199 182 224 183 242 0 0 150 235 1 174 0 174 0 1 0 0 0]
	return []byte{7, 80, 48, 0, 0, 8, 194, 179, 181, 199, 182, 224, 183, 242, 0, 0, 150, 235, 1, 174, 0, 174, 0, 1, 0, 0, 0}
}

type GCAddMonsterFromBurrowing struct {
	ObjectID    uint32       // object id
	MonsterType uint16       // monster type
	MonsterName string       // monster name
	MainColor   uint16       // monster main color
	SubColor    uint16       // monster sub color
	X           uint8        // x coord.
	Y           uint8        // y coord.
	Dir         uint8        // monster direction
	EffectInfo  []EffectInfo // effects info on monster
	CurrentHP   uint16       // current hp
	MaxHP       uint16       // maximum hp
}

func (monster *GCAddMonsterFromBurrowing) Id() PacketID {
	return PACKET_GC_ADD_MONSTER_FROM_BURROWING
}
func (monster *GCAddMonsterFromBurrowing) String() string {
	return "add monster from burrowing"
}
func (monster *GCAddMonsterFromBurrowing) Bytes() []byte {
	// 185 0 27 0 0 0 48 176 47 0 0 73 0 8 200 248 182 224 210 193 182 247 36 231 240 24 148 227 4 0 156 0 156 0
	// 185 0 27 0 0 0 48 62 48 0 0 213 0 8 185 197 181 194 203 185 182 161 53 0 0 0 137 238 0 0 54 1 54 1
	return []byte{48, 62, 48, 0, 0, 213, 0, 8, 185, 197, 181, 194, 203, 185, 182, 161, 53, 0, 0, 0, 137, 238, 0, 0, 54, 1, 54, 1}
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

func (monster *GCAddMonster) Id() PacketID {
	return PACKET_GC_ADD_MONSTER
}
func (monster *GCAddMonster) String() string {
	return "add monster"
}
func (monster *GCAddMonster) Bytes() []byte {
	//[47 218 47 0 0 223 0 6 196 218 185 254 203 185 7 0 174 0 102 79 5 0 133 0 133 0 0]
	//[123 166 47 0 0 72 0 4 192 188 197 181 5 137 133 0 164 214 6 0 156 0 156 0 0]
	return
}
