package main

import (
	"github.com/tiancaiamao/ouster/packet/darkeden"
)

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
