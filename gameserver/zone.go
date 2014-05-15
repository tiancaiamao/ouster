package main

import (
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
	"math/rand"
	"time"
)

type Zone struct {
	*data.Map

	players  []*Player
	monsters []Monster
	aoi      *aoi.CellAoi

	quit      chan struct{}
	event     chan interface{}
	heartbeat <-chan time.Time
}

const maskNPC uint32 = 1 << 31

func New(m *data.Map) *Zone {
	ret := new(Zone)
	ret.Map = m
	aoi := aoi.NewCellAoi(m.Width, m.Height, 32, 32)
	players := make([]*Player, 0, 200)

	num := 0
	for _, mi := range m.MonsterInfo {
		num += int(mi.Count)
	}
	monsters := make([]Monster, num)

	idx := 0
	for _, mi := range m.MonsterInfo {
		tp := data.MonsterType2MonsterInfo[mi.MonsterType]
		for i := 0; i < int(mi.Count); i++ {
			var x, y int
			// set monster's position and so on
			for {
				x = rand.Intn(int(m.Width))
				y = rand.Intn(int(m.Height))

				flag := m.Data[x*int(m.Width)+y]
				if flag == 0x0 && !ret.Blocked(uint16(x), uint16(y)) {
					break
				}
			}

			monster := &monsters[idx]
			monster.aoi = aoi.Add(uint16(x), uint16(y), uint32(idx)|ObjectIDMaskNPC)
			monster.MonsterType = mi.MonsterType
			monster.STR = tp.STR
			monster.DEX = tp.DEX
			monster.INT = tp.INTE
			monster.HP = tp.STR*4 + uint16(tp.Level)
			idx++
		}
	}

	ret.aoi = aoi
	ret.players = players
	ret.monsters = monsters
	ret.quit = make(chan struct{})
	ret.event = make(chan interface{})
	ret.heartbeat = time.Tick(50 * time.Millisecond)

	return ret
}

func (m *Zone) Blocked(x, y uint16) bool {
	return false
}

func (m *Zone) Player(playerId uint32) *Player {
	if playerId >= uint32(len(m.players)) {
		return nil
	}
	return m.players[playerId]
}

func (m *Zone) Monster(idx uint32) *Monster {
	if idx >= uint32(len(m.monsters)) {
		return nil
	}
	return &m.monsters[idx]
}

func (m *Zone) String() string {
	return m.Map.Name
}

func (m *Zone) HeartBeat() {
	for i := 0; i < len(m.monsters); i++ {
		monster := &m.monsters[i]
		if (monster.flag & flagActive) != 0 {
			monster.HeartBeat(m)
		}
	}
}

func (m *Zone) Login(player *Player) error {
	idx := len(m.players)
	m.players = append(m.players, player)

	player.Entity = m.aoi.Add(145, 237, uint32(idx))
	player.zone = m

	return nil
}
