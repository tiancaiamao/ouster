package main

import (
	"github.com/tiancaiamao/ouster/packet/darkeden"
)

type SkillPropertyType uint8

const (
	SKILL_PROPERTY_TYPE_MELEE SkillPropertyType = iota
	SKILL_PROPERTY_TYPE_MAGIC
	SKILL_PROPERTY_TYPE_PHYSIC
)

const (
	SKILL_RAPID_GLIDING uint16 = 203
	SKILL_METEOR_STRIKE uint16 = 180
	SKILL_INVISIBILITY  uint16 = 100
	SKILL_PARALYZE      uint16 = 89
	SKILL_BLOOD_SPEAR   uint16 = 97
)

type SkillInfo struct {
	Type        SkillPropertyType
	Name        string
	Level       uint
	MinDamage   uint
	MaxDamage   uint
	MinDelay    uint
	MinCastTime uint
	MaxCastTime uint

	ConsumeMP uint

	P2M func(*Player, *Monster)
	M2P func(*Monster, *Player)
	P2P func(*Player, *Player)
}

var skillTable map[uint16]*SkillInfo

func init() {
	skillTable = make(map[uint16]*SkillInfo)

	skillTable[SKILL_BLOOD_SPEAR] = &SkillInfo{
		Type:      SKILL_PROPERTY_TYPE_MAGIC,
		Name:      "Bloody Spear",
		ConsumeMP: 60,
		P2M:       BloodSpearP2M,
	}
	skillTable[SKILL_PARALYZE] = &SkillInfo{
		Type:      SKILL_PROPERTY_TYPE_MAGIC,
		Name:      "Paralyze",
		ConsumeMP: 30,
		P2M:       ParalyzeP2M,
	}
}

func BloodSpearP2M(player *Player, monster *Monster) {
	fail := &darkeden.GCSkillFailed1Packet{
		SkillType: SKILL_BLOOD_SPEAR,
	}
	player.send <- fail
}

func ParalyzeP2M(player *Player, monster *Monster) {
	ok := &darkeden.GCSkillToObjectOK1{
		SkillType:      SKILL_PARALYZE,
		TargetObjectID: monster.Id(),
		Duration:       40,
	}
	player.send <- ok
}
