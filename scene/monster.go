package scene

import (
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/data"
)

const (
	flagDead = 1 << iota
	flagActive
)

type Monster struct {
	Id uint32

	pos   ouster.FPoint
	level int
	hp    int

	damage  int
	defence int
	tohit   int
	dodge   int

	// mask the monster's current status, flagDead means it's dead.
	// flagActive means it's active by player...
	flag uint8

	// if flag & flagDead, reborn after that times heartbeat
	reborn int
	// the config information that generate this monster
	meta interface{}
}

func (m *Monster) Init(meta interface{}) {
	switch meta.(type) {
	case data.Enemy:
		e := meta.(data.Enemy)
		m.pos.X = float32(e.Location.X)
		m.pos.Y = float32(e.Location.Y)
	case data.EnemyGroup:

	}

	m.meta = meta
}

func (m *Monster) HeartBeat() {

}
