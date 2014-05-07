package scene

import (
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/packet/darkeden"
	"github.com/tiancaiamao/ouster/player"
	"log"
)

func loop(m *Zone) {
	for {
		for id, player := range m.players {
			select {
			case msg := <-player.read:
				m.processPlayerInput(uint32(id), msg)
			default:
				break
			}
		}

		select {
		case <-m.quit:
			// 处理退出消息
		case <-m.event:
			// 处理地图事件...比如boss刷了之类的
		case <-m.heartbeat:
			//50ms的心跳,目前还不确定做什么...npc,怪物...player的逻辑放到
			//自身的goroutine中
			m.HeartBeat()
		}
	}
}

func (m *Zone) Go() {
	go loop(m)
}

func (m *Zone) processPlayerInput(playerId uint32, msg interface{}) {
	switch msg.(type) {
	case darkeden.CGMovePacket:

	case packet.CMovePacket:
		log.Println("scene receive and process a CMovePacket")
		raw := msg.(packet.CMovePacket)
		handle := m.Player(playerId)
		pc := handle.pc
		switch pc.State {
		case player.STAND, player.MOVE:
			pc.State = player.MOVE
			handle.to.X = raw.X
			handle.to.Y = raw.Y

			// boardcast to it's nearby players
			smove := packet.SMovePacket{
				Id:  playerId,
				Cur: handle.pos,
				To:  handle.to,
			}
			nearby := pc.NearBy()
			for _, playerId := range nearby {
				p := m.Player(playerId)
				if p != nil {
					p.write <- smove
				}
			}
		}
		handle.write <- player.CMovePacketAck{}
	case player.SkillEffect:
		log.Println("scene receive and process a SkillEffect")
		raw := msg.(player.SkillEffect)
		handle := m.Player(playerId)
		pc := handle.pc

		nearby := pc.NearBy()
		for _, playerId := range nearby {
			p := m.Player(playerId)
			if p != nil {
				p.write <- packet.SkillTargetEffectPacket{
					Skill: raw.Id,
					From:  playerId,
					To:    raw.To,
					Hurt:  raw.Hurt,
					Succ:  raw.Succ,
				}
			}
		}
	}
}
