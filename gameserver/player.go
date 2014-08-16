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
    X   int
    Y   int
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

type Player struct {
    aoi.Entity
    Creature

    PCType byte
    // field from data.PCInfo
    Name               string
    Sex                uint8
    BatColor           uint16
    SkinColor          uint16
    HairColor          uint16
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

    skillslot []SkillSlot

    conn         net.Conn
    packetReader *darkeden.Reader
    packetWriter *darkeden.Writer

    client <-chan packet.Packet
    send   chan<- packet.Packet

    nearby    map[uint32]struct{}
    heartbeat <-chan time.Time
    ticker    uint32

    computation chan func()
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

func (player *Player) SkillSlot(SkillType uint16) *SkillSlot {
    for i := 0; i < len(player.skillslot); i++ {
        if player.skillslot[i].SkillType == SkillType {
            return &player.skillslot[i]
        }
    }
    return nil
}

func NewPlayer(conn net.Conn) *Player {
    return &Player{
        conn:        conn,
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

func (player *Player) Load(name string, typ darkeden.PCType) error {
    f, err := os.Open(os.Getenv("HOME") + "/.ouster/player/" + name)
    if err != nil {
        panic(err)
        player.PCType = 'V'
        player.Name = name
        // player.Level = 150
        player.SkinColor = 420
        player.Alignment = 7500
        // player.STR = [3]uint16{20, 20, 20}
        // player.DEX = [3]uint16{20, 20, 20}
        // player.INT = [3]uint16{20, 20, 20}
        // player.HP = [2]uint16{472, 472}
        player.Rank = 50
        player.RankExp = 10700
        player.Exp = 125
        player.Fame = 282
        player.Sight = 13
        player.Bonus = 9999
        player.Competence = 1
        player.GuildMemberRank = 4
        player.AdvancementLevel = 100

        // player.ZoneID =           23
        // player.ZoneX =            145
        // player.ZoneY =            237
        return err
    }
    defer f.Close()
    decoder := json.NewDecoder(f)

    switch typ {
    case darkeden.PC_VAMPIRE:
        err = loadVampire(player, decoder)
    case darkeden.PC_OUSTER:
        err = loadOuster(player, decoder)
    case darkeden.PC_SLAYER:
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

    player.PCType = 'O'
    player.Name = pcInfo.Name
    // player.Level = pcInfo.Level
    player.Sex = pcInfo.Sex
    player.HairColor = pcInfo.HairColor
    player.MasterEffectColor = pcInfo.MasterEffectColor
    player.Alignment = pcInfo.Alignment
    // player.STR = pcInfo.STR
    // player.DEX = pcInfo.DEX
    // player.INT = pcInfo.INT
    // player.HP = pcInfo.HP
    // player.MP = pcInfo.MP
    player.Rank = pcInfo.Rank
    player.RankExp = pcInfo.RankExp
    player.Exp = pcInfo.Exp
    player.Fame = pcInfo.Fame
    player.Sight = pcInfo.Sight
    player.Bonus = pcInfo.Bonus
    player.Competence = pcInfo.Competence
    player.GuildMemberRank = pcInfo.GuildMemberRank
    player.AdvancementLevel = pcInfo.AdvancementLevel

    var skillInfo darkeden.OusterSkillInfo
    err = decoder.Decode(&skillInfo)
    if err != nil {
        return err
    }

    player.skillslot = make([]SkillSlot, len(skillInfo.SubOusterSkillInfoList))
    for i := 0; i < len(skillInfo.SubOusterSkillInfoList); i++ {
        v := &skillInfo.SubOusterSkillInfoList[i]
        player.skillslot[i].SkillType = v.SkillType
        player.skillslot[i].ExpLevel = v.ExpLevel
        player.skillslot[i].Interval = v.Interval
        player.skillslot[i].CastingTime = v.CastingTime
    }

    scene := zoneTable[pcInfo.ZoneID]
    scene.Login(player, pcInfo.ZoneX, pcInfo.ZoneY)
    return nil
}

func loadVampire(player *Player, decoder *json.Decoder) error {
    var pcInfo data.PCVampireInfo
    err := decoder.Decode(&pcInfo)
    if err != nil {
        return err
    }

    player.PCType = 'V'
    player.Name = pcInfo.Name
    // player.Level = pcInfo.Level
    player.Sex = pcInfo.Sex
    player.SkinColor = pcInfo.SkinColor
    player.Alignment = pcInfo.Alignment
    // player.STR = pcInfo.STR
    // player.DEX = pcInfo.DEX
    // player.INT = pcInfo.INT
    // player.HP = pcInfo.HP
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
    scene.Login(player, pcInfo.ZoneX, pcInfo.ZoneY)
    return nil
}

func (player *Player) Save() {
    f, err := os.Create(os.Getenv("HOME") + "/.ouster/player/" + player.Name)
    if err != nil {
        return
    }
    encoder := json.NewEncoder(f)

    pcInfo := player.PCInfo()
    skillInfo := player.SkillInfo()

    encoder.Encode(pcInfo)
    encoder.Encode(skillInfo)

    f.Close()
}

func (player *Player) SkillInfo() darkeden.SkillInfo {
    switch player.PCType {
    case 'V':
        var ret darkeden.VampireSkillInfo
        ret.LearnNewSkill = false
        skillList := make([]darkeden.SubVampireSkillInfo, len(player.skillslot))
        for i := 0; i < len(player.skillslot); i++ {
            slot := &player.skillslot[i]
            skillList[i].SkillType = slot.SkillType
            skillList[i].Interval = slot.Interval
            skillList[i].CastingTime = slot.CastingTime
        }

        ret.SubVampireSkillInfoList = skillList
        return ret
    case 'O':
        var ret darkeden.OusterSkillInfo
        ret.LearnNewSkill = false
        skillList := make([]darkeden.SubOusterSkillInfo, len(player.skillslot))
        for i := 0; i < len(player.skillslot); i++ {
            slot := &player.skillslot[i]
            skillList[i].SkillType = slot.SkillType
            skillList[i].ExpLevel = slot.ExpLevel
            skillList[i].Interval = slot.Interval
            skillList[i].CastingTime = slot.CastingTime
        }

        ret.SubOusterSkillInfoList = skillList
        return ret
    case 'S':
    }
    return nil
}

func (player *Player) PCInfo() data.PCInfo {
    switch player.PCType {
    case 'V':
        return &data.PCVampireInfo{
            ObjectID: player.Id(),
            Name:     player.Name,
            // Level:    player.Level,
            Sex: player.Sex,

            BatColor:          player.BatColor,
            SkinColor:         player.SkinColor,
            MasterEffectColor: player.MasterEffectColor,

            Alignment: player.Alignment,
            // STR:       player.STR,
            // DEX:       player.DEX,
            // INT:			 player.INT,

            // HP: player.HP,

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

            ZoneID: player.Scene.ZoneID,
            ZoneX:  player.X(),
            ZoneY:  player.Y(),
        }
    case 'O':
        info := &data.PCOusterInfo{
            ObjectID: player.Id(),
            Name:     player.Name,
            // Level:    player.Level,
            Sex: player.Sex,

            HairColor:         player.HairColor,
            MasterEffectColor: player.MasterEffectColor,

            Alignment: player.Alignment,
            // STR:       player.STR,
            // DEX:       player.DEX,
            // INT:       player.INT,

            // HP: player.HP,
            // MP: player.MP,

            Rank:    player.Rank,
            RankExp: player.RankExp,

            Exp:          player.Exp,
            Fame:         player.Fame,
            Gold:         player.Gold,
            Sight:        player.Sight,
            Bonus:        player.Bonus,
            SilverDamage: player.SilverDamage,

            Competence: player.Competence,
            GuildID:    player.GuildID,

            GuildMemberRank: player.GuildMemberRank,
            UnionID:         player.UnionID,

            AdvancementLevel:   player.AdvancementLevel,
            AdvancementGoalExp: player.AdvancementGoalExp,

            ZoneID: player.Scene.ZoneID,
            ZoneX:  player.X(),
            ZoneY:  player.Y(),
        }

        if info.SkillBonus == 0 {
            info.SkillBonus = 9999
            log.Println("SKillBonus =========== 0!!!")
        }
        if info.GuildID == 0 {
            info.GuildID = 66
            log.Println("GuildID =========== 0!!!")
        }
        return info
    case 'S':
    }

    panic("not reached")
    return nil
}

func Encrypt(ZoneID uint16, ServerID uint16) uint8 {
    return uint8(((ZoneID >> 8) ^ ZoneID) ^ ((ServerID + 1) << 4))
}

func (player *Player) handleClientMessage(pkt packet.Packet) {
    switch pkt.Id() {
    case darkeden.PACKET_CG_CONNECT:
        raw := pkt.(*darkeden.CGConnectPacket)
        player.Load(raw.PCName, darkeden.PCType(raw.PCType))

        info := &darkeden.GCUpdateInfoPacket{
            PCType: player.PCType,
            PCInfo: player.PCInfo(),
            ZoneID: player.Scene.ZoneID,
            ZoneX:  player.X(),
            ZoneY:  player.Y(),

            GameTime: darkeden.GameTimeType{
                Year:  1983,
                Month: 8,
                Day:   19,

                Hour:   12,
                Minute: 28,
                Second: 16,
            },

            DarkLevel:  13,
            LightLevel: 6,

            MonsterTypes: []uint16{5, 6, 7, 8},

            Premium: 17,
            NicknameInfo: darkeden.NicknameInfo{
                NicknameID: 32560,
            },

            GuildUnionUserType: 2,
        }

        code := Encrypt(player.Scene.ZoneID, 1)
        player.packetReader.Code = code
        player.packetWriter.Code = code

        if info.PCType == 'O' {
            info.GearInfo = darkeden.GearInfo{
                GearSlotInfoList: []darkeden.GearSlotInfo{
                    darkeden.GearSlotInfo{
                        ObjectID:   12494,
                        ItemClass:  59,
                        ItemType:   14,
                        Durability: 6700,
                        Grade:      4,
                        ItemNum:    1,

                        SlotID: 3,
                    },
                },
            }
        }

        player.send <- info
        player.send <- &darkeden.GCPetInfoPacket{}
    case darkeden.PACKET_CG_READY:
        log.Println("get a CG Ready Packet!!!")
        player.send <- &darkeden.GCSetPositionPacket{
            X:   player.X(),
            Y:   player.Y(),
            Dir: 2,
        }

        var skillInfo darkeden.GCSkillInfoPacket
        switch player.PCType {
        case 'V':
            skillInfo.PCType = darkeden.PC_VAMPIRE
        case 'O':
            skillInfo.PCType = darkeden.PC_OUSTER
        case 'S':
            skillInfo.PCType = darkeden.PC_SLAYER
        }
        skillInfo.PCSkillInfoList = []darkeden.SkillInfo{
            player.SkillInfo(),
        }
        player.send <- &skillInfo
    case darkeden.PACKET_CG_MOVE:
        player.Scene.agent <- AgentMessage{
            Player: player,
            Msg:    pkt,
        }
    case darkeden.PACKET_CG_SAY:
        say := pkt.(*darkeden.CGSayPacket)
        log.Println("say:", say.Message)
    case darkeden.PACKET_CG_ATTACK:
        // attack := pkt.(darkeden.CGAttackPacket)
        //    log.Println(" attack monster ", attack.ObjectID)
        //    target := player.Scene.objects[attack.ObjectID]
        //    if monster, ok := target.(*Monster); ok {
        //        hit := HitTest(player.ToHit, monster.Defense)
        //        if hit {
        //            player.send <- darkeden.GCAttackMeleeOK1{
        //                ObjectID: monster.Id(),
        //            }
        //
        //            damage := 1
        //            if player.Damage > monster.Protection {
        //                damage = int(player.Damage - monster.Protection)
        //            }
        //
        //            log.Println("send attack SkillOutput to scene..........")
        //            player.Scene.agent <- AgentMessage{
        //                Player: player,
        //                Msg: SkillOutput{
        //                    MonsterID: attack.ObjectID,
        //                    Damage:    damage,
        //                },
        //            }
        //        } else {
        //            player.send <- &darkeden.GCSkillFailed1Packet{}
        //        }
        //    }
    case darkeden.PACKET_CG_SKILL_TO_SELF:
        skill := pkt.(darkeden.CGSkillToSelfPacket)
        switch skill.SkillType {
        case SKILL_INVISIBILITY:
            ok := &darkeden.GCSkillToSelfOK1{
                SkillType: SKILL_INVISIBILITY,
                CEffectID: 181,
                Duration:  0,
                Grade:     0,
            }
            ok.Short = make(map[darkeden.ModifyType]uint16)
            ok.Short[12] = 180 + 256
            player.send <- ok
        default:
            log.Println("unknown SkillToSelf type:", skill.SkillType)
        }
    case darkeden.PACKET_CG_SKILL_TO_OBJECT:
        skill := pkt.(darkeden.CGSkillToObjectPacket)
        player.SkillToObject(skill)
    case darkeden.PACKET_CG_SKILL_TO_TILE:
        skill := pkt.(darkeden.CGSkillToTilePacket)
        player.SkillToTile(skill)
    case darkeden.PACKET_CG_BLOOD_DRAIN:
    case darkeden.PACKET_CG_VERIFY_TIME:
    case darkeden.PACKET_CG_LOGOUT:
        player.Save()
        return
    }
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

func (player *Player) SkillToTile(packet darkeden.CGSkillToTilePacket) {
    skillInfo, ok := skillTable[packet.SkillType]
    if !ok {
        log.Println("unknown SkillToTie type:", packet.SkillType, packet.X, packet.Y, packet.CEffectID)
        return
    }

    handler := skillInfo.Handler
    tileHandler, ok := handler.(PlayerSkillToTileHandler)
    if !ok {
        log.Println("error ", packet.SkillType, "not implement SkillTileHandler!!!")
        return
    }
    tileHandler.ExecuteP2T(player, packet.X, packet.Y)
}

func (player *Player) SkillToObject(packet darkeden.CGSkillToObjectPacket) {
    // skillInfo, ok := skillTable[packet.SkillType]
    // if !ok {
    //     log.Println("unknown SkillToObject type:", packet.SkillType, packet.TargetObjectID, packet.CEffectID)
    //     return
    // }

    // target := player.Scene.objects[packet.TargetObjectID]
    // if monster, ok := target.(*Monster); ok {
    //     handler := skillInfo.Handler
    //     if toObjectHandler, ok := handler.(PlayerSkillToMonsterHandler); ok {
    //         toObjectHandler.ExecuteP2M(player, monster)
    //     } else {
    //         log.Println("can't execute skill ", packet.SkillType)
    //     }
    // } else {
    //
    // }
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
    // obj := this.Scene.objects[id]
    // if _, ok := obj.(*Monster); ok {
    //     log.Println("it's a monster...send message")
    //     monster := obj.(*Monster)
    //     if _, ok := this.nearby[id]; !ok {
    //         this.nearby[id] = struct{}{}
    //
    //         addMonster := &darkeden.GCAddMonster{
    //             ObjectID:    uint32(id),
    //             MonsterType: monster.MonsterType,
    //             MonsterName: "test",
    //             X:           monster.X(),
    //             Y:           monster.Y(),
    //             Dir:         2,
    //             CurrentHP:   monster.HP[ATTR_CURRENT],
    //             MaxHP:       monster.MaxHP(),
    //         }
    //
    //         this.send <- addMonster
    //         monster.flag |= flagActive
    //         log.Println("monster ", id, "set to active", monster.flag)
    //         monster.Enemies = append(monster.Enemies, this.Id())
    //     } else {
    //
    //     }
    // }
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
