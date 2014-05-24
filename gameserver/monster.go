package main

import (
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet/darkeden"
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

	Enemies []uint32
	Owner *Player

	isEventMonster bool
	isChief        bool
	isMaster       bool
	bTreasure      bool

	LastKiller uint32
}

func (m *Monster) MaxHP() uint16 {
	mi := data.MonsterType2MonsterInfo[m.MonsterType]
	return m.STR*4 + uint16(mi.Level)
}

// a state machine
func (m *Monster) HeartBeat(mp *Scene) {
	m.ticker++
	if m.ticker == 100 {
		m.ticker = 0
		targetID := m.Enemies[0]
		targetObj := mp.objects[targetID]
		mi := data.MonsterType2MonsterInfo[m.MonsterType]
		if _, ok := targetObj.(*Player); ok {
			pc := targetObj.(*Player)
			x := m.X()
			y := m.Y()
			dx := int(pc.X()) - int(x)
			dy := int(pc.Y()) - int(y)
			d2 := dx*dx + dy*dy
			if d2 <= mi.MeleeRange*mi.MeleeRange {
				// attack player
			} else if d2 > mi.MissileRange*mi.MissileRange {
				m.flag = m.flag &^ flagActive
				m.Enemies = m.Enemies[:0]
			} else {
				dir := dir(dx, dy)
				x = uint8(int(x) + dirMoveMask[dir].X)
				y = uint8(int(y) + dirMoveMask[dir].Y)
				mp.Update(m.Entity, x, y)
				pc.send <- darkeden.GCMovePacket{
					ObjectID: m.Id(),
					X:        uint8(x),
					Y:        uint8(y),
					Dir:      dir,
				}
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
