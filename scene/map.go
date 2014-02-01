package scene

import (
	"errors"
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/data"
)

// pos of a player is part of map, not player!
type Player struct {
	pos ouster.Point
	ch  chan interface{}
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

func (m *Map) PlayerPosition(playerId uint32) (ouster.Point, error) {
	p := m.players.players[playerId]
	if p.ch == nil {
		return ouster.Point{}, errors.New("no player correspond to this playerId")
	}
	return p.pos, nil
}

func (m *Map) HeartBeat() {

}

func (m *Map) Login(pos ouster.Point, ch chan interface{}) (playerId uint32, ok bool) {
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
			m.players.players = append(m.players.players, Player{pos, ch})
			return uint32(len(m.players.players) - 1), ok
		}
	}
	return m.players.slot, true
}
