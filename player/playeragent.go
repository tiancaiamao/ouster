package player

import (
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
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
	client chan interface{}
	aoi    <-chan interface{}
	scene  chan interface{}
	nearby []uint32
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
	}
}

func (this *Player) handleSceneMessage(msg interface{}) {
	switch msg.(type) {
	case packet.MovePacket:
		this.client <- packet.Packet{packet.PMove, msg}
	}
}

func (player *Player) Go() {
	go player.loop()
	go func(player *Player) {
		for {
			data := <-player.client
			pkt := data.(packet.Packet)
			err := packet.Write(player.conn, pkt.Id, pkt.Obj)
			if err != nil {
				continue
			}
		}
	}(player)
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
