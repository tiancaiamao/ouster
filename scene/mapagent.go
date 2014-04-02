package scene

import (
	// "runtime"
	"github.com/tiancaiamao/ouster/packet"
)

func loop(m *Map) {
	for {
		for id, player := m.players.Begin(); m.players.Valid(); id, player = m.players.Next() {
			select {
			case msg := <-player.ch:
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
	case packet.MovePacket:
		raw := msg.(packet.MovePacket)
		player := m.Player(playerId)
		switch player.state {
		case STAND, MOVE:
			player.state = MOVE
			player.to.X = raw.X
			player.to.Y = raw.Y

			// boardcast to it's nearby players
			nearby := player.this.NearBy()
			for _, playerId := range nearby {
				p := m.Player(playerId)
				if p != nil {
					p.ch <- msg
				}
			}
		}
	}
}
