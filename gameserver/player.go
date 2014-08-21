package main

import (
    "bytes"
    "encoding/json"
    "github.com/tiancaiamao/ouster"
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/packet"
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
    X   int
    Y   int
}

const (
    ATTR_CURRENT = iota
    ATTR_MAX
    ATTR_BASE
)

type PlayerStatus uint8

const (
    GPS_NONE = iota
    GPS_BEGIN_SESSION
    GPS_WAITING_FOR_CG_READY
    GPS_NORMAL
    GPS_IGNORE_ALL
    GPS_AFTER_SENDING_GL_INCOMING_CONNECTION
    GPS_END_SESSION
)

// Player负责网络相关的处理，接收消息包，发送消息包
// 也负责玩家数据完全加载之前的封包的处理
// 直接加载完成之后，控制权转交给agent
type Player struct {
    PlayerStatus PlayerStatus
    // aoi.Entity
    // Creature

    // PCType byte
    // field from data.PCInfo
    // Name               string
    // Sex                uint8
    // BatColor           uint16
    // SkinColor          uint16
    // HairColor          uint16
    // MasterEffectColor  uint8
    // Alignment          uint32
    // Rank               uint8
    // RankExp            uint32
    // Exp                uint32
    // Fame               uint32
    // Gold               uint32
    // Sight              uint8
    // Bonus              uint16
    // HotKey             [8]uint16
    // SilverDamage       uint16
    // Competence         uint8
    // GuildID            uint16
    // GuildName          string
    // GuildMemberRank    uint8
    // UnionID            uint32
    // AdvancementLevel   uint8
    // AdvancementGoalExp uint32

    // Scene *Scene
    // carried []int
    // skillslot []SkillSlot

    conn         net.Conn
    packetReader *packet.Reader
    packetWriter *packet.Writer

    client <-chan packet.Packet
    send   chan<- packet.Packet

    nearby map[uint32]struct{}
}

// 两个后台channel，负责与客户端的通信，分别读和写
func (player *Player) Go() {
    read := make(chan packet.Packet, 1)
    write := make(chan packet.Packet, 1)
    player.send = write
    player.client = read

    // open a goroutine to read from conn
    go func() {
        reader := packet.NewReader()
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
        writer := packet.NewWriter()
        player.packetWriter = writer
        for {
            pkt := <-write
            // log.Println("write channel get a pkt ", pkt.String())
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
}

type SkillSlot struct {
    SkillType uint16
    ExpLevel  uint16

    LastUse  time.Time
    Cooling  uint16
    Duration uint16

    Interval    uint32
    CastingTime uint32
}

// func (player *Player) SkillSlot(SkillType uint16) *SkillSlot {
//     for i := 0; i < len(player.skillslot); i++ {
//         if player.skillslot[i].SkillType == SkillType {
//             return &player.skillslot[i]
//         }
//     }
//     return nil
// }

func NewPlayer(conn net.Conn) *Player {
    ret := &Player{
        conn:   conn,
        nearby: make(map[uint32]struct{}),
    }
    go ret.Go()
    return ret
}

func (player *Player) NearBy() map[uint32]struct{} {
    return player.nearby
}

// TODO
func (player *Player) sendPacket(pkt packet.Packet) {

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

func (player *Player) Load(name string, typ packet.PCType) error {
    f, err := os.Open(os.Getenv("HOME") + "/.ouster/player/" + name)
    if err != nil {
        panic(err)
        // player.Name = name
        // player.Level = 150
        // player.SkinColor = 420
        // player.Alignment = 7500
        // player.STR = [3]uint16{20, 20, 20}
        // player.DEX = [3]uint16{20, 20, 20}
        // player.INT = [3]uint16{20, 20, 20}
        // player.HP = [2]uint16{472, 472}
        // player.Rank = 50
        // player.RankExp = 10700
        // player.Exp = 125
        // player.Fame = 282
        // player.Sight = 13
        // player.Bonus = 9999
        // player.Competence = 1
        // player.GuildMemberRank = 4
        // player.AdvancementLevel = 100

        // player.ZoneID =           23
        // player.ZoneX =            145
        // player.ZoneY =            237
        return err
    }
    defer f.Close()
    decoder := json.NewDecoder(f)

    switch typ {
    case packet.PC_VAMPIRE:
        err = loadVampire(player, decoder)
    case packet.PC_OUSTER:
        err = loadOuster(player, decoder)
    case packet.PC_SLAYER:
    }

    // player.ToHit = player.BaseToHit()
    // player.Defense = player.BaseDefense()
    // player.Protection = player.BaseProtection()
    // player.Damage = player.BaseDamage()

    return err
}

func loadOuster(player *Player, decoder *json.Decoder) error {
    var pcInfo data.PCOusterInfo
    err := decoder.Decode(&pcInfo)
    if err != nil {
        return err
    }

    // player.PCType = 'O'
    // player.Name = pcInfo.Name
    // player.Level = pcInfo.Level
    // player.Sex = pcInfo.Sex
    // player.HairColor = pcInfo.HairColor
    // player.MasterEffectColor = pcInfo.MasterEffectColor
    // player.Alignment = pcInfo.Alignment
    // player.STR = pcInfo.STR
    // player.DEX = pcInfo.DEX
    // player.INT = pcInfo.INT
    // player.HP = pcInfo.HP
    // player.MP = pcInfo.MP
    // player.Rank = pcInfo.Rank
    // player.RankExp = pcInfo.RankExp
    // player.Exp = pcInfo.Exp
    // player.Fame = pcInfo.Fame
    // player.Sight = pcInfo.Sight
    // player.Bonus = pcInfo.Bonus
    // player.Competence = pcInfo.Competence
    // player.GuildMemberRank = pcInfo.GuildMemberRank
    // player.AdvancementLevel = pcInfo.AdvancementLevel

    var skillInfo packet.OusterSkillInfo
    err = decoder.Decode(&skillInfo)
    if err != nil {
        return err
    }

    // player.skillslot = make([]SkillSlot, len(skillInfo.SubOusterSkillInfoList))
    // for i := 0; i < len(skillInfo.SubOusterSkillInfoList); i++ {
    //     v := &skillInfo.SubOusterSkillInfoList[i]
    //     player.skillslot[i].SkillType = v.SkillType
    //     player.skillslot[i].ExpLevel = v.ExpLevel
    //     player.skillslot[i].Interval = v.Interval
    //     player.skillslot[i].CastingTime = v.CastingTime
    // }

    // scene := zoneTable[pcInfo.ZoneID]
    // scene.Login(player, pcInfo.ZoneX, pcInfo.ZoneY)
    return nil
}

func loadVampire(player *Player, decoder *json.Decoder) error {
    var pcInfo data.PCVampireInfo
    err := decoder.Decode(&pcInfo)
    if err != nil {
        return err
    }

    // player.PCType = 'V'
    // player.Name = pcInfo.Name
    // player.Level = pcInfo.Level
    // player.Sex = pcInfo.Sex
    //  player.SkinColor = pcInfo.SkinColor
    //  player.Alignment = pcInfo.Alignment
    // player.STR = pcInfo.STR
    // player.DEX = pcInfo.DEX
    // player.INT = pcInfo.INT
    // player.HP = pcInfo.HP
    // player.Rank = pcInfo.Rank
    // player.RankExp = pcInfo.RankExp
    // player.Exp = pcInfo.Exp
    // player.Fame = pcInfo.Fame
    // player.Sight = pcInfo.Sight
    // player.Bonus = pcInfo.Bonus
    // player.Competence = pcInfo.Competence
    // player.GuildMemberRank = pcInfo.GuildMemberRank
    // player.AdvancementLevel = pcInfo.AdvancementLevel
    //
    // scene := zoneTable[pcInfo.ZoneID]
    // scene.Login(player, pcInfo.ZoneX, pcInfo.ZoneY)
    return nil
}

// func (player *Player) Save() {
//		 f, err := os.Create(os.Getenv("HOME") + "/.ouster/player/" + player.Name)
//		 if err != nil {
//				 return
//		 }
//		 encoder := json.NewEncoder(f)
//
//		 // pcInfo := player.PCInfo()
//		 skillInfo := player.SkillInfo()
//
//		 // encoder.Encode(pcInfo)
//		 encoder.Encode(skillInfo)
//
//		 f.Close()
// }

// func (player *Player) SkillInfo() packet.SkillInfo {
//		 switch player.PCType {
//		 case 'V':
//				 var ret packet.VampireSkillInfo
//				 ret.LearnNewSkill = false
//				 skillList := make([]packet.SubVampireSkillInfo, len(player.skillslot))
//				 for i := 0; i < len(player.skillslot); i++ {
//						 slot := &player.skillslot[i]
//						 skillList[i].SkillType = slot.SkillType
//						 skillList[i].Interval = slot.Interval
//						 skillList[i].CastingTime = slot.CastingTime
//				 }
//
//				 ret.SubVampireSkillInfoList = skillList
//				 return ret
//		 case 'O':
//				 var ret packet.OusterSkillInfo
//				 ret.LearnNewSkill = false
//				 skillList := make([]packet.SubOusterSkillInfo, len(player.skillslot))
//				 for i := 0; i < len(player.skillslot); i++ {
//						 slot := &player.skillslot[i]
//						 skillList[i].SkillType = slot.SkillType
//						 skillList[i].ExpLevel = slot.ExpLevel
//						 skillList[i].Interval = slot.Interval
//						 skillList[i].CastingTime = slot.CastingTime
//				 }
//
//				 ret.SubOusterSkillInfoList = skillList
//				 return ret
//		 case 'S':
//		 }
//		 return nil
// }

func Encrypt(ZoneID uint16, ServerID uint16) uint8 {
    return uint8(((ZoneID >> 8) ^ ZoneID) ^ ((ServerID + 1) << 4))
}

func (player *Player) BroadcastPacket(x uint8, y uint8, pkt packet.Packet) {
    // player.Scene.Nearby(x, y, func(watcher aoi.Entity, marker aoi.Entity) {
    //				 id := marker.Id()
    //				 if id != player.Id() {
    //						 object := player.Scene.objects[id]
    //						 if nearby, ok := object.(*Player); ok {
    //								 nearby.send <- pkt
    //						 }
    //				 }
    //		 })
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
