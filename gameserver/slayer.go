package main

import (
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
    Creature

    Name            string
    Competence      byte
    CompetenceShape byte

    Sex       Sex
    HairStyle HairStyle

    HairColor Color_t
    SkinColor Color_t

    Alignment Alignment_t

    STR [3]Attr_t
    DEX [3]Attr_t
    INI [3]Attr_t

    HP  [3]HP_t
    MP  [3]MP_t

    Damage        [3]Damage_t
    ToHit         [2]ToHit_t
    Defense       [2]Defense_t
    Protection    [2]Protection_t
    AttackSpeed   [2]Speed_t
    CriticalRatio int

    Vision     [2]Vision_t
    SkillPoint SkillPoint_t

    Fame Fame_t
    Gold Gold_t

    SkillDomainLevels [SKILL_DOMAIN_VAMPIRE]SkillLevel_t

    GoalExp [SKILL_DOMAIN_VAMPIRE]Exp_t

    // 技能糟
    SkillSlot map[SkillType_t]*SlayerSkillSlot

    WearItem [SLAYER_WEAR_MAX]*Item

    // 摩托车
    // Motocycle *Motocycle

    HPStealAmount Steal_t
    HPStealRatio  Steal_t
    MPStealAmount Steal_t
    MPStealRatio  Steal_t

    HPRegen Regen_t
    MPRegen Regen_t

    Luck Luck_t

    MPRegenTime time.Time
}

func (slayer Slayer) CreatureClass() CreatureClass {
    return CREATURE_CLASS_SLAYER
}
