package darkeden

import (
	"time"
)

type GCMoveOKPacket struct{}

func (moveOk GCMoveOKPacket) Id() PacketID {
	return PACKET_GC_MOVE_OK
}
func (moveOk GCMoveOKPacket) String() string {
	return "move ok"
}
func (moveOK GCMoveOKPacket) Bytes() []byte {
	return []byte{}
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
