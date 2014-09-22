package main

import (
    . "github.com/tiancaiamao/ouster/util"
)

var (
    VampireTimebandFactor   = [4]int{125, 100, 125, 150}
    MonsterTimebandFactor   = [4]int{75, 50, 75, 100}
    AttrExpTimebandFactor   = [4]int{100, 100, 100, 150}
    DomainExpTimebandFactor = [4]int{100, 100, 100, 150}
)

type SkillPropertyType uint8

const (
    SKILL_PROPERTY_TYPE_MELEE SkillPropertyType = iota
    SKILL_PROPERTY_TYPE_MAGIC
    SKILL_PROPERTY_TYPE_PHYSIC
)

const (
    SKILL_ATTACK_MELEE       = 0
    SKILL_UN_BURROW          = 107
    SKILL_UN_TRANSFORM       = 108
    SKILL_UN_INVISIBILITY    = 109
    SKILL_RAPID_GLIDING      = 203
    SKILL_METEOR_STRIKE      = 180
    SKILL_INVISIBILITY       = 100
    SKILL_PARALYZE           = 89
    SKILL_BLOOD_SPEAR        = 97
    SKILL_ABSORB_SOUL        = 246
    SKILL_SUMMON_SYLPH       = 247
    SKILL_SHARP_HAIL         = 348 // 尖锐冰雹
    SKILL_FLOURISH           = 219 // 活跃攻击
    SKILL_DESTRUCTION_SPEAR  = 298 //致命爆发
    SKILL_SHARP_CHAKRAM      = 295 // 税利之轮
    SKILL_EVADE              = 220 // 回避术
    SKILL_FIRE_OF_SOUL_STONE = 227
    SKILL_ICE_OF_SOUL_STONE  = 228
    SKILL_SAND_OF_SOUL_STONE = 229
    SKILL_TELEPORT           = 280 // 瞬间移动
    SKILL_DUCKING_WALLOP     = 302 // 光速冲击
    SKILL_DISTANCE_BLITZ     = 304 // 雷神斩
)

type SkillFormula interface {
    ComputeOutput(*Creature, *Creature) SkillOutput
}

type PlayerSkillToTileHandler interface {
    ExecuteP2T(*Player, uint8, uint8)
}

type MonsterSkillToTileHandler interface {
    ExecuteM2T(*Monster, uint8, uint8)
}

type PlayerSkillToMonsterHandler interface {
    ExecuteP2M(*Player, *Monster)
}

type MonsterSkillToPlayerHandler interface {
    ExecuteM2P(*Monster, *Player)
}

type PlayerSkillToPlayerHandler interface {
    ExecuteP2P(*Player, *Player)
}

type SkillToSelfHandler interface {
    Execute(*Player)
}

type Skill struct {
}

// 派生类中重写这个函数
func (skill Skill) ComputeOutput(*SkillInput, *SkillOutput) {}
func (skill Skill) Check(skillType SkillType_t, agent *Agent) bool {
    // skillSlot := agent.hasSkill(skillType)
    // skillInfo := skillInfoTable[skillType]
    requireMP := decreaseConsumeMP(agent)

    // var hitBonus int
    // if agent.hasRankBonus(RANK_BONUS_KNOWLEDGE_OF_INNATE) {
    //     rankBonus := agent.getRankBonus(RANK_BONUS_KNOWLEDGE_OF_INNATE)
    //     hitBonus = rankBonus.getPoint()
    // }

    manaCheck := hasEnoughMana(agent, requireMP)
    if !manaCheck {
        return false
    }

    // timeCheck := verifyRuntime(skillSlot)
    // if !timcCheck {
    //     return false
    // }

    // hitRoll := HitRoll.isSuccessMagic(ouster, skillInfo, skillSlot, hitBonus)
    // if !hitRoll {
    //     return false
    // }

    pc := agent.PlayerCreatureInstance()
    effected := pc.IsFlag(EFFECT_CLASS_INVISIBILITY) || pc.IsFlag(EFFECT_CLASS_HAS_FLAG) || pc.IsFlag(EFFECT_CLASS_HAS_SWEEPER)
    if effected {
        return false
    }

    decreaseMana(agent, requireMP)
    // skillSlot.setRunTime(output.Delay)
    return true
}

var skillTable map[SkillType_t]SkillHandler

func init() {
    skillTable = make(map[SkillType_t]SkillHandler)

    skillTable[SKILL_ATTACK_MELEE] = AttackMelee{}
    // skillTable[SKILL_BLOOD_SPEAR] = &SkillInfo{
    //     Type:      SKILL_PROPERTY_TYPE_MAGIC,
    //     Name:      "Bloody Spear",
    //     ConsumeMP: 60,
    // }
    // skillTable[SKILL_PARALYZE] = &SkillInfo{
    //     Type:      SKILL_PROPERTY_TYPE_MAGIC,
    //     Name:      "Paralyze",
    //     ConsumeMP: 30,
    // }
    // skillTable[SKILL_RAPID_GLIDING] = &SkillInfo{
    //     Type:      SKILL_PROPERTY_TYPE_MAGIC,
    //     Name:      "Rapid Gliding",
    //     ConsumeMP: 23,
    // }
    // skillTable[SKILL_INVISIBILITY] = &SkillInfo{
    //     Type:      SKILL_PROPERTY_TYPE_MAGIC,
    //     Name:      "Invisibility",
    //     ConsumeMP: 36,
    // }
    // skillTable[SKILL_METEOR_STRIKE] = &SkillInfo{
    //     Type:      SKILL_PROPERTY_TYPE_MAGIC,
    //     Name:      "Meteor Strike",
    //     ConsumeMP: 53,
    // }
    //
    // skillTable[SKILL_SHARP_HAIL] = &SkillInfo{
    //     Type:      SKILL_PROPERTY_TYPE_PHYSIC,
    //     Name:      "Sharp Hail",
    //     ConsumeMP: 20,
    // }
}

type Invisibility struct {
    Skill
}

type MeteorStrike struct {
    Skill
}
type Paralyze struct {
    Skill
}
type AttackMelee struct {
    Skill
}

func AOE(scene *Scene,
    player *Player,
    tileX uint8,
    tileY uint8,
    skill *SkillInfo,
    slot *SkillSlot,
) {
    // middleX := uint8((int(player.X) + int(tileX)) / 2)
    // middleY := uint8((int(player.Y) + int(tileY)) / 2)
    // var middleX, middleY uint8
    // scene.Nearby(middleX,
    // middleY,
    // func(ignore aoi.Entity, target aoi.Entity) {
    // id := target.Id()
    // if id == player.Id() {
    //     return
    // }

    // obj := scene.objects[id]

    // broadcast skill to nearby player
    // if player, ok := obj.(*Player); ok {
    // 	if Distance2(target.X(), target.Y(), middleX, middleY) <= 144 {
    // 		player.send <- &darkeden.GCSkillToTileOK5{
    // 			ObjectID:  player.Id(),
    // 			SkillType: slot.SkillType,
    // 			X:         player.X(),
    // 			Y:         player.Y(),
    // 			Range:     skill.Range,
    // 			Duration:  slot.Duration,
    // 			// CreatureList []uint32
    // 		}
    // 	}
    //
    // 	if Distance2(target.X(), target.Y(), player.X(), player.Y()) <= 64 {
    // 		player.send <- &darkeden.GCSkillToTileOK4{
    // 			SkillType: slot.SkillType,
    // 			X:         player.X(),
    // 			Y:         player.Y(),
    // 			Range:     skill.Range,
    // 			Duration:  slot.Duration,
    // 			// CreatureList []uint32
    // 		}
    // 	}
    //
    // 	if Distance2(target.X(), target.Y(), tileX, tileY) <= 64 {
    // 		player.send <- &darkeden.GCSkillToTileOK3{
    // 			ObjectID:  target.Id(),
    // 			SkillType: slot.SkillType,
    // 			X:         player.X(),
    // 			Y:         player.Y(),
    // 		}
    // 	}
    // }

    // attack nearby enemy
    // if inRange(skill.Range, target.X(), target.Y(), tileX, tileY) {
    // switch obj.(type) {
    // case *Monster:
    //     monster := obj.(*Monster)
    //     handler := skill.Handler.(SkillFormula)
    //     output := handler.ComputeOutput(&player.Creature, &monster.Creature)
    //     // scene is the owner of monster
    //     scene.agent <- AgentMessage{
    //         Player: player,
    //         Msg:    output,
    //     }
    // case *Player:
    //
    // }
    // }
    // })
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
