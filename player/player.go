package player

import (
	"github.com/tiancaiamao/ouster"
	// "github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/packet/darkeden"
	// "github.com/tiancaiamao/ouster/skill"
	"log"
	"net"
	"time"
)

type PlayerClass uint8

const (
	_     = iota
	BRUTE = iota
)

type PlayerState uint8

const (
	STAND PlayerState = iota
	MOVE
)

// package scene have import player, if player import scene it would be a circle
// so use a interface to avoid direct use of scene.
type scene interface {
	Pos(uint32) (ouster.FPoint, error)
	To(uint32) (ouster.FPoint, error)
	Creature(uint32) ouster.Creature
	String() string
}

const (
	LEFT      = 53
	RIGHT     = 49
	UP        = 34
	DOWN      = 55
	LEFTUP    = 50
	RIGHTUP   = 48
	LEFTDOWN  = 52
	RIGHTDOWN = 54
)

// mostly the same as data.Player, but this is in memory instead.
type Player struct {
	Id    uint32 // set by scene.Login
	Scene scene  // set by scene.Login

	name  string
	class PlayerClass
	hp    int
	mp    int
	speed float32
	level int

	strength     int
	agility      int
	intelligence int

	carried []int

	conn   net.Conn
	client <-chan packet.Packet
	send   chan<- packet.Packet
	aoi    <-chan uint32

	read      <-chan interface{}
	write     chan<- interface{}
	nearby    []uint32
	heartbeat <-chan time.Time
	ticker    uint32

	// Own by scene...write allowed only by scene agent
	State PlayerState
}

// implement Create
func (player *Player) Agility() int {
	return player.agility
}

func (player *Player) Strength() int {
	return player.strength
}

func (player *Player) Intelligence() int {
	return player.intelligence
}

func (player *Player) Damage() int {
	return player.strength
}

func (player *Player) Dodge() int {
	return player.agility
}

func (player *Player) ToHit() int {
	return player.agility
}

func (player *Player) HP() int {
	return 4*player.strength + player.level
}

// provide for scene to use
func (player *Player) Speed() float32 {
	return player.speed
}

func (player *Player) Defense() int {
	return player.strength
}

func New(conn net.Conn, a <-chan uint32, rd <-chan interface{}, wr chan<- interface{}) *Player {
	return &Player{
		name:  "test",
		class: 1,
		hp:    110,
		mp:    110,
		conn:  conn,
		speed: 0.5,

		aoi:   a,
		read:  rd,
		write: wr,

		heartbeat: time.Tick(50 * time.Millisecond),
	}
}

func (player *Player) NearBy() []uint32 {
	return player.nearby
}

func (player *Player) handleClientMessage(pkt packet.Packet) {
	var flag bool
	hp := darkeden.GCStatusCurrentHP{
		ObjectID:  2351,
		CurrentHP: 133,
	}

	switch pkt.Id() {
	case darkeden.PACKET_CG_CONNECT:
		player.send <- &darkeden.GCUpdateInfoPacket{}
	case darkeden.PACKET_CG_READY:
		log.Println("get a CG Ready Packet!!!")
		player.send <- &darkeden.GCSetPositionPacket{
			X:   145,
			Y:   237,
			Dir: 2,
		}
	case darkeden.PACKET_CG_MOVE:
		player.write <- pkt
		move := pkt.(darkeden.CGMovePacket)
		moveOk := darkeden.GCMoveOKPacket{
			Dir: move.Dir,
			X:   move.X,
			Y:   move.Y,
		}
		player.send <- moveOk

		if !flag {
			flag = true
			addMonster := &darkeden.GCAddMonster{
				ObjectID:    hp.ObjectID,
				MonsterType: 223,
				MonsterName: "test",
				MainColor:   7,
				SubColor:    174,
				X:           146,
				Y:           238,
				Dir:         2,
				CurrentHP:   133,
				MaxHP:       133,
			}
			addBat := &darkeden.GCAddBat{
				ObjectID:    2352,
				MonsterName: "bat",
				X:           149,
				Y:           242,
				Dir:         1,
				CurrentHP:   111,
				MaxHP:       133,
				GuildID:     1,
			}
			player.send <- addBat
			player.send <- addMonster
		}

	case darkeden.PACKET_CG_ATTACK:
		if hp.CurrentHP > 0 {
			hp.CurrentHP -= 5
			player.send <- hp
			player.send <- darkeden.GCAttackMeleeOK1(hp.ObjectID)
		}

	case darkeden.PACKET_CG_BLOOD_DRAIN:
	case darkeden.PACKET_CG_VERIFY_TIME:
	}
}

type BaseAttack struct{}

func (_ BaseAttack) ExecuteTarget(from, to ouster.Creature) (int, bool) {
	return 10, true
}

type SkillEffect struct {
	Id   int
	To   uint32
	Succ bool
	Hurt int
}

// func (player *Player) execute(pkt packet.SkillPacket) {
// 	skl := skill.Query(pkt.Id)
// 	switch skl.(type) {
// 	case skill.SelfSkill:
//
// 	case skill.TargetSkill:
// 		skill := skl.(skill.TargetSkill)
// 		target := player.Scene.Creature(pkt.Target)
// 		hurt, ok := skill.ExecuteTarget(player, target)
//
// 		player.write <- SkillEffect{
// 			Id:   pkt.Id,
// 			To:   pkt.Target,
// 			Succ: ok,
// 			Hurt: hurt,
// 		}
// 	case skill.RegionSkill:
// 	}
// }

type CMovePacketAck struct{}

func (this *Player) handleSceneMessage(msg interface{}) {
}

func (this *Player) heartBeat() {
	this.ticker++
}
