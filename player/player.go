package player

import (
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/skill"
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
	Creature(int) ouster.Creature
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

	strength     int
	agility      int
	intelligence int

	carried []int

	conn   net.Conn
	client <-chan interface{} // actually, a XXXPacket struct
	send   chan<- packet.Packet
	aoi    <-chan uint32

	read      <-chan interface{} // alloc in player.New
	write     chan<- interface{} // alloc in player.New
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

// provide for scene to use
func (player *Player) Speed() float32 {
	return player.speed
}

func New(playerData *data.Player, conn net.Conn, a <-chan uint32, rd <-chan interface{}, wr chan<- interface{}) *Player {
	return &Player{
		name:    playerData.Name,
		class:   PlayerClass(playerData.Class),
		hp:      playerData.HP,
		mp:      playerData.MP,
		carried: playerData.Carried,
		conn:    conn,

		aoi:   a,
		read:  rd,
		write: wr,

		heartbeat: time.Tick(50 * time.Microsecond),
	}
}

func (player *Player) NearBy() []uint32 {
	return player.nearby
}

func (this *Player) handleClientMessage(msg interface{}) {
	switch msg.(type) {
	case packet.CMovePacket:
		move := msg.(packet.CMovePacket)
		this.write <- move
	case packet.PlayerInfoPacket:
		info := msg.(packet.PlayerInfoPacket)
		for k, _ := range info {
			switch k {
			case "name":
				info["name"] = this.name
			case "hp":
				info["hp"] = this.hp
			case "mp":
				info["mp"] = this.mp
			case "speed":
				info["speed"] = this.speed
			case "pos":
				info["pos"], _ = this.Scene.Pos(this.Id)
			case "scene":
				info["scene"] = this.Scene.String()
			}
		}
		this.send <- packet.Packet{packet.PPlayerInfo, info}
	case packet.SkillPacket:
		raw := msg.(packet.SkillPacket)
		this.execute(raw)
	}
}

type BaseAttack struct{}

func (_ BaseAttack) ExecuteTarget(from, to ouster.Creature) (int, bool) {
	return 10, true
}

type SkillEffect struct {
	Id   int
	To   int
	Succ bool
	Hurt int
}

func (player *Player) execute(pkt packet.SkillPacket) {
	skl := skill.Query(pkt.Id)
	switch skl.(type) {
	case skill.SelfSkill:

	case skill.TargetSkill:
		skill := skl.(skill.TargetSkill)
		target := player.Scene.Creature(pkt.Target)
		hurt, ok := skill.ExecuteTarget(player, target)

		player.write <- SkillEffect{
			Id:   pkt.Id,
			To:   pkt.Target,
			Succ: ok,
			Hurt: hurt,
		}
	case skill.RegionSkill:
	}
}

type CMovePacketAck struct{}

func (this *Player) handleSceneMessage(msg interface{}) {
	switch msg.(type) {
	case CMovePacketAck:
		sendPosSync(this)
	case packet.SMovePacket:
		pkt := packet.Packet{
			Id:  packet.PSMove,
			Obj: msg,
		}
		this.send <- pkt
	case packet.SkillTargetEffectPacket:
		pkt := packet.Packet{
			Id:  packet.PSkillTargetEffect,
			Obj: msg,
		}
		this.send <- pkt
	}
}

func sendPosSync(this *Player) {
	var posSync packet.PosSyncPacket
	posSync.Cur, _ = this.Scene.Pos(this.Id)
	posSync.To, _ = this.Scene.To(this.Id)

	pkt := packet.Packet{
		Id:  packet.PPosSync,
		Obj: posSync,
	}
	this.send <- pkt
}

func (this *Player) heartBeat() {
	this.ticker++

	// send PosSync every 400 ms
	if this.State == MOVE && (this.ticker&8) == 0 {
		sendPosSync(this)
	}
}
