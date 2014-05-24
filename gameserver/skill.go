package main

import (
	"github.com/tiancaiamao/ouster/packet/darkeden"
)

var skillP2M map[uint16]func(*Player, *Monster)
var skillM2P map[uint16]func(*Monster, *Player)
var skillP2P map[uint16]func(*Player, *Player)

func init() {
	skillP2M = make(map[uint16]func(*Player, *Monster))
	skillM2P = make(map[uint16]func(*Monster, *Player))
	skillP2P = make(map[uint16]func(*Player, *Player))
	
	skillP2M[darkeden.SKILL_BLOOD_SPEAR] = BloodSpearP2M
	skillP2M[darkeden.SKILL_PARALYZE] = ParalyzeP2M
}

func BloodSpearP2M(player *Player, monster *Monster) {
		
}

func ParalyzeP2M(player *Player, monster *Monster) {
	ok := &darkeden.GCSkillToObjectOK1{
		SkillType:      darkeden.SKILL_PARALYZE,
		TargetObjectID: monster.Id(),
		Duration:       40,
	}
	player.send <- ok
}