package player

import (
	"github.com/tiancaiamao/ouster"
	// "github.com/tiancaiamao/ouster/data"
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
	client <-chan darkeden.Packet
	send   chan<- darkeden.Packet
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

func (player *Player) handleClientMessage(pkt darkeden.Packet) {
	switch pkt.Id() {
	case darkeden.PACKET_CG_CONNECT:
		player.send <- &darkeden.GCUpdateInfoPacket{}
	case darkeden.PACKET_CG_READY:
		log.Println("get a CG Ready Packet!!!")
		// don't what's this packet
		//	player.conn.Write([]byte{85, 1, 3, 0, 0, 0, 2, 145, 237, 2, 62, 1, 4, 0, 0, 0, 3, 0, 0, 0, 0, 105, 1, 3, 0, 0, 0, 4, 1, 1, 0, 0, 61, 1, 1, 0, 0, 0, 5, 0, 19, 1, 17, 0, 0, 0, 6, 5, 18, 40, 0, 19, 50, 0, 20, 80, 0, 16, 33, 0, 17, 35, 0, 0, 173, 0, 26, 0, 0, 0, 7, 80, 48, 0, 0, 8, 194, 179, 181, 199, 182, 224, 183, 242, 0, 0, 150, 235, 1, 174, 0, 174, 0, 1, 0, 0, 0, 34, 1, 6, 0, 0, 0, 8, 31, 0, 100, 0, 0, 0, 32, 1, 15, 0, 0, 0, 9, 3, 0, 0, 5, 0, 1, 0, 1, 30, 0, 2, 0, 1, 30, 0, 249, 0, 1, 0, 0, 0, 10, 0, 248, 0, 1, 0, 0, 0, 11, 0, 31, 1, 24, 0, 0, 0, 12, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 185, 0, 27, 0, 0, 0, 48, 176, 47, 0, 0, 73, 0, 8, 200, 248, 182, 224, 210, 193, 182})
		player.conn.Write([]byte{85, 1, 3, 0, 0, 0, 2, 145, 237, 2})
	case darkeden.PACKET_CG_MOVE:
		player.write <- pkt
	case darkeden.PACKET_CG_ATTACK:
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
