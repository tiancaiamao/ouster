package player

import (
	"github.com/tiancaiamao/ouster/packet"
)

func (this *Player) loop() {
	var msg interface{}
	for {
		select {
		case msg = <-this.client:
			this.handleClientMessage(msg)
		case <-this.Scene2player:
			this.handleSceneMessage(msg)
		case <-this.aoi:
			// 来自aoi的消息
		}
	}
}

func (this *Player) handleSceneMessage(msg interface{}) {
	switch msg.(type) {
	case packet.CMovePacket:
		this.send <- packet.Packet{packet.PCMove, msg}
	}
}

func (player *Player) Go() {
	go player.loop()

	ch := make(chan packet.Packet)
	player.send = ch

	for {
		pkt := <-ch
		err := packet.Write(player.conn, pkt.Id, pkt.Obj)
		if err != nil {
			continue
		}
	}
}
