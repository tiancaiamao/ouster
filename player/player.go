package player

import (
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
	"net"
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

// mostly the same as data.Player, but this is in memory instead.
type Player struct {
	Id    uint32 // set by scene.Login
	Scene string // set by scene.Login
	name  string
	class PlayerClass
	hp    int
	mp    int
	speed float32

	carried []int

	conn         net.Conn
	client       <-chan interface{} // actually, a XXXPacket struct
	send         chan<- packet.Packet
	aoi          <-chan interface{}
	Scene2player chan interface{} // alloc in player.New
	Player2scene chan interface{} // alloc in player.New
	nearby       []uint32

	// Own by scene...write allowed only by scene agent
	Pos   ouster.FPoint
	State PlayerState
	To    ouster.FPoint
}

// provide for scene to use
func (player *Player) Speed() float32 {
	return player.speed
}

func New(playerData *data.Player, conn net.Conn) *Player {
	return &Player{
		name:         playerData.Name,
		class:        PlayerClass(playerData.Class),
		hp:           playerData.HP,
		mp:           playerData.MP,
		carried:      playerData.Carried,
		conn:         conn,
		Scene2player: make(chan interface{}),
		Player2scene: make(chan interface{}),
	}
}

func (player *Player) NearBy() []uint32 {
	return player.nearby
}

func (this *Player) handleClientMessage(msg interface{}) {
	switch msg.(type) {
	case packet.CMovePacket:
		move := msg.(packet.CMovePacket)
		this.Player2scene <- move
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
				info["pos"] = this.Pos
			case "scene":
				info["scene"] = this.Scene
			}
		}
		this.send <- packet.Packet{packet.PPlayerInfo, info}
	}
}

type CMovePacketAck struct{}

func (this *Player) handleSceneMessage(msg interface{}) {
	switch msg.(type) {
	case CMovePacketAck:
		posSync := packet.PosSyncPacket{
			Cur: this.Pos,
			To:  this.To,
		}
		pkt := packet.Packet{
			Id:  packet.PPosSync,
			Obj: posSync,
		}
		this.send <- pkt
	}
}
