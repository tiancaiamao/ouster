package main

import (
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet/darkeden"
)

func loop(m *Zone) {
	for {
		for id, player := range m.players {
			select {
			case msg := <-player.agent2scene:
				m.processPlayerInput(uint32(id), msg)
			default:
				break
			}
		}

		select {
		case <-m.quit:
		case <-m.event:
		case <-m.heartbeat:
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
		move := msg.(darkeden.CGMovePacket)
		player := m.Player(playerId)
		player.send <- darkeden.GCMoveOKPacket{
			X:   move.X,
			Y:   move.Y,
			Dir: move.Dir,
		}
		m.aoi.Nearby(uint16(move.X), uint16(move.Y), func(entity *aoi.Entity) {
			dx := int(move.X) - int(entity.X())
			dy := int(move.Y) - int(entity.Y())
			if dx*dx+dy*dy <= 64 {
				player.handleAoiMessage(ObjectIDType(entity.Id()))
			}
		})
	}
}

var (
	maps map[string]*Zone
)

func Initialize() {
	maps = make(map[string]*Zone)
	maps["limbo_lair_se"] = New(&data.LimboLairSE)

	for _, m := range maps {
		m.Go()
	}
}

func Query(mapName string) *Zone {
	return maps[mapName]
}
