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

	P2Z func(*Player, *Scene, uint, uint)
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
	player.send <- &darkeden.GCSkillToObjectOK1{
		SkillType:      SKILL_BLOOD_SPEAR,
		CEffectID:      0,
		TargetObjectID: monster.Id(),
	}
	player.Scene.BroadcastPacket(player.X(), player.Y(), &darkeden.GCSkillToObjectOK3{
		ObjectID:  player.Id(),
		SkillType: SKILL_BLOOD_SPEAR,
		TargetX:   monster.X(),
		TargetY:   monster.Y(),
	})
	player.Scene.BroadcastPacket(monster.X(), monster.Y(), &darkeden.GCSkillToObjectOK4{
		ObjectID:  player.Id(),
		SkillType: SKILL_BLOOD_SPEAR,
	})
	
	damage := player.STR[ATTR_CURRENT]/6 + player.INT[ATTR_CURRENT]/2 + player.DEX[ATTR_CURRENT]/12
	if damage >= 180 {
		damage = 180
	}
}

func ParalyzeP2M(player *Player, monster *Monster) {
	ok := &darkeden.GCSkillToObjectOK1{
		SkillType:      SKILL_PARALYZE,
		TargetObjectID: monster.Id(),
		Duration:       (3 + player.INT[ATTR_CURRENT]/15) * 10,
	}
	player.send <- ok
}

type SkillOutput struct {
	MonsterID uint32
	Damage    int
	Duration  int
}
