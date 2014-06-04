package main

import (
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/packet/darkeden"
)

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
