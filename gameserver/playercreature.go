package main

type PlayerCreatureInterface interface {
    CreatureInterface
    PlayerCreatureInstance() *PlayerCreature
}

// 多重继承Creature和Player
type PlayerCreature struct {
    Creature
    Player

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
