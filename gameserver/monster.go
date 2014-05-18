package main

import (
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
//	"github.com/tiancaiamao/ouster/packet/darkeden"
)

const (
	flagDead = 1 << iota
	flagActive
)

type Monster struct {
	aoi.Entity

	// mask the monster's current status, flagDead means it's dead.
	// flagActive means it's active by player...
	flag uint8

	ticker uint16

	MonsterType uint16
	Name        string

	STR          uint16
	DEX          uint16
	INT          uint16
	HP           uint16
	Defense      uint16
	Protection   uint16
	ToHit        uint16
	Damage       uint16
	MeleeRange   int
	MissileRange int

	Enemies []ObjectIDType

	isEventMonster bool
	isChief        bool
	isMaster       bool
	bTreasure      bool

	LastKiller uint32
}

// a state machine
func (m *Monster) HeartBeat(mp *Scene) {
	m.ticker++
	if m.ticker == 200 {
		m.ticker = 0
		targetID := m.Enemies[0]
		mi := data.MonsterType2MonsterInfo[m.MonsterType]
		if targetID.Player() {
			pc := mp.Player(targetID.Index())
			x := m.X()
			y := m.Y()
			dx := int(pc.X()) - int(x)
			dy := int(pc.Y()) - int(y)
			if dx*dx+dy*dy <= mi.MeleeRange*mi.MeleeRange {
				// attack player
			} else {
				switch {
				case dx > 0:
					x++
				case dx < 0:
					x--
				}
				switch {
				case dy > 0:
					y++
				case dy < 0:
					y--
				}
				mp.Update(m.Entity, x, y)

//				pc.send <- darkeden.GCMovePacket{
//					ObjectID: m.Id(),
//					X:        uint8(x),
//					Y:        uint8(y),
//					Dir:      dir(dx, dy),
//				}
			}
		}
	}
}

func dir(dx int, dy int) uint8 {
	var ret uint8
	switch {
	case dx > 0 && dy > 0:
		ret = RIGHTDOWN
	case dx > 0 && dy == 0:
		ret = RIGHT
	case dx > 0 && dy < 0:
		ret = RIGHTUP
	case dx < 0 && dy > 0:
		ret = LEFTDOWN
	case dx < 0 && dy == 0:
		ret = LEFT
	case dx < 0 && dy > 0:
		ret = LEFTUP
	case dx == 0 && dy > 0:
		ret = DOWN
	case dx == 0 && dy < 0:
		ret = UP
	}
	return ret
}
