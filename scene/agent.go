package scene

import (
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/player"
	"log"
)

func loop(m *Map) {
	for {
		for id, player := range m.players {
			select {
			case msg := <-player.Player2scene:
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

func (m *Map) Go() {
	go loop(m)
}

func (m *Map) processPlayerInput(playerId uint32, msg interface{}) {
	switch msg.(type) {
	case packet.CMovePacket:
		log.Println("scene receive and process a CMovePacket")
		raw := msg.(packet.CMovePacket)
		pc := m.Player(playerId)
		switch pc.State {
		case player.STAND, player.MOVE:
			pc.State = player.MOVE
			pc.To.X = raw.X
			pc.To.Y = raw.Y

			// boardcast to it's nearby players
			nearby := pc.NearBy()
			for _, playerId := range nearby {
				p := m.Player(playerId)
				if p != nil {
					p.Scene2player <- msg
				}
			}
		}
		pc.Scene2player <- player.CMovePacketAck{}
	}
}
