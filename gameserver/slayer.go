package main

import (
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
    "time"
)

const (
    SLAYER_WEAR_HEAD = iota
    SLAYER_WEAR_NECK
    SLAYER_WEAR_BODY
    SLAYER_WEAR_LEFTHAND
    SLAYER_WEAR_RIGHTHAND
    SLAYER_WEAR_HAND3
    SLAYER_WEAR_BELT
    SLAYER_WEAR_LEG
    SLAYER_WEAR_WRIST1
    SLAYER_WEAR_WRIST2
    SLAYER_WEAR_FINGER1
    SLAYER_WEAR_FINGER2
    SLAYER_WEAR_FINGER3
    SLAYER_WEAR_FINGER4
    SLAYER_WEAR_FOOT
    SLAYER_WEAR_ZAP1
    SLAYER_WEAR_ZAP2
    SLAYER_WEAR_ZAP3
    SLAYER_WEAR_ZAP4
    SLAYER_WEAR_PDA
    SLAYER_WEAR_SHOULDER
    SLAYER_WEAR_MAX
)

type Slayer struct {
    PlayerCreature

    Competence      byte
    CompetenceShape byte

    Sex       Sex_t
    HairStyle HairStyle

    HairColor Color_t
    SkinColor Color_t

    Alignment Alignment_t

    MP  [3]MP_t

    Damage        [3]Damage_t
    ToHit         [2]ToHit_t
    Defense       [2]Defense_t
    Protection    [2]Protection_t
    AttackSpeed   [2]Speed_t
    CriticalRatio int

    Vision     [2]Vision_t
    SkillPoint SkillPoint_t

    DomainLevels [6]SkillLevel_t
    DomainExps   [6]SkillExp_t
    HotKey       [4]SkillType_t

    Fame Fame_t
    Gold Gold_t

    AttrBonus Bonus_t

    SkillDomainLevels [SKILL_DOMAIN_VAMPIRE]SkillLevel_t

    GoalExp [SKILL_DOMAIN_VAMPIRE]Exp_t

    // 技能糟
    SkillSlot map[SkillType_t]*SlayerSkillSlot

    WearItem [SLAYER_WEAR_MAX]*Item

    // 摩托车
    // Motocycle *Motocycle

    STRExp Exp_t
    DEXExp Exp_t
    INIExp Exp_t

    GuildName string

    HPStealAmount Steal_t
    HPStealRatio  Steal_t
    MPStealAmount Steal_t
    MPStealRatio  Steal_t

    HPRegen Regen_t
    MPRegen Regen_t

    Luck Luck_t

    MPRegenTime time.Time
}

func (slayer *Slayer) CreatureClass() CreatureClass {
    return CREATURE_CLASS_SLAYER
}

func (slayer *Slayer) PCInfo() data.PCInfo {
    ret := &data.PCSlayerInfo{
        ObjectID:          slayer.ObjectID,
        Name:              slayer.Name,
        Sex:               slayer.Sex,
        HairStyle:         slayer.HairStyle,
        HairColor:         slayer.HairColor,
        SkinColor:         slayer.SkinColor,
        MasterEffectColor: slayer.MasterEffectColor,
        Alignment:         slayer.Alignment,

        STRExp:  slayer.STRExp,
        DEXExp:  slayer.DEXExp,
        INIExp:  slayer.INIExp,
        Rank:    slayer.Rank,
        RankExp: slayer.RankExp,

        Fame: slayer.Fame,
        Gold: slayer.Gold,

        Sight: slayer.Sight,

        Competence:      slayer.Competence,
        GuildID:         slayer.GuildID,
        GuildName:       slayer.GuildName,
        GuildMemberRank: slayer.GuildMemberRank,

        UnionID:            slayer.UnionID,
        AdvancementLevel:   slayer.AdvancementLevel,
        AdvancementGoalExp: slayer.AdvancementGoalExp,
        AttrBonus:          slayer.AttrBonus,
    }

    for i := 0; i < 3; i++ {
        ret.STR[i] = slayer.STR[i]
        ret.DEX[i] = slayer.DEX[i]
        ret.INI[i] = slayer.INI[i]
    }

    for i := 0; i < 2; i++ {
        ret.HP[i] = slayer.HP[i]
        ret.MP[i] = slayer.MP[i]
    }

    for i := 0; i < 6; i++ {
        ret.DomainLevels[i] = slayer.DomainLevels[i]
        ret.DomainExps[i] = slayer.DomainExps[i]
    }

    for i := 0; i < 4; i++ {
        ret.HotKey[i] = slayer.HotKey[i]
    }

    return ret
}

// TODO
func (slayer *Slayer) computeDamage(creature CreatureInterface, bCritical bool) Damage_t {
    return 0
}

func (slayer *Slayer) getProtection() Protection_t {
    return slayer.Protection[ATTR_CURRENT]
}

// TODO
func (slayer *Slayer) getWearItem(class int) ItemInterface {
    return nil
}

// TODO
func increaseDomainExp(pSlayer *Slayer,
    SKILL_DOMAIN_BLADE SkillPoint_t,
    Point int,
    _GCAttackMeleeOK1 packet.Packet,
    level Level_t) {

}
