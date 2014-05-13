package data

type MonsterInfo struct {
	Name         string
	Level        uint8
	STR          uint16
	DEX          uint16
	INTE         uint16
	BodySize     uint
	HP           uint16
	Exp          uint16
	MColor       uint16
	SColor       uint16
	Sight        uint8
	MeleeRange   int
	MissileRange int
}

var DeadBody MonsterInfo = MonsterInfo{
	Name:         "DeadBody",
	Level:        1,
	STR:          10,
	DEX:          10,
	INTE:         1,
	BodySize:     1,
	MColor:       1,
	SColor:       1,
	Sight:        5,
	MeleeRange:   1,
	MissileRange: 6,
}

var TurningDead MonsterInfo = MonsterInfo{
	Name:         "TurningDead",
	Level:        1,
	STR:          15,
	DEX:          10,
	INTE:         1,
	BodySize:     1,
	MColor:       1,
	SColor:       1,
	Sight:        6,
	MeleeRange:   1,
	MissileRange: 6,
}

var TurningSoul MonsterInfo = MonsterInfo{
	Name:         "TurningSoul",
	Level:        1,
	STR:          20,
	DEX:          10,
	INTE:         1,
	BodySize:     1,
	MColor:       1,
	SColor:       1,
	Sight:        7,
	MeleeRange:   1,
	MissileRange: 6,
}

var Kid MonsterInfo = MonsterInfo{
	Name:         "Kid",
	Level:        2,
	STR:          30,
	DEX:          20,
	INTE:         10,
	BodySize:     1,
	MColor:       1,
	SColor:       1,
	Sight:        8,
	MeleeRange:   1,
	MissileRange: 6,
}

var Soldier MonsterInfo = MonsterInfo{
	Name:         "Soldier",
	Level:        5,
	STR:          40,
	DEX:          25,
	INTE:         10,
	BodySize:     1,
	MColor:       1,
	SColor:       1,
	Sight:        8,
	MeleeRange:   1,
	MissileRange: 666666,
}

var JuniorWolfArch MonsterInfo = MonsterInfo{
	Name:         "Junior WolfArch",
	Level:        1,
	STR:          10,
	DEX:          10,
	INTE:         5,
	BodySize:     1,
	MColor:       1,
	SColor:       1,
	Sight:        13,
	MeleeRange:   1,
	MissileRange: 6,
}

var MonsterType2MonsterInfo map[uint16]*MonsterInfo

func init() {
	MonsterType2MonsterInfo = make(map[uint16]*MonsterInfo)
	MonsterType2MonsterInfo[4] = &DeadBody
	MonsterType2MonsterInfo[29] = &DeadBody
	MonsterType2MonsterInfo[30] = &DeadBody
	MonsterType2MonsterInfo[31] = &DeadBody
	MonsterType2MonsterInfo[32] = &DeadBody
	MonsterType2MonsterInfo[33] = &DeadBody
	MonsterType2MonsterInfo[34] = &DeadBody
	MonsterType2MonsterInfo[35] = &DeadBody
	MonsterType2MonsterInfo[36] = &DeadBody
	MonsterType2MonsterInfo[37] = &DeadBody

	MonsterType2MonsterInfo[5] = &TurningDead
	MonsterType2MonsterInfo[38] = &TurningDead
	MonsterType2MonsterInfo[39] = &TurningDead
	MonsterType2MonsterInfo[40] = &TurningDead
	MonsterType2MonsterInfo[41] = &TurningDead
	MonsterType2MonsterInfo[42] = &TurningDead
	MonsterType2MonsterInfo[43] = &TurningDead
	MonsterType2MonsterInfo[44] = &TurningDead
	MonsterType2MonsterInfo[45] = &TurningDead
	MonsterType2MonsterInfo[46] = &TurningDead

	MonsterType2MonsterInfo[7] = &TurningSoul
	MonsterType2MonsterInfo[56] = &TurningSoul
	MonsterType2MonsterInfo[57] = &TurningSoul
	MonsterType2MonsterInfo[58] = &TurningSoul
	MonsterType2MonsterInfo[59] = &TurningSoul
	MonsterType2MonsterInfo[60] = &TurningSoul
	MonsterType2MonsterInfo[61] = &TurningSoul
	MonsterType2MonsterInfo[62] = &TurningSoul
	MonsterType2MonsterInfo[63] = &TurningSoul
	MonsterType2MonsterInfo[64] = &TurningSoul

	MonsterType2MonsterInfo[6] = &Kid
	MonsterType2MonsterInfo[47] = &Kid
	MonsterType2MonsterInfo[48] = &Kid
	MonsterType2MonsterInfo[49] = &Kid
	MonsterType2MonsterInfo[50] = &Kid
	MonsterType2MonsterInfo[51] = &Kid
	MonsterType2MonsterInfo[52] = &Kid
	MonsterType2MonsterInfo[53] = &Kid
	MonsterType2MonsterInfo[54] = &Kid
	MonsterType2MonsterInfo[55] = &Kid

	MonsterType2MonsterInfo[9] = &Soldier
	MonsterType2MonsterInfo[74] = &Soldier
	MonsterType2MonsterInfo[75] = &Soldier
	MonsterType2MonsterInfo[76] = &Soldier
	MonsterType2MonsterInfo[77] = &Soldier
	MonsterType2MonsterInfo[78] = &Soldier
	MonsterType2MonsterInfo[79] = &Soldier
	MonsterType2MonsterInfo[80] = &Soldier
	MonsterType2MonsterInfo[81] = &Soldier
	MonsterType2MonsterInfo[82] = &Soldier

	MonsterType2MonsterInfo[687] = &JuniorWolfArch
}
