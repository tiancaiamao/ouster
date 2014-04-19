package scene

import (
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/player"
	"math"
	"time"
)

type Handle struct {
	// Player don't expose public field, but provide getter
	// so we can read but not write in this package
	pc *player.Player

	aoi   chan<- uint32
	write chan<- interface{}
	read  <-chan interface{}
	pos   ouster.FPoint
	to    ouster.FPoint
}

type Map struct {
	data.Map
	players []Handle
	aoi     *aoi.Aoi

	quit      chan struct{}
	event     chan interface{}
	heartbeat <-chan time.Time
}

func New(m *data.Map) *Map {
	ret := new(Map)
	ret.players = make([]Handle, 0, 200)
	ret.aoi = aoi.New()

	ret.quit = make(chan struct{})
	ret.event = make(chan interface{})
	ret.heartbeat = time.Tick(50 * time.Microsecond)

	return ret
}

func (m *Map) Player(playerId uint32) *Handle {
	if playerId >= uint32(len(m.players)) {
		return nil
	}
	return &m.players[playerId]
}

func (m *Map) Pos(playerId uint32) (ouster.FPoint, error) {
	handle := m.Player(playerId)
	if handle == nil {
		return ouster.FPoint{}, ouster.NewError("query a non-exist id")
	}
	return handle.pos, nil
}

func (m *Map) To(playerId uint32) (ouster.FPoint, error) {
	handle := m.Player(playerId)
	if handle == nil {
		return ouster.FPoint{}, ouster.NewError("query a non-exist id")
	}
	return handle.to, nil
}

func (m *Map) Creature(id int) ouster.Creature {
	handle := m.Player(uint32(id))
	return handle.pc
}

func (m *Map) String() string {
	return m.Map.Name
}

func (m *Map) HeartBeat() {
	for id, handle := range m.players {
		if handle.pc == nil {
			continue
		}

		// process player move
		pc := handle.pc
		if pc.State == player.MOVE {
			v := pc.Speed()
			if ouster.Distance(handle.pos, handle.to) < v {
				pc.State = player.STAND
				handle.pos.X = handle.to.X
				handle.pos.Y = handle.to.Y
			} else {
				dx := handle.to.X - handle.pos.X
				dy := handle.to.Y - handle.pos.Y
				angle := math.Atan2(float64(dy), float64(dx))
				vx := v * float32(math.Cos(angle))
				vy := v * float32(math.Sin(angle))

				handle.pos.X += vx
				handle.to.Y += vy
			}

			// aoi update
			m.aoi.Update(uint32(id), aoi.ModeWatcher|aoi.ModeMarker, aoi.FPoint(handle.pos))
		}
	}

	m.aoi.Message(func(watcher uint32, marker uint32) {
		if watcher >= 0 && watcher < uint32(len(m.players)) {
			handle := m.players[watcher]
			if handle.pc != nil {
				handle.aoi <- marker
			}
		}
	})
}

func (m *Map) Login(player *player.Player, pos ouster.FPoint, a chan<- uint32, rd <-chan interface{}, wr chan<- interface{}) error {
	var handle Handle
	handle.pc = player
	handle.pos = pos
	handle.to = pos
	handle.read = rd
	handle.write = wr

	idx := len(m.players)
	m.players = append(m.players, handle)

	player.Id = uint32(idx)
	player.Scene = m

	m.aoi.Update(player.Id, aoi.ModeWatcher|aoi.ModeMarker, aoi.FPoint(pos))

	return nil
}
