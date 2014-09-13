package main

import (
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/log"
    . "github.com/tiancaiamao/ouster/util"
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
    PlayerCreature

    Competence      byte
    CompetenceShape byte

    Sex       Sex_t
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

func (vampire *Vampire) CreatureClass() CreatureClass {
    return CREATURE_CLASS_VAMPIRE
}

func (vampire *Vampire) PCInfo() data.PCInfo {
    if vampire == nil || vampire.Scene == nil {
		log.Errorln("fuck...Scene为空谁让你调这个函数了？")
        return nil
    }

    ret := &data.PCVampireInfo{
        ObjectID: vampire.ObjectID,
        Name:     vampire.Name,
        Level:    vampire.Level,
        Sex:      vampire.Sex,

        BatColor:  vampire.BatColor,
        SkinColor: vampire.SkinColor,

        Alignment: vampire.Alignment,
        Rank:      vampire.Rank,
        RankExp:   vampire.RankExp,

        Exp:   vampire.Exp,
        Fame:  vampire.Fame,
        Gold:  vampire.Gold,
        Sight: vampire.Sight,
        Bonus: vampire.Bonus,
        // HotKey:       vampire.HotKey,
        SilverDamage: vampire.SilverDamage,

        Competence: vampire.Competence,
        GuildID:    vampire.GuildID,

        GuildMemberRank: vampire.GuildMemberRank,
        UnionID:         vampire.UnionID,

        AdvancementLevel:   vampire.AdvancementLevel,
        AdvancementGoalExp: vampire.AdvancementGoalExp,

        ZoneID: vampire.Scene.ZoneID,
        ZoneX:  Coord_t(vampire.X),
        ZoneY:  Coord_t(vampire.Y),
    }

    log.Debugln("run here.......")
    ret.STR[ATTR_CURRENT] = vampire.STR[ATTR_CURRENT]
    ret.STR[ATTR_MAX] = vampire.STR[ATTR_MAX]
    ret.DEX[ATTR_CURRENT] = vampire.DEX[ATTR_CURRENT]
    ret.DEX[ATTR_MAX] = vampire.DEX[ATTR_MAX]
    ret.INI[ATTR_CURRENT] = vampire.INI[ATTR_CURRENT]
    ret.INI[ATTR_MAX] = vampire.INI[ATTR_MAX]
    ret.HP[ATTR_CURRENT] = vampire.HP[ATTR_CURRENT]
    ret.HP[ATTR_MAX] = vampire.HP[ATTR_MAX]

    return ret
}
