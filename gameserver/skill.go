package main

import (
	"github.com/tiancaiamao/ouster/aoi"
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

	SKILL_ABSORB_SOUL  uint16 = 246
	SKILL_SUMMON_SYLPH uint16 = 247
	SKILL_FLOURISH     uint16 = 219
)

type SkillFormula interface {
	ComputeOutput(*Creature, *Creature) SkillOutput
}

type SkillToTileHandler interface {
	ExecuteP2T(*Player, uint8, uint8)
	ExecuteM2T(*Monster, uint8, uint8)
}

type SkillToObjectHandler interface {
	ExecuteP2M(*Player, *Monster)
	ExecuteM2P(*Monster, *Player)
	ExecuteP2P(*Player, *Player)
}

type SkillToSelfHandler interface {
	Execute(*Player)
}

type SkillInfo struct {
	Type        SkillPropertyType
	Name        string
	Level       uint
	MinDamage   uint
	MaxDamage   uint
	MinDelay    uint
	MinCastTime uint
	MaxCastTime uint
	ConsumeMP   uint
	Handler     interface{}
}

type BloodSpearHandler struct{}

func (blood BloodSpearHandler) ExecuteP2M(player *Player, monster *Monster) {
	player.send <- &darkeden.GCSkillToObjectOK1{
		SkillType:      SKILL_BLOOD_SPEAR,
		CEffectID:      0,
		TargetObjectID: monster.Id(),
	}
	player.BroadcastPacket(player.X(), player.Y(), &darkeden.GCSkillToObjectOK3{
		ObjectID:  player.Id(),
		SkillType: SKILL_BLOOD_SPEAR,
		TargetX:   monster.X(),
		TargetY:   monster.Y(),
	})
	player.BroadcastPacket(monster.X(), monster.Y(), &darkeden.GCSkillToObjectOK4{
		ObjectID:  player.Id(),
		SkillType: SKILL_BLOOD_SPEAR,
	})

	damage := player.STR[ATTR_CURRENT]/6 + player.INT[ATTR_CURRENT]/2 + player.DEX[ATTR_CURRENT]/12
	if damage >= 180 {
		damage = 180
	}
	player.Scene.agent <- AgentMessage{
		Player: player,
		Msg: SkillOutput{
			MonsterID: monster.Id(),
			Damage:    int(damage),
			Duration:  10,
		},
	}
}
func (blood BloodSpearHandler) ExecuteM2P(monster *Monster, player *Player) {
	// TODO
}
func (blood BloodSpearHandler) ExecuteP2P(p1 *Player, p2 *Player) {
	// TODO
}

var skillTable map[uint16]*SkillInfo

func init() {
	skillTable = make(map[uint16]*SkillInfo)

	skillTable[SKILL_BLOOD_SPEAR] = &SkillInfo{
		Type:      SKILL_PROPERTY_TYPE_MAGIC,
		Name:      "Bloody Spear",
		ConsumeMP: 60,
		Handler:   BloodSpearHandler{},
	}
	skillTable[SKILL_PARALYZE] = &SkillInfo{
		Type:      SKILL_PROPERTY_TYPE_MAGIC,
		Name:      "Paralyze",
		ConsumeMP: 30,
		Handler:   ParalyzeHandler{},
	}
	skillTable[SKILL_RAPID_GLIDING] = &SkillInfo{
		Type:      SKILL_PROPERTY_TYPE_MAGIC,
		Name:      "Rapid Gliding",
		ConsumeMP: 23,
		Handler:   RapidGlidingHandler{},
	}
	skillTable[SKILL_INVISIBILITY] = &SkillInfo{
		Type:      SKILL_PROPERTY_TYPE_MAGIC,
		Name:      "Invisibility",
		ConsumeMP: 36,
		Handler:   InvisibilityHandler{},
	}
	skillTable[SKILL_METEOR_STRIKE] = &SkillInfo{
		Type:      SKILL_PROPERTY_TYPE_MAGIC,
		Name:      "Meteor Strike",
		ConsumeMP: 53,
		Handler:   MeteorStrikeHandler{},
	}
}

type InvisibilityHandler struct{}

func (ignore InvisibilityHandler) Execute(player *Player) {
	// TODO
}

type RapidGlidingHandler struct{}

func (ignore RapidGlidingHandler) ExecuteP2T(player *Player, x uint8, y uint8) {
	fastMove := &darkeden.GCFastMovePacket{
		ObjectID:  player.Id(),
		FromX:     player.X(),
		FromY:     player.Y(),
		ToX:       x,
		ToY:       y,
		SkillType: SKILL_RAPID_GLIDING,
	}
	player.Scene.agent <- AgentMessage{
		Player: player,
		Msg:    fastMove,
	}

	ok := &darkeden.GCSkillToTileOK1{
		SkillType: SKILL_RAPID_GLIDING,
		Duration:  10,
		Range:     1,
		X:         x,
		Y:         y,
	}
	player.send <- ok
}
func (ignore RapidGlidingHandler) ExecuteM2T(monster *Monster, x uint8, y uint8) {
	// TODO
}

type MeteorStrikeHandler struct{}

func (ignore MeteorStrikeHandler) ExecuteP2T(player *Player, x uint8, y uint8) {
	player.Scene.Nearby(x, y, func(watcher aoi.Entity, marker aoi.Entity) {
		if x >= marker.X()-1 &&
			x <= marker.X()+1 &&
			y >= marker.Y()-1 &&
			y <= marker.Y()+1 {
			id := marker.Id()
			obj := player.Scene.objects[id]
			switch obj.(type) {
			case *Monster:
				monster := obj.(*Monster)
				skillOutput := ignore.ComputeOutput(&player.Creature, &monster.Creature)
				player.Scene.agent <- AgentMessage{
					Player: player,
					Msg:    skillOutput,
				}
			case *Player:

			}
		}
	})
	ok := &darkeden.GCSkillToTileOK1{
		SkillType: SKILL_METEOR_STRIKE,
		Duration:  10,
		Range:     3,
		X:         x,
		Y:         y,
	}
	player.send <- ok
}
func (ignore MeteorStrikeHandler) ComputeOutput(c1 *Creature, c2 *Creature) SkillOutput {
	return SkillOutput{
		Damage: int(float32(c1.Level)*0.8) + int(c1.STR[ATTR_CURRENT]+c1.DEX[ATTR_CURRENT])/6,
	}
}

type ParalyzeHandler struct{}

func (ignore ParalyzeHandler) ExecuteP2M(player *Player, monster *Monster) {
	skillOutput := ignore.ComputeOutput(&player.Creature, &monster.Creature)
	ok := &darkeden.GCSkillToObjectOK1{
		SkillType:      SKILL_PARALYZE,
		TargetObjectID: monster.Id(),
		Duration:       uint16(skillOutput.Duration),
	}
	player.send <- ok
	player.Scene.BroadcastPacket(player.X(), player.Y(), &darkeden.GCSkillToObjectOK3{
		ObjectID:  player.Id(),
		SkillType: SKILL_PARALYZE,
		TargetX:   monster.X(),
		TargetY:   monster.Y(),
	})
	player.Scene.BroadcastPacket(monster.X(), monster.Y(), &darkeden.GCSkillToObjectOK4{
		ObjectID:  monster.Id(),
		SkillType: SKILL_PARALYZE,
		Duration:  uint16(skillOutput.Duration),
	})
}
func (ignore ParalyzeHandler) ExecuteM2P(monster *Monster, player *Player) {
	// TODO
}
func (ignore ParalyzeHandler) ExecuteP2P(p1 *Player, p2 *Player) {
	// TODO
}
func (ignore ParalyzeHandler) ComputeOutput(c1 *Creature, c2 *Creature) SkillOutput {
	return SkillOutput{
		Duration: int((3 + c1.INT[ATTR_CURRENT]/15) * 10),
	}
}

type SkillOutput struct {
	MonsterID uint32
	Damage    int
	Duration  int
}
