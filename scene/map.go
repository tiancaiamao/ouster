package scene

import (
	"errors"
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/player"
	"math"
)

type PlayerState uint8

const (
	STAND PlayerState = iota
	MOVE
)

// pos attribute of a player is part of map, not player!
// while ch is reference of player.Player, not part of map!
// all Move related attribute of player are here.
type Player struct {
	pos   ouster.FPoint
	state PlayerState
	to    ouster.FPoint
	ch    chan interface{}
	this  *player.Player
}

type PlayerArray struct {
	players []Player
	slot    uint32
	empty   int
	iter    uint32
}

func (this *PlayerArray) Begin() (uint32, *Player) {
	this.iter = uint32(0)
	return this.Next()
}

func (this *PlayerArray) Next() (uint32, *Player) {
	for ; this.iter < uint32(len(this.players)); this.iter++ {
		if this.players[this.iter].ch != nil {
			return this.iter, &this.players[this.iter]
		}
	}
	return this.iter, nil
}

func (this *PlayerArray) Valid() bool {
	return this.iter == uint32(len(this.players)-1)
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
		players: make([]Player, 0, 200),
	}

	ret.quit = make(chan struct{})
	ret.event = make(chan interface{})
	ret.heartbeat = make(chan struct{})
	return ret
}

func (m *Map) PlayerPosition(playerId uint32) (ouster.FPoint, error) {
	p := m.players.players[playerId]
	if p.ch == nil {
		return ouster.FPoint{}, errors.New("no player correspond to this playerId")
	}
	return p.pos, nil
}

func (m *Map) Player(playerId uint32) *Player {
	if playerId >= uint32(len(m.players.players)) {
		return nil
	}
	p := m.players.players[playerId]
	if p.ch == nil {
		return nil
	}
	
	return &m.players.players[playerId]
}

func (m *Map) HeartBeat() {
	for _, player := m.players.Begin(); m.players.Valid(); _, player = m.players.Next() {
		if player.state == MOVE {
			v := player.this.Speed()
			if ouster.Distance(player.pos, player.to) < v {
				player.state = STAND
				player.pos.X = player.to.X
				player.pos.Y = player.to.Y
			} else {
				dx := player.to.X - player.pos.X
				dy := player.to.Y - player.pos.Y
				angle := math.Atan2(float64(dy), float64(dx))
				vx := v * float32(math.Cos(angle))
				vy := v * float32(math.Sin(angle))

				player.pos.X += vx
				player.pos.Y += vy
			}
		}
	}
}

func (m *Map) Login(pos ouster.FPoint, ch chan interface{}) (playerId uint32, ok bool) {
	if m.players.slot == uint32(len(m.players.players)-1) {
		if m.players.empty*4 > len(m.players.players) {
			for i := m.players.slot; i < uint32(len(m.players.players)); i++ {
				if m.players.players[i].ch == nil {
					m.players.players[i].ch = ch
					m.players.players[i].pos = pos
					m.players.slot = i
					return i, true
				}
			}
		} else {
			m.players.players = append(m.players.players, Player{
				pos:   pos,
				ch:    ch,
				state: STAND,
			})
			return uint32(len(m.players.players) - 1), ok
		}
	}
	return m.players.slot, true
}
