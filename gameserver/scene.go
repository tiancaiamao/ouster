package main

import (
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/aoi/cell"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet/darkeden"
	"log"
	"math/rand"
	"time"
)

type Object interface {
	Id() uint32
}

type Scene struct {
	*data.Map
	objects  []Object
	players  []*Player
	monsters []Monster
	aoi.Aoi

	quit      chan struct{}
	event     chan interface{}
	heartbeat <-chan time.Time
}

func (scene *Scene) AddObject(obj Object) uint32 {
	idx := len(scene.objects)
	if idx < 10000 {
		scene.objects = append(scene.objects, obj)
	} else {
	}
	return uint32(idx)
}

func NewScene(m *data.Map) *Scene {
	ret := new(Scene)
	ret.Map = m
	aoi := cell.New(m.Width, m.Height, 32, 32)
	players := make([]*Player, 0, 200)

	num := 0
	for _, mi := range m.MonsterInfo {
		num += int(mi.Count)
	}
	monsters := make([]Monster, num)
	ret.objects = make([]Object, 0, num+200)

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
			id := ret.AddObject(monster)
			monster.Entity = aoi.Add(uint8(x), uint8(y), id)
			monster.MonsterType = mi.MonsterType
			monster.STR = tp.STR
			monster.DEX = tp.DEX
			monster.INT = tp.INTE
			monster.Defense = monster.DEX / 2
			monster.Protection = monster.STR
			monster.HP = tp.STR*4 + uint16(tp.Level)
			idx++
		}
	}

	ret.Aoi = aoi
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

func (m *Scene) String() string {
	return m.Map.Name
}

func (m *Scene) HeartBeat() {
	m.Message(func(watcher aoi.Entity, marker aoi.Entity) {
		wId := watcher.Id()
		mId := marker.Id()

		wObj := m.objects[wId]
		mObj := m.objects[mId]

		if wp, ok := wObj.(*Player); ok {
			switch mObj.(type) {
			case *Monster:
				// monster active by player
				monster := mObj.(*Monster)
				monster.flag |= flagActive
				monster.Owner = wp
				monster.Enemies = append(monster.Enemies, wId)
			case *Player:
				player := mObj.(*Player)
				player.handleAoiMessage(wId)
			}
		}

		if _, ok := wObj.(*Monster); ok {
			if _, ok2 := mObj.(*Player); ok2 {
				player := mObj.(*Player)
				player.handleAoiMessage(wId)
			}
		}

		for i := 0; i < len(m.monsters); i++ {
			monster := &m.monsters[i]
			if (monster.flag & flagActive) != 0 {
				monster.HeartBeat(m)
			}
		}
		return
	})
}

func (m *Scene) Login(player *Player) error {
	m.players = append(m.players, player)

	id := m.AddObject(player)
	player.Entity = m.Add(145, 237, id)
	player.Scene = m

	return nil
}

func loop(m *Scene) {
	for {
		for _, player := range m.players {
			select {
			case msg := <-player.agent2scene:
				m.processPlayerInput(player.Id(), msg)
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
		log.Println("scene receive a CGMovePacket:", move.X, move.Y, move.Dir)
		obj := m.objects[playerId]
		player := obj.(*Player)

		if move.Dir >= 8 {
			moveErr := darkeden.GCMoveErrorPacket{
				player.X(),
				player.Y(),
			}
			player.send <- moveErr
		}

		move.X = uint8(int(move.X) + dirMoveMask[move.Dir].X)
		move.Y = uint8(int(move.Y) + dirMoveMask[move.Dir].Y)

		m.Update(player.Entity, move.X, move.Y)
		player.send <- darkeden.GCMoveOKPacket{
			X:   move.X,
			Y:   move.Y,
			Dir: move.Dir,
		}
	}
}

var (
	maps map[string]*Scene
)

func Initialize() {
	maps = make(map[string]*Scene)
	maps["limbo_lair_se"] = NewScene(&data.LimboLairSE)

	for _, m := range maps {
		m.Go()
	}
}

func Query(mapName string) *Scene {
	return maps[mapName]
}
