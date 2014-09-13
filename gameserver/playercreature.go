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
)

type PlayerCreatureInterface interface {
    CreatureInterface
    PlayerCreatureInstance() *PlayerCreature

    PCInfo() data.PCInfo
}

type PlayerCreature struct {
    Creature

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
        // return loadSlayer(decoder)
    }
    return nil, 0, errors.New("player type error!")
}

func loadOuster(decoder *json.Decoder) (ouster *Ouster, zoneID ZoneID_t, err error) {
    var pcInfo data.PCOusterInfo
    err = decoder.Decode(&pcInfo)
    if err != nil {
        return
    }

    ouster = new(Ouster)
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

    var skillInfo packet.OusterSkillInfo
    err = decoder.Decode(&skillInfo)
    if err != nil {
        return
    }

    // player.skillslot = make([]SkillSlot, len(skillInfo.SubOusterSkillInfoList))
    // for i := 0; i < len(skillInfo.SubOusterSkillInfoList); i++ {
    //     v := &skillInfo.SubOusterSkillInfoList[i]
    //     player.skillslot[i].SkillType = v.SkillType
    //     player.skillslot[i].ExpLevel = v.ExpLevel
    //     player.skillslot[i].Interval = v.Interval
    //     player.skillslot[i].CastingTime = v.CastingTime
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
