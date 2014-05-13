package scene

import (
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
	// "github.com/tiancaiamao/ouster/player"
	// "log"
	// "math"
	"math/rand"
	"time"
)

type Handle struct {
	// Player don't expose public field, but provide getter
	// so we can read but not write in this package
	pc *Player

	aoi   chan<- uint32
	write chan<- interface{}
	read  <-chan interface{}

	pos ouster.FPoint
	to  ouster.FPoint
}

// id 0xxxxxx player
// id 1xxxxxx non-player
// id 10xxxxx npc
// id 11xxxxx monster
// id 110xxxx

type Zone struct {
	*data.Map

	players  []Handle
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
	players := make([]Handle, 0, 200)

	num := 0
	for _, mi := range m.MonsterInfo {
		num += int(mi.Count)
	}
	monsters := make([]Monster, num)

	idx := 0
	for _, mi := range m.MonsterInfo {
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

func (m *Zone) Player(playerId uint32) *Handle {
	if playerId >= uint32(len(m.players)) {
		return nil
	}
	return &m.players[playerId]
}

func (m *Zone) Pos(playerId uint32) (ouster.FPoint, error) {
	handle := m.Player(playerId)
	if handle == nil {
		return ouster.FPoint{}, ouster.NewError("query a non-exist id")
	}
	return handle.pos, nil
}

func (m *Zone) To(playerId uint32) (ouster.FPoint, error) {
	handle := m.Player(playerId)
	if handle == nil {
		return ouster.FPoint{}, ouster.NewError("query a non-exist id")
	}
	return handle.to, nil
}

func (m *Zone) Monster(id uint32) *Monster {
	if id >= uint32(len(m.monsters)) {
		return nil
	}
	return &m.monsters[id]
}

func (m *Zone) Creature(id uint32) ouster.Creature {
	handle := m.Player(id)
	return handle.pc
}

func (m *Zone) String() string {
	return m.Map.Name
}

func (m *Zone) movePC() {
	for id := 0; id < len(m.players); id++ {
		handle := &m.players[id]
		if handle.pc == nil {
			continue
		}

		// pc := handle.pc
		// if pc.State == player.MOVE {
		// 	v := pc.Speed()
		// 	if ouster.Distance2(handle.pos, handle.to) <= v*v {
		// 		pc.State = player.STAND
		// 		handle.pos.X = handle.to.X
		// 		handle.pos.Y = handle.to.Y
		// 		// pc.SendPosSync()
		// 	} else {
		// 		dx := handle.to.X - handle.pos.X
		// 		dy := handle.to.Y - handle.pos.Y
		// 		angle := math.Atan2(float64(dy), float64(dx))
		// 		vx := v * float32(math.Cos(angle))
		// 		vy := v * float32(math.Sin(angle))

		// newX := uint16(handle.pos.X + vx)
		// newY := uint16(handle.pos.Y + vy)

		// idx := newX*m.Width + newY
		// for _, layer := range m.Layers {
		//					 if layer.Type == data.BACKGROUND {
		//						 flag := layer.Data[idx]
		//						 if flag != 0 {
		//							 // encounter a obscure
		//							 pc.State = player.STAND
		//							 handle.pos.X = handle.to.X
		//							 handle.pos.Y = handle.to.Y
		//							 pc.SendPosSync()
		//						 }
		//					 }
		//					 if layer.Type == data.COLLISION {
		//						 flag := layer.Data[idx]
		//						 if flag != 0 {
		//							 // encounter a obscure
		//							 pc.State = player.STAND
		//							 handle.pos.X = handle.to.X
		//							 handle.pos.Y = handle.to.Y
		//							 pc.SendPosSync()
		//						 } else {
		//							 layer.Data[idx] = 1
		//						 }
		//					 }
		//				 }

		// 	handle.pos.X += vx
		// 	handle.pos.Y += vy
		// }

		// aoi update
		// m.aoi.Update(uint32(id), aoi.ModeWatcher|aoi.ModeMarker, aoi.FPoint(handle.pos))
		// }
	}
}

func (m *Zone) moveMonster() {
	for _, monster := range m.monsters {
		// if (monster.flag & flagDead) != 0 {
		// 	monster.reborn++
		// 	if monster.reborn >= 100 {
		// 		monster.flag = monster.flag &^ flagDead
		// 	}
		// 	continue
		// }
		//
		// if (monster.flag & flagActive) == 0 {
		// 	continue
		// }

		monster.HeartBeat(m)
	}
}

func (m *Zone) HeartBeat() {
	m.movePC()
	m.moveMonster()

	// m.aoi.Message(func(watcher uint32, marker uint32) {
	// 	// watcher is a player
	// 	if (watcher & maskNPC) == 0 {
	// 		handle := &m.players[watcher]
	// 		if handle.pc != nil {
	// 			handle.aoi <- marker
	// 		}
	// 	}
	//
	// 	// marker is a monster
	// 	if (marker & maskNPC) != 0 {
	// 		id := marker &^ maskNPC
	// 		monster := &m.monsters[id]
	// 		if (monster.flag & flagActive) == 0 {
	// 			monster.flag |= flagActive
	// 			monster.target = watcher
	// 		}
	// 	}
	// })
}

func (m *Zone) Login(player *Player, x uint16, y uint16, a chan<- ObjectIDType, rd <-chan interface{}, wr chan<- interface{}) error {
	var handle Handle
	handle.pc = player
	handle.pos.X = float32(x)
	handle.pos.Y = float32(y)
	handle.to = handle.pos
	handle.read = rd
	handle.write = wr

	idx := len(m.players)
	m.players = append(m.players, handle)

	player.Id = uint32(idx)
	// player.Scene = m

	// m.aoi.Update(player.Id, aoi.ModeWatcher|aoi.ModeMarker, aoi.FPoint(pos))

	return nil
}
