package main

import (
    "encoding/json"
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
    "os"
)

type PlayerCreatureInterface interface {
    CreatureInterface
    PlayerCreatureInstance() *PlayerCreature
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
    AdvancementLevel   uint8
    AdvancementGoalExp uint32
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

func (pc PlayerCreature) PlayerCreatureInstance() *PlayerCreature {
    return &pc
}

func LoadPlayerCreature(name string, typ packet.PCType) (PlayerCreatureInterface, error) {
    f, err := os.Open(os.Getenv("HOME") + "/.ouster/player/" + name)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    decoder := json.NewDecoder(f)

    var ret PlayerCreatureInterface
    switch typ {
    case packet.PC_VAMPIRE:
        err = loadVampire(decoder)
    case packet.PC_OUSTER:
        ret, err = loadOuster(decoder)
    case packet.PC_SLAYER:
    }

    return ret, err
}

func loadOuster(decoder *json.Decoder) (ouster Ouster, err error) {
    var pcInfo data.PCOusterInfo
    err = decoder.Decode(&pcInfo)
    if err != nil {
        return
    }

    ouster.Name = pcInfo.Name
    // ouster.Level = pcInfo.Level
    // ouster.Sex = pcInfo.Sex
    // ouster.HairColor = pcInfo.HairColor
    // ouster.MasterEffectColor = pcInfo.MasterEffectColor
    // ouster.Alignment = pcInfo.Alignment
    // ouster.STR = pcInfo.STR
    // ouster.DEX = pcInfo.DEX
    // ouster.INT = pcInfo.INT
    // ouster.HP = pcInfo.HP
    // ouster.MP = pcInfo.MP
    // ouster.Rank = pcInfo.Rank
    // ouster.RankExp = pcInfo.RankExp
    // ouster.Exp = pcInfo.Exp
    // ouster.Fame = pcInfo.Fame
    // ouster.Sight = pcInfo.Sight
    // ouster.Bonus = pcInfo.Bonus
    // ouster.Competence = pcInfo.Competence
    // ouster.GuildMemberRank = pcInfo.GuildMemberRank
    // ouster.AdvancementLevel = pcInfo.AdvancementLevel

    // ouster.ZoneID =           23
    // ouster.X =            145
    // ouster.Y =            237

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

    // scene := zoneTable[pcInfo.ZoneID]
    // scene.Login(player, pcInfo.ZoneX, pcInfo.ZoneY)
    return
}

func loadVampire(decoder *json.Decoder) error {
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
