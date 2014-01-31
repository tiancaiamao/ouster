package ouster

func loop(m) {
	for {
		for playerId := range m.Players {
			c := m.getPlayerChan()
			select {
			case msg := <-c:
				processPlayerInput(playerId, msg)
			}
		}
		select {
		case <-quit:
			// 处理退出消息
		case <-event:
			// 处理地图事件...比如boss刷了之类的
		case <-heartbeat:
			//100us的心跳,目前还不确定做什么...npc,怪物...player的逻辑放到自身的goroutine中
			m.HeartBeat()
		default:
			runtime.Sched()
		}
	}
}

func MapGoroutine(m *Map) {
	go loop(m)
}

func processPlayerInput(uint32 playerId, msg interface{}) {
	switch msg.(type) {
	case packet.Move:
		raw := packet.Move(msg)
		
	case packet.Skill:
		
	}
}
