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

	SKILL_ABSORB_SOUL       uint16 = 246
	SKILL_SUMMON_SYLPH      uint16 = 247
	SKILL_SHARP_HAIL        uint16 = 348 // 尖锐冰雹
	SKILL_FLOURISH          uint16 = 219 // 活跃攻击
	SKILL_DESTRUCTION_SPEAR uint16 = 298 //致命爆发
	SKILL_SHARP_CHAKRAM     uint16 = 295 // 税利之轮
	SKILL_EVADE             uint16 = 220 // 回避术

	SKILL_FIRE_OF_SOUL_STONE uint16 = 227
	SKILL_ICE_OF_SOUL_STONE  uint16 = 228
	SKILL_SAND_OF_SOUL_STONE uint16 = 229

	SKILL_TELEPORT       uint16 = 280 // 瞬间移动
	SKILL_DUCKING_WALLOP uint16 = 302 // 光速冲击
	SKILL_DISTANCE_BLITZ uint16 = 304 // 雷神斩

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
	Range       uint8
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

	skillTable[SKILL_SHARP_HAIL] = &SkillInfo{
		Type:      SKILL_PROPERTY_TYPE_PHYSIC,
		Name:      "Sharp Hail",
		ConsumeMP: 20,
		Handler:   SharpHailHandler{},
	}
}

type SharpHailHandler struct{}

func (ignore SharpHailHandler) ExecuteP2T(player *Player, x uint8, y uint8) {
	slot := player.SkillSlot(SKILL_SHARP_HAIL)

	ok := &darkeden.GCSkillToTileOK1{
		SkillType: SKILL_SHARP_HAIL,
		Duration:  10,
		Range:     5,
		X:         x,
		Y:         y,
	}
	player.send <- ok

	AOE(player.Scene, player, x, y, skillTable[SKILL_SHARP_HAIL], slot)
}
func (ignore SharpHailHandler) ExecuteM2T(monster *Monster, x uint8, y uint8) {}

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

func AOE(scene *Scene,
	player *Player,
	tileX uint8,
	tileY uint8,
	skill *SkillInfo,
	slot *SkillSlot,
) {
	middleX := uint8((int(player.X()) + int(tileX)) / 2)
	middleY := uint8((int(player.Y()) + int(tileY)) / 2)

	scene.Nearby(middleX,
		middleY,
		func(ignore aoi.Entity, target aoi.Entity) {
			id := target.Id()
			obj := scene.objects[id]

			// broadcast skill to nearby player
			if player, ok := obj.(*Player); ok {
				if Distance2(target.X(), target.Y(), middleX, middleY) <= 144 {
					player.send <- &darkeden.GCSkillToTileOK5{
						ObjectID:  player.Id(),
						SkillType: slot.SkillType,
						X:         player.X(),
						Y:         player.Y(),
						Range:     skill.Range,
						Duration:  slot.Duration,
						// CreatureList []uint32
					}
				}

				if Distance2(target.X(), target.Y(), player.X(), player.Y()) <= 64 {
					player.send <- &darkeden.GCSkillToTileOK4{
						SkillType: slot.SkillType,
						X:         player.X(),
						Y:         player.Y(),
						Range:     skill.Range,
						Duration:  slot.Duration,
						// CreatureList []uint32
					}
				}

				if Distance2(target.X(), target.Y(), tileX, tileY) <= 64 {
					player.send <- &darkeden.GCSkillToTileOK3{
						ObjectID:  target.Id(),
						SkillType: slot.SkillType,
						X:         player.X(),
						Y:         player.Y(),
					}
				}
			}

			// attack nearby enemy
			if inRange(skill.Range, target.X(), target.Y(), tileX, tileY) {
				switch obj.(type) {
				case *Monster:
					monster := obj.(*Monster)
					handler := skill.Handler.(SkillFormula)
					output := handler.ComputeOutput(&player.Creature, &monster.Creature)
					// scene is the owner of monster
					scene.agent <- AgentMessage{
						Player: player,
						Msg:    output,
					}
				case *Player:

				}
			}
		})
}

func Distance2(x1, y1, x2, y2 uint8) int {
	d1 := int(x1) - int(x2)
	d2 := int(y1) - int(y2)
	return d1*d1 + d2*d2
}

func inRange(radius uint8, x uint8, y uint8, middleX uint8, middleY uint8) bool {
	return x >= middleX-radius/2 &&
		x <= middleX+radius/2 &&
		y >= middleY-radius/2 &&
		y <= middleY-radius/2
}
