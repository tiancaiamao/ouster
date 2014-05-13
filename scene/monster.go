package scene

import (
	// "github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/aoi"
	// "github.com/tiancaiamao/ouster/data"
	// "math"
)

const (
	flagDead = 1 << iota
	flagActive
)

type MonsterState uint8

const (
	stateIdle = iota
	state
)

type Monster struct {
	aoi *aoi.Entity

	// mask the monster's current status, flagDead means it's dead.
	// flagActive means it's active by player...
	flag uint8

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

	isEventMonster bool
	isChief        bool
	isMaster       bool
	bTreasure      bool

	LastKiller uint32
}

// a state machine
func (m *Monster) HeartBeat(mp *Zone) {
	// handle := mp.Player(m.target)
	// // pc := handle.pc
	// d := ouster.Distance2(m.pos, handle.pos)
	// if d < 10 {
	// 	// msg := packet.SkillTargetEffectPacket{
	// 	// 	Skill: 1,
	// 	// 	From:  m.Id,
	// 	// 	To:    m.target,
	// 	// 	Hurt:  1,
	// 	// 	Succ:  true,
	// 	// }
	// 	// nearby := pc.NearBy()
	// 	// boardcast(nearby, msg, mp)
	// } else {
	// 	dx := handle.pos.X - m.pos.X
	// 	dy := handle.pos.Y - m.pos.Y
	// 	angle := math.Atan2(float64(dy), float64(dx))
	// 	vx := m.speed * float32(math.Cos(angle))
	// 	vy := m.speed * float32(math.Sin(angle))
	//
	// 	m.pos.X += vx
	// 	m.pos.Y += vy
	// }
}

func boardcast(nearby []uint32, msg interface{}, mp *Zone) {
	for _, playerId := range nearby {
		p := mp.Player(playerId)
		if p != nil {
			p.write <- msg
		}
	}
}
