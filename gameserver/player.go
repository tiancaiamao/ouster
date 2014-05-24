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
	"os"
	"encoding/json"
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
	level uint8

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
	
	computation chan func()
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
		computation: make(chan func()),
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

func LoadPlayer(name string) *darkeden.GCUpdateInfoPacket {	
	info := &darkeden.GCUpdateInfoPacket{
		PCType: 'V',
		PCInfo: darkeden.PCInfo{
			Name:             name,
			Level:            150,
			Sex:              0,
			SkinColor:        420,
			Alignment:        7500,
			STR:              [3]uint16{20, 20, 20},
			DEX:              [3]uint16{20, 20, 20},
			INT:              [3]uint16{20, 20, 20},
			HP:               [2]uint16{472, 472},
			Rank:             50,
			RankExp:          10700,
			Exp:              125,
			Fame:             282,
			Sight:            13,
			Bonus:            9999,
			Competence:       1,
			GuildMemberRank:  4,
			AdvancementLevel: 100,
		},
		ZoneID: 21,
		ZoneX:  145,
		ZoneY:  237,
	}
	
	f, err := os.Open(os.Getenv("HOME")+"/.ouster/player/"+name)
	if err != nil {
		return info
	}
	
	decoder := json.NewDecoder(f)
	var ret darkeden.GCUpdateInfoPacket
	err = decoder.Decode(&ret)
	if err != nil {
		return info
	}
	
	return &ret
}

func (player *Player) handleClientMessage(pkt packet.Packet) {
	switch pkt.Id() {
	case darkeden.PACKET_CG_CONNECT:
		raw := pkt.(*darkeden.CGConnectPacket)
		info := LoadPlayer(raw.PCName)
		info.PCInfo.ObjectID = player.Id()
		player.send <- info
		player.send <- &darkeden.GCPetInfoPacket{}
	case darkeden.PACKET_CG_READY:
		log.Println("get a CG Ready Packet!!!")
		player.send <- &darkeden.GCSetPositionPacket{
			X:   145,
			Y:   237,
			Dir: 2,
		}
		player.send <- &darkeden.GCSkillInfoPacket{
			PCType: darkeden.PC_VAMPIRE,
			PCSkillInfoList: []darkeden.SkillInfo{
				darkeden.VampireSkillInfo{
					LearnNewSkill: false,
					SubVampireSkillInfoList: []darkeden.SubVampireSkillInfo{
						darkeden.SubVampireSkillInfo{
							SkillType:   darkeden.SKILL_RAPID_GLIDING,
							Interval:    50,
							CastingTime: 31,
						},
						darkeden.SubVampireSkillInfo{
							SkillType:   darkeden.SKILL_METEOR_STRIKE,
							Interval:    10,
							CastingTime: 4160749567,
						},
						darkeden.SubVampireSkillInfo{
							SkillType:   darkeden.SKILL_INVISIBILITY,
							Interval:    30,
							CastingTime: 11,
						},
						darkeden.SubVampireSkillInfo{
							SkillType:   darkeden.SKILL_PARALYZE,
							Interval:    60,
							CastingTime: 41,
						},
						darkeden.SubVampireSkillInfo{
							SkillType:   darkeden.SKILL_BLOOD_SPEAR,
							Interval:    60,
							CastingTime: 41,
						},
					},
				},
			},
		}
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
	case darkeden.PACKET_CG_SKILL_TO_SELF:
		skill := pkt.(darkeden.CGSkillToSelfPacket)
		switch skill.SkillType {
		case darkeden.SKILL_INVISIBILITY:
			ok := &darkeden.GCSkillToSelfOK1{
				SkillType: darkeden.SKILL_INVISIBILITY,
				CEffectID: 181,
				Duration:  0,
				Grade:     0,
			}
			ok.Short = make(map[darkeden.ModifyType]uint16)
			ok.Short[12] = 180 + 256
			player.send <- ok
		}
	case darkeden.PACKET_CG_SKILL_TO_OBJECT:
		skill := pkt.(darkeden.CGSkillToObjectPacket)
		player.SkillToObject(skill)
	case darkeden.PACKET_CG_SKILL_TO_TILE:
		skill := pkt.(darkeden.CGSkillToTilePacket)
		switch skill.SkillType {
		case darkeden.SKILL_RAPID_GLIDING:
			fastMove := &darkeden.GCFastMovePacket{
				ObjectID:  player.Id(),
				FromX:     player.X(),
				FromY:     player.Y(),
				ToX:       skill.X,
				ToY:       skill.Y,
				SkillType: skill.SkillType,
			}
			player.send <- fastMove
			ok := &darkeden.GCSkillToTileOK1{
				SkillType: skill.SkillType,
				CEffectID: skill.CEffectID,
				Duration:  10,
				Range:     1,
				X:         skill.X,
				Y:         skill.Y,
			}
			player.send <- ok
		}

	case darkeden.PACKET_CG_BLOOD_DRAIN:
	case darkeden.PACKET_CG_VERIFY_TIME:
	case darkeden.PACKET_CG_LOGOUT:
		info := LoadPlayer(player.name)
		info.PCInfo.Level = player.level
		info.PCInfo.HP[0] = uint16(player.hp)	
		f, err := os.Create(os.Getenv("HOME")+"/.ouster/player/"+player.name)
		if err != nil {
			return
		}
		
		encoder := json.NewEncoder(f)
		err = encoder.Encode(info)
		f.Close()
		return
	}
}

func (player *Player) SkillToObject(packet darkeden.CGSkillToObjectPacket) {
	target := player.Scene.objects[packet.TargetObjectID]
	if monster, ok := target.(*Monster); ok {
		if skillExecutable, ok := skillP2M[packet.SkillType]; ok {
			if monster.Owner == player || monster.Owner == nil {
				skillExecutable(player, monster)
			} else {
					monster.Owner.computation <- func() {
						skillExecutable(player, monster)
					}
			}	
		} else {
			log.Println("can't execute skill ", packet.SkillType)
		}
	} else {
		
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
		case f, _:= <-this.computation:
			f()
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
