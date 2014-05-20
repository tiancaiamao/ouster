package main

import (
	"bytes"
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/packet/darkeden"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	LEFT      = 0
	RIGHT     = 4
	UP        = 6
	DOWN      = 2
	LEFTUP    = 7
	RIGHTUP   = 5
	LEFTDOWN  = 1
	RIGHTDOWN = 3
)

type Point struct {
	X int
	Y int
}

var dirMoveMask [8]Point

func init() {
	dirMoveMask[RIGHTUP] = Point{1, -1}
	dirMoveMask[LEFT] = Point{-1, 0}
	dirMoveMask[RIGHT] = Point{1, 0}
	dirMoveMask[LEFTDOWN] = Point{-1, 1}
	dirMoveMask[DOWN] = Point{0, 1}
	dirMoveMask[RIGHTDOWN] = Point{1, 1}
	dirMoveMask[UP] = Point{0, -1}
	dirMoveMask[LEFTUP] = Point{-1, -1}
}

type Player struct {
	aoi.Entity
	Scene *Scene

	name  string
	hp    int
	mp    int
	speed float32
	level int

	STR        uint16
	DEX        uint16
	INT        uint16
	Defense    uint16
	Protection uint16
	ToHit      uint16
	Damage     uint16

	carried []int

	conn   net.Conn
	client <-chan packet.Packet
	send   chan<- packet.Packet

	agent2scene chan interface{}
	nearby      map[uint32]struct{}
	heartbeat   <-chan time.Time
	ticker      uint32
}

func NewPlayer(conn net.Conn) *Player {
	return &Player{
		name:       "test",
		hp:         110,
		mp:         110,
		conn:       conn,
		speed:      0.5,
		STR:        20,
		DEX:        20,
		INT:        20,
		Defense:    10,
		Protection: 20,
		ToHit:      30,
		Damage:     25,

		agent2scene: make(chan interface{}),
		nearby:      make(map[uint32]struct{}),
		heartbeat:   time.Tick(50 * time.Millisecond),
	}
}

func (player *Player) NearBy() map[uint32]struct{} {
	return player.nearby
}

// if tohit == dodge, the default formula is 0.85
// if tohit < dodge, then tohit / dodge should be primary factor, also take other factor into consideration
// if tohit > dodge, then the differential should be important, also dodge.
func HitTest(tohit uint16, dodge uint16) bool {
	var prob float32
	if tohit < dodge {
		prob = 0.85*float32(tohit)/float32(dodge) - 0.15*float32(dodge-tohit)/float32(tohit)
	} else {
		prob = 0.85 + 0.15*float32(tohit-dodge)/float32(dodge)
	}

	return rand.Float32() < prob
}

func (player *Player) handleClientMessage(pkt packet.Packet) {
	switch pkt.Id() {
	case darkeden.PACKET_CG_CONNECT:
		player.send <- &darkeden.GCUpdateInfoPacket{}
		player.send <- &darkeden.GCPetInfoPacket{}
	case darkeden.PACKET_CG_READY:
		log.Println("get a CG Ready Packet!!!")
		player.send <- &darkeden.GCSetPositionPacket{
			X:   145,
			Y:   237,
			Dir: 2,
		}

		//		player.send <- &darkeden.GCAddMonster {
		//			ObjectID: 127,
		//			MonsterType: 7,
		//			MonsterName: "test",
		//			X: 150,
		//			Y: 240,
		//			Dir: 3,
		//			CurrentHP: 77,
		//			MaxHP: 120,
		//		}
	case darkeden.PACKET_CG_MOVE:
		player.agent2scene <- pkt
	case darkeden.PACKET_CG_SAY:
		say := pkt.(*darkeden.CGSayPacket)
		log.Println("say:", say.Message)
	case darkeden.PACKET_CG_ATTACK:
		attack := pkt.(darkeden.CGAttackPacket)
		log.Println(" attack monster ", attack.ObjectID)
		target := player.Scene.objects[attack.ObjectID]
		if monster, ok := target.(*Monster); ok {
			hit := HitTest(player.ToHit, monster.Defense)
			if hit {
				player.send <- darkeden.GCAttackMeleeOK1{
					ObjectID: monster.Id(),
				}
				damage := uint16(1)
				if player.Damage > monster.Protection {
					damage = player.Damage - monster.Protection
				}

				if monster.HP > damage {
					monster.HP -= damage
					player.send <- darkeden.GCStatusCurrentHP{
						ObjectID:  monster.Id(),
						CurrentHP: monster.HP,
					}
				} else {
					player.send <- &darkeden.GCAddMonsterCorpse{
						ObjectID:    monster.Id(),
						MonsterType: monster.MonsterType,
						MonsterName: monster.Name,
						X:           monster.X(),
						Y:           monster.Y(),
						Dir:         2,
						LastKiller:  player.Id(),
					}
					player.send <- darkeden.GCCreatureDiedPacket(monster.Id())
				}
			} else {
			}
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

// called in scene
func (this *Player) handleAoiMessage(id uint32) {
	obj := this.Scene.objects[id]
	if _, ok := obj.(*Monster); ok {
		log.Println("it's a monster...send message")
		monster := obj.(*Monster)
		if _, ok := this.nearby[id]; !ok {
			this.nearby[id] = struct{}{}

			addMonster := &darkeden.GCAddMonster{
				ObjectID:    uint32(id),
				MonsterType: monster.MonsterType,
				MonsterName: "test",
				X:           monster.X(),
				Y:           monster.Y(),
				Dir:         2,
				CurrentHP:   monster.HP,
				MaxHP:       monster.MaxHP(),
			}

			this.send <- addMonster
			monster.flag |= flagActive
			log.Println("monster ", id, "set to active", monster.flag)
			monster.Enemies = append(monster.Enemies, this.Id())
		} else {

		}
	}
}

func (this *Player) heartBeat() {
	this.ticker++
}

func (this *Player) loop() {
	// var msg interface{}
	for {
		select {
		case msg, ok := <-this.client:
			if !ok {
				// kick the player off...
				return
			} else {
				this.handleClientMessage(msg)
			}
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
			read <- data
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
