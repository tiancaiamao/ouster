package player

import (
	"bytes"
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/packet/darkeden"
	"log"
)

func (this *Player) loop() {
	var msg interface{}
	for {
		select {
		case msg, ok := <-this.client:
			if !ok {
				// kick the player off...
				return
			} else {
				this.handleClientMessage(msg)
			}
		case msg = <-this.read:
			this.handleSceneMessage(msg)
		case id := <-this.aoi:
			this.nearby = append(this.nearby, id)
		case <-this.heartbeat:
			this.heartBeat()
		}
	}
}

func (player *Player) Go() {
	read := make(chan packet.Packet, 1)
	write := make(chan packet.Packet, 1)
	player.send = write
	player.client = read

	// open a goroutine to read from conn
	go func() {
		reader := darkeden.NewReader()
		for {
			data, err := reader.Read(player.conn)
			if err != nil {
				log.Println(err)
				player.conn.Close()
				close(read)
				return
			}
			// log.Println("packet before send to chan", data)
			read <- data
			// log.Println("packet after send to chan", data)
		}
	}()

	// open a goroutine to write to conn
	go func() {
		writer := darkeden.NewWriter()
		for {
			pkt := <-write
			log.Println("write channel get a pkt ", pkt.String())
			err := writer.Write(player.conn, pkt)
			if err != nil {
				log.Println(err)
				continue
			}

			buf := &bytes.Buffer{}
			writer.Write(buf, pkt)
			log.Println("send packet to client: ", buf.Bytes())
		}
	}()

	player.loop()
}
