package scene

import (
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/player"
	"math"
)

type PlayerArray struct {
	players []*player.Player
	slot    uint32
	empty   int
	iter    uint32
}

func (this *PlayerArray) Begin() (uint32, *player.Player) {
	this.iter = uint32(0)
	return this.Next()
}

func (this *PlayerArray) Next() (uint32, *player.Player) {
	for ; this.iter < uint32(len(this.players)); this.iter++ {
		if this.players[this.iter] != nil {
			return this.iter, this.players[this.iter]
		}
	}
	return this.iter, nil
}

func (this *PlayerArray) Valid() bool {
	return this.iter < uint32(len(this.players))
}

type Map struct {
	data.Map
	players *PlayerArray
	aoi     *aoi.Aoi

	quit      chan struct{}
	event     chan interface{}
	heartbeat chan struct{}
}

func New(m *data.Map) *Map {
	ret := new(Map)
	ret.players = &PlayerArray{
		players: make([]*player.Player, 0, 200),
	}

	ret.quit = make(chan struct{})
	ret.event = make(chan interface{})
	ret.heartbeat = make(chan struct{})
	return ret
}

func (m *Map) Player(playerId uint32) *player.Player {
	if playerId >= uint32(len(m.players.players)) {
		return nil
	}
	return m.players.players[playerId]
}

func (m *Map) HeartBeat() {
	for _, pc := m.players.Begin(); m.players.Valid(); _, pc = m.players.Next() {
		if pc.State == player.MOVE {
			v := pc.Speed()
			if ouster.Distance(pc.Pos, pc.To) < v {
				pc.State = player.STAND
				pc.Pos.X = pc.To.X
				pc.Pos.Y = pc.To.Y
			} else {
				dx := pc.To.X - pc.Pos.X
				dy := pc.To.Y - pc.Pos.Y
				angle := math.Atan2(float64(dy), float64(dx))
				vx := v * float32(math.Cos(angle))
				vy := v * float32(math.Sin(angle))

				pc.Pos.X += vx
				pc.Pos.Y += vy
			}
		}
	}
}

func (m *Map) Login(player *player.Player) error {
	var idx uint32
	if m.players.slot == uint32(len(m.players.players)-1) {
		if m.players.empty*4 > len(m.players.players) {
			for i := m.players.slot; i < uint32(len(m.players.players)); i++ {
				if m.players.players[i] == nil {
					m.players.players[i] = player
					m.players.slot = i

					idx = i
				}
			}
		} else {
			m.players.players = append(m.players.players, player)
			idx = uint32(len(m.players.players) - 1)
		}
	}

	player.Id = idx
	player.Scene = m.Title
	return nil
}
