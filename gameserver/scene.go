package main

import (
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet/darkeden"
	"math/rand"
	"time"
)

type Scene struct {
	*data.Map

	players  []*Player
	monsters []Monster
	aoi      *aoi.CellAoi

	quit      chan struct{}
	event     chan interface{}
	heartbeat <-chan time.Time
}

const maskNPC uint32 = 1 << 31

func New(m *data.Map) *Scene {
	ret := new(Scene)
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

func (m *Scene) Blocked(x, y uint16) bool {
	return false
}

func (m *Scene) Player(playerId uint32) *Player {
	if playerId >= uint32(len(m.players)) {
		return nil
	}
	return m.players[playerId]
}

func (m *Scene) Monster(idx uint32) *Monster {
	if idx >= uint32(len(m.monsters)) {
		return nil
	}
	return &m.monsters[idx]
}

func (m *Scene) String() string {
	return m.Map.Name
}

func (m *Scene) HeartBeat() {
	for i := 0; i < len(m.monsters); i++ {
		monster := &m.monsters[i]
		if (monster.flag & flagActive) != 0 {
			monster.HeartBeat(m)
		}
	}
}

func (m *Scene) Login(player *Player) error {
	idx := len(m.players)
	m.players = append(m.players, player)

	player.Entity = m.aoi.Add(145, 237, uint32(idx))
	player.Scene = m

	return nil
}

func loop(m *Scene) {
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

func (m *Scene) Go() {
	go loop(m)
}

func (m *Scene) processPlayerInput(playerId uint32, msg interface{}) {
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
	maps map[string]*Scene
)

func Initialize() {
	maps = make(map[string]*Scene)
	maps["limbo_lair_se"] = New(&data.LimboLairSE)

	for _, m := range maps {
		m.Go()
	}
}

func Query(mapName string) *Scene {
	return maps[mapName]
}
