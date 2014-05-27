package main

import (
	"bytes"
	"encoding/json"
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/packet/darkeden"
	"log"
	"math/rand"
	"net"
	"os"
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

const (
	ATTR_CURRENT = iota
	ATTR_MAX
	ATTR_BASE
)

type Creature struct {
	Level      uint8
	STR        [3]uint16
	DEX        [3]uint16
	INT        [3]uint16
	HP         [2]uint16
	MP         [2]uint16
	Defense    uint16
	Protection uint16
	ToHit      uint16
	Damage     uint16
}

type Player struct {
	aoi.Entity
	Creature

	PCType byte
	// field from data.PCInfo
	Name               string
	Sex                uint8
	BatColor           uint16
	SkinColor          uint16
	MasterEffectColor  uint8
	Alignment          uint32
	Rank               uint8
	RankExp            uint32
	Exp                uint32
	Fame               uint32
	Gold               uint32
	Sight              uint8
	Bonus              uint16
	HotKey             [8]uint16
	SilverDamage       uint16
	Competence         uint8
	GuildID            uint16
	GuildName          string
	GuildMemberRank    uint8
	UnionID            uint32
	AdvancementLevel   uint8
	AdvancementGoalExp uint32

	Scene *Scene

	carried []int

	conn   net.Conn
	packetReader *darkeden.Reader
	packetWriter *darkeden.Writer
	
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
		conn:        conn,
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

func (player *Player) Load(name string) error {
	info := data.PCInfo{
		PCType:           'V',
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
		ZoneID:           23,
		ZoneX:            145,
		ZoneY:            237,
	}

	var pcInfo data.PCInfo
	f, err := os.Open(os.Getenv("HOME") + "/.ouster/player/" + name)
	if err != nil {
		pcInfo = info
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&pcInfo)
	if err != nil {
		pcInfo = info
	}
	f.Close()

	player.PCType = pcInfo.PCType
	player.Name = pcInfo.Name
	player.Level = pcInfo.Level
	player.Sex = pcInfo.Sex
	player.SkinColor = pcInfo.SkinColor
	player.Alignment = pcInfo.Alignment
	player.STR = pcInfo.STR
	player.DEX = pcInfo.DEX
	player.INT = pcInfo.INT
	player.HP = pcInfo.HP
	player.Rank = pcInfo.Rank
	player.RankExp = pcInfo.RankExp
	player.Exp = pcInfo.Exp
	player.Fame = pcInfo.Fame
	player.Sight = pcInfo.Sight
	player.Bonus = pcInfo.Bonus
	player.Competence = pcInfo.Competence
	player.GuildMemberRank = pcInfo.GuildMemberRank
	player.AdvancementLevel = pcInfo.AdvancementLevel

	scene := zoneTable[pcInfo.ZoneID]
	scene.Login(player)
	scene.Update(player.Entity, pcInfo.ZoneX, pcInfo.ZoneY)
	return err
}

func (player *Player) Save() {
	info := player.PCInfo()
	f, err := os.Create(os.Getenv("HOME") + "/.ouster/player/" + player.Name)
	if err != nil {
		return
	}

	encoder := json.NewEncoder(f)
	err = encoder.Encode(info)
	f.Close()
}

func (player *Player) PCInfo() *data.PCInfo {
	return &data.PCInfo{
		ObjectID: player.Id(),
		Name:     player.Name,
		Level:    player.Level,
		Sex:      player.Sex,

		BatColor:          player.BatColor,
		SkinColor:         player.SkinColor,
		MasterEffectColor: player.MasterEffectColor,

		Alignment: player.Alignment,
		STR:       player.STR,
		DEX:       player.DEX,
		INT:       player.INT,

		HP: player.HP,

		Rank:    player.Rank,
		RankExp: player.RankExp,

		Exp:          player.Exp,
		Fame:         player.Fame,
		Gold:         player.Gold,
		Sight:        player.Sight,
		Bonus:        player.Bonus,
		HotKey:       player.HotKey,
		SilverDamage: player.SilverDamage,

		Competence: player.Competence,
		GuildID:    player.GuildID,

		GuildMemberRank: player.GuildMemberRank,
		UnionID:         player.UnionID,

		AdvancementLevel:   player.AdvancementLevel,
		AdvancementGoalExp: player.AdvancementGoalExp,
	}
}

func Encrypt(ZoneID uint16, ServerID uint16) uint8 {
	return uint8(((ZoneID >> 8) ^ ZoneID) ^ ((ServerID+1) << 4))
}

func (player *Player) handleClientMessage(pkt packet.Packet) {
	switch pkt.Id() {
	case darkeden.PACKET_CG_CONNECT:
		raw := pkt.(*darkeden.CGConnectPacket)
		player.Load(raw.PCName)

		info := &darkeden.GCUpdateInfoPacket{
			PCType: player.PCType,
			PCInfo: *player.PCInfo(),
			ZoneID: player.Scene.ZoneID,
			ZoneX:  player.X(),
			ZoneY:  player.Y(),
		}
		
		code := Encrypt(player.Scene.ZoneID, 1)
		player.packetReader.Code = code
		player.packetWriter.Code = code
		
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

				if monster.HP[ATTR_CURRENT] > damage {
					monster.HP[ATTR_CURRENT] -= damage
					player.send <- darkeden.GCStatusCurrentHP{
						ObjectID:  monster.Id(),
						CurrentHP: monster.HP[ATTR_CURRENT],
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
		player.Save()
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
				CurrentHP:   monster.HP[ATTR_CURRENT],
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
		case f, _ := <-this.computation:
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
		player.packetReader = reader
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
		player.packetWriter = writer
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
