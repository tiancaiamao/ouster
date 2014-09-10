package main

import (
    "github.com/tiancaiamao/ouster/log"
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
    "io"
    "math/rand"
    "net"
    "time"
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
type Player struct {
    PlayerStatus PlayerStatus
    ZoneID       ZoneID_t
    OldZoneID    ZoneID_t

    conn         net.Conn
    packetReader *packet.Reader
    packetWriter *packet.Writer

    client <-chan packet.Packet
    send   chan<- packet.Packet
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

func InitPlayer(player *Player, conn net.Conn) {
    player.PlayerStatus = GPS_BEGIN_SESSION
    player.conn = conn

    read := make(chan packet.Packet, 1)
    write := make(chan packet.Packet, 1)
    player.send = write
    player.client = read

    go func() {
        reader := packet.NewReader()
        player.packetReader = reader
        for {
            data, err := reader.Read(player.conn)
            if err != nil {
                if _, ok := err.(packet.NotImplementError); ok {
                    log.Println("读到一个未实现的包:", data.PacketID())
                } else {
                    if err == io.EOF {
                        log.Println("后台gouroutine读客户端失败了:", err)
                        player.conn.Close()
                        close(read)
                        return
                    } else {
                        log.Println("这是一个严重的错误:", err)
                        return
                    }
                }
            }
            log.Println("读到了一个packet:", data)
            read <- data
        }
    }()

    go func() {
        writer := packet.NewWriter()
        player.packetWriter = writer
        for {
            pkt := <-write
            log.Debugln("write channel get a pkt ", pkt)
            err := writer.Write(player.conn, pkt)
            if err != nil {
                log.Println(err)
                continue
            }
        }
    }()
}

func (player *Player) sendPacket(pkt packet.Packet) {
    player.send <- pkt
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

type SkillEffect struct {
    Id   int
    To   uint32
    Succ bool
    Hurt int
}
