package scene

import (
	"github.com/tiancaiamao/ouster"
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
}
