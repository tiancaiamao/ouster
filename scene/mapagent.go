package scene

import (
	// "runtime"
	// "github.com/tiancaiamao/ouster/packet"
)

func loop(m *Map) {
	for {		
		for id, player := m.players.Begin(); m.players.Valid(); id, player = m.players.Next() {
			select {
			case msg := <-player.ch:
				processPlayerInput(id, msg)
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

func processPlayerInput(playerId uint32, msg interface{}) {
	switch msg.(type) {
	// case packet.Move:
	// 	raw := packet.Move(msg)

	// case packet.Skill:

	}
}
