package player

import (
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/scene"
	"net"
)

type PlayerClass uint8

const (
	_     = iota
	BRUTE = iota
)

// mostly the same as data.Player, but this is in memory instead.
type Player struct {
	id    uint32 // alloc by scene
	name  string
	class PlayerClass
	hp    int
	mp    int
	speed float32

	carried []int

	conn   net.Conn
	client <-chan interface{}
	send   chan<- packet.Packet
	aoi    <-chan interface{}
	scene  chan interface{}
	nearby []uint32
	m      *scene.Map
}

func (player *Player) Speed() float32 {
	return player.speed
}

func New(playerId uint32, playerData *data.Player, conn net.Conn, scene chan interface{}) *Player {
	return &Player{
		id:      playerId,
		name:    playerData.Name,
		class:   PlayerClass(playerData.Class),
		hp:      playerData.HP,
		mp:      playerData.MP,
		carried: playerData.Carried,
		conn:    conn,
	}
}

func (player *Player) Init(playerId uint32, playerData *data.Player, conn net.Conn, scene chan interface{}, m *scene.Map) {
	player.id = playerId
	player.name = playerData.Name
	player.class = PlayerClass(playerData.Class)
	player.hp = playerData.HP
	player.mp = playerData.MP
	player.carried = playerData.Carried
	player.conn = conn
	player.scene = scene
	player.m = m
}

func (this *Player) loop() {
	var msg interface{}
	for {
		select {
		case msg = <-this.client:
			this.handleClientMessage(msg)
		case <-this.scene:
			this.handleSceneMessage(msg)
		case <-this.aoi:
			// 来自aoi的消息
		}
	}
}

func (this *Player) handleClientMessage(msg interface{}) {
	switch msg.(type) {
	case packet.MovePacket:
		move := msg.(packet.MovePacket)
		this.scene <- move
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
				info["pos"], _ = this.m.PlayerPosition(this.id)
			}
		}
		this.send <- packet.Packet{packet.PPlayerInfoPacket, info}
	}
}

func (this *Player) handleSceneMessage(msg interface{}) {
	switch msg.(type) {
	case packet.MovePacket:
		this.send <- packet.Packet{packet.PMove, msg}
	}
}

func (player *Player) Go() {
	go player.loop()
	
	ch := make(chan packet.Packet)
	player.send = ch
	go func(player *Player, ch packet.Packet) {
		for {
			data := <-ch
			err := packet.Write(player.conn, pkt.Id, pkt.Obj)
			if err != nil {
				continue
			}
		}
	}(player, ch)
	for {
		data, err := packet.Read(player.conn)
		if err != nil {
			// write a reset packet...
			continue
		}
		player.client <- data
	}
}

func (player *Player) NearBy() []uint32 {
	return player.nearby
}
