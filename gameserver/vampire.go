package main

import (
    "time"
)

const (
    VAMPIRE_WEAR_NECK = iota
    VAMPIRE_WEAR_BODY
    VAMPIRE_WEAR_WRIST1
    VAMPIRE_WEAR_WRIST2
    VAMPIRE_WEAR_FINGER1
    VAMPIRE_WEAR_FINGER2
    VAMPIRE_WEAR_FINGER3
    VAMPIRE_WEAR_FINGER4
    VAMPIRE_WEAR_EARRING1
    VAMPIRE_WEAR_EARRING2
    VAMPIRE_WEAR_LEFTHAND
    VAMPIRE_WEAR_RIGHTHAND
    VAMPIRE_WEAR_AMULET1
    VAMPIRE_WEAR_AMULET2
    VAMPIRE_WEAR_AMULET3
    VAMPIRE_WEAR_AMULET4
    VAMPIRE_WEAR_ZAP1
    VAMPIRE_WEAR_ZAP2
    VAMPIRE_WEAR_ZAP3
    VAMPIRE_WEAR_ZAP4
    VAMPIRE_WEAR_DERMIS
    VAMPIRE_WEAR_PERSONA
    VAMPIRE_VAMPIRE_WEAR_MAX
)

type Vampire struct {
    Creature

    Name            string
    Competence      byte
    CompetenceShape byte

    Sex       Sex
    BatColor  Color_t
    SkinColor Color_t

    Alignment Alignment_t

    STR [3]Attr_t
    DEX [3]Attr_t
    INI [3]Attr_t

    HP  [3]HP_t

    Damage        [3]Damage_t
    ToHit         [2]ToHit_t
    Defense       [2]Defense_t
    Protection    [2]Protection_t
    AttackSpeed   [2]Speed_t
    CriticalRatio int

    GoalExp Exp_t
    Level   Level_t

    Bonus Bonus_t
    Gold  Gold_t
    Fame  Fame_t

    VisionWidth  ZoneCoord_t
    VisionHeight ZoneCoord_t

    SkillSlot map[SkillType_t]*VampireSkillSlot
    WearItem  [VAMPIRE_VAMPIRE_WEAR_MAX]*Item

    SilverDamage Silver_t

    HPStealAmount Steal_t
    HPStealRatio  Steal_t

    HPRegen      Regen_t
    HPRegenBonus Regen_t

    Luck        Luck_t
    HPRegenTime time.Time
}

func (vampire Vampire) CreatureClass() CreatureClass {
    return CREATURE_CLASS_VAMPIRE
}
