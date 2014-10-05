package main

import (
    "encoding/json"
    "errors"
    "github.com/tiancaiamao/ouster/config"
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/log"
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
    "os"
    "time"
)

type PlayerCreatureInterface interface {
    CreatureInterface
    PlayerCreatureInstance() *PlayerCreature

    PCInfo() data.PCInfo
    SkillInfo() data.PCSkillInfo
    computeDamage(CreatureInterface, bool) Damage_t
    heartbeat()
}

type PlayerCreature struct {
    Creature

    Name string
    // 装备
    // 物品
    // 钱
    GuildID GuildID_t
    // 阶级
    Rank    Rank_t
    RankExp RankExp_t

    Exp Exp_t

    GuildMemberRank    GuildMemberRank_t
    UnionID            uint32
    AdvancementLevel   Level_t
    AdvancementGoalExp Exp_t
    // 任务管理

    MagicBonusDamage  Damage_t
    PhysicBonusDamage Damage_t

    MagicDamageReduce  Damage_t
    PhysicDamageReduce Damage_t

    // 宠物信息
    // 昵称
    // 圣书
    BaseLuck Luck_t
    // 商店

    PowerPoint int
    IsAdvanced bool
    // 转职
    AdvancedSTR       Attr_t
    AdvancedDEX       Attr_t
    AdvancedINT       Attr_t
    AdvancedAttrBonus Attr_t

    MasterEffectColor uint8

    STR [3]Attr_t
    DEX [3]Attr_t
    INI [3]Attr_t

    HP  [3]HP_t
}

func (pc *PlayerCreature) PlayerCreatureInstance() *PlayerCreature {
    return pc
}

func LoadPlayerCreature(name string, typ packet.PCType) (ptr PlayerCreatureInterface, zoneID ZoneID_t, err error) {
    fileName := config.DataDir + name
    log.Debugln("load文件", fileName)
    f, err := os.Open(fileName)
    if err != nil {
        return
    }
    defer f.Close()

    decoder := json.NewDecoder(f)

    switch typ {
    case packet.PC_VAMPIRE:
        return loadVampire(decoder)
    case packet.PC_OUSTER:
        return loadOuster(decoder)
    case packet.PC_SLAYER:
        return loadSlayer(decoder)
    }
    return nil, 0, errors.New("player type error!")
}

func loadOuster(decoder *json.Decoder) (ouster *Ouster, zoneID ZoneID_t, err error) {
    var pcInfo data.PCOusterInfo
    err = decoder.Decode(&pcInfo)
    if err != nil {
        log.Errorln("decode pcinfo failed")
        return
    }

    ouster = new(Ouster)
    ouster.Init()
    ouster.Name = pcInfo.Name
    ouster.Level = pcInfo.Level
    ouster.HairColor = pcInfo.HairColor
    ouster.MasterEffectColor = pcInfo.MasterEffectColor
    ouster.Alignment = pcInfo.Alignment
    ouster.STR = pcInfo.STR
    ouster.DEX = pcInfo.DEX
    ouster.INI = pcInfo.INI
    ouster.HP[ATTR_MAX] = pcInfo.HP[ATTR_MAX]
    ouster.HP[ATTR_CURRENT] = pcInfo.HP[ATTR_CURRENT]
    ouster.MP[ATTR_MAX] = pcInfo.MP[ATTR_MAX]
    ouster.MP[ATTR_CURRENT] = pcInfo.MP[ATTR_CURRENT]
    ouster.Rank = pcInfo.Rank
    ouster.RankExp = pcInfo.RankExp
    ouster.Exp = pcInfo.Exp
    ouster.Fame = pcInfo.Fame
    ouster.Sight = pcInfo.Sight
    ouster.Bonus = pcInfo.Bonus
    ouster.Competence = pcInfo.Competence
    ouster.GuildMemberRank = pcInfo.GuildMemberRank
    ouster.AdvancementLevel = pcInfo.AdvancementLevel

    zoneID = pcInfo.ZoneID
    ouster.X = pcInfo.ZoneX
    ouster.Y = pcInfo.ZoneY

    var skillInfo data.OusterSkillInfo
    err = decoder.Decode(&skillInfo)
    if err != nil {
        return
    }

    ouster.SkillSlot = make(map[SkillType_t]*OusterSkillSlot)
    for _, v := range skillInfo.SubOusterSkillInfoList {
        skillslot := &OusterSkillSlot{
            // Name        string
            SkillType:   v.SkillType,
            ExpLevel:    v.ExpLevel,
            Interval:    time.Duration(v.Interval) * time.Millisecond,
            CastingTime: time.Duration(v.CastingTime) * time.Millisecond,
            // RunTime     time.Time
        }
        ouster.SkillSlot[v.SkillType] = skillslot
    }

    return
}

func loadSlayer(decoder *json.Decoder) (slayer *Slayer, zoneID ZoneID_t, err error) {
    var pcInfo data.PCSlayerInfo
    err = decoder.Decode(&pcInfo)
    if err != nil {
        log.Errorln("decode pcinfo failed")
        return
    }

    slayer = new(Slayer)
    slayer.Init()
    slayer.Name = pcInfo.Name
    slayer.HairColor = pcInfo.HairColor
    slayer.MasterEffectColor = pcInfo.MasterEffectColor
    slayer.Alignment = pcInfo.Alignment
    slayer.STR = pcInfo.STR
    slayer.DEX = pcInfo.DEX
    slayer.INI = pcInfo.INI
    slayer.HP[ATTR_MAX] = pcInfo.HP[ATTR_MAX]
    slayer.HP[ATTR_CURRENT] = pcInfo.HP[ATTR_CURRENT]
    slayer.MP[ATTR_MAX] = pcInfo.MP[ATTR_MAX]
    slayer.MP[ATTR_CURRENT] = pcInfo.MP[ATTR_CURRENT]
    slayer.Rank = pcInfo.Rank
    slayer.RankExp = pcInfo.RankExp
    slayer.Fame = pcInfo.Fame
    slayer.Sight = pcInfo.Sight
    slayer.Competence = pcInfo.Competence
    slayer.GuildMemberRank = pcInfo.GuildMemberRank
    slayer.AdvancementLevel = pcInfo.AdvancementLevel

    for i := 0; i < 6; i++ {
        slayer.DomainLevels[i] = pcInfo.DomainLevels[i]
        slayer.DomainExps[i] = pcInfo.DomainExps[i]
    }

    zoneID = pcInfo.ZoneID
    slayer.X = ZoneCoord_t(pcInfo.ZoneX)
    slayer.Y = ZoneCoord_t(pcInfo.ZoneY)

    // var skillInfo packet.OusterSkillInfo
    // err = decoder.Decode(&skillInfo)
    // if err != nil {
    //     return
    // }
    //
    // slayer.SkillSlot = make(map[SkillType_t]*OusterSkillSlot)
    // for _, v := range skillInfo.SubOusterSkillInfoList {
    //     skillslot := &OusterSkillSlot{
    //         // Name        string
    //         SkillType:   v.SkillType,
    //         ExpLevel:    v.ExpLevel,
    //         Interval:    time.Duration(v.Interval) * time.Millisecond,
    //         CastingTime: time.Duration(v.CastingTime) * time.Millisecond,
    //         // RunTime     time.Time
    //     }
    //     slayer.SkillSlot[v.SkillType] = skillslot
    // }

    return
}

func loadVampire(decoder *json.Decoder) (vampire *Vampire, zoneID ZoneID_t, err error) {
    var pcInfo data.PCVampireInfo
    err = decoder.Decode(&pcInfo)
    if err != nil {
        return
    }

    vampire = new(Vampire)
    vampire.Init()
    vampire.Name = pcInfo.Name
    vampire.Level = pcInfo.Level
    vampire.Sex = pcInfo.Sex
    vampire.SkinColor = pcInfo.SkinColor
    vampire.Alignment = pcInfo.Alignment
    vampire.STR = pcInfo.STR
    vampire.DEX = pcInfo.DEX
    vampire.INI = pcInfo.INI
    vampire.HP[ATTR_MAX] = pcInfo.HP[ATTR_MAX]
    vampire.HP[ATTR_CURRENT] = pcInfo.HP[ATTR_CURRENT]
    vampire.Rank = pcInfo.Rank
    vampire.RankExp = pcInfo.RankExp
    vampire.Exp = pcInfo.Exp
    vampire.Fame = pcInfo.Fame
    vampire.Sight = pcInfo.Sight
    vampire.Bonus = pcInfo.Bonus
    vampire.Competence = pcInfo.Competence
    vampire.GuildMemberRank = pcInfo.GuildMemberRank
    vampire.AdvancementLevel = pcInfo.AdvancementLevel

    zoneID = pcInfo.ZoneID
    vampire.X = ZoneCoord_t(pcInfo.ZoneX)
    vampire.Y = ZoneCoord_t(pcInfo.ZoneY)

    if vampire.X == 0 || vampire.Y == 0 {
        log.Debugln("不科学呀", vampire.X, vampire.Y)
    }
    return
}
