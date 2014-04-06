package scene

import (
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/player"
)

func loop(m *Map) {
	for {
		for id, player := m.players.Begin(); m.players.Valid(); id, player = m.players.Next() {
			select {
			case msg := <-player.Player2scene:
				m.processPlayerInput(id, msg)
			}
		}

		select {
		case <-m.quit:
			// 处理退出消息
		case <-m.event:
			// 处理地图事件...比如boss刷了之类的
		case <-m.heartbeat:
			//100us的心跳,目前还不确定做什么...npc,怪物...player的逻辑放到自身的goroutine中
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
	}
}
