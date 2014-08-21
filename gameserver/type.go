package main

type ObjectID_t uint32
type ZoneCoord_t int16
type Dir_t uint8
type Sight_t uint8
type Resist_t uint16

type GuildID_t uint16
type MagicDomain int

const (
    MAGIC_DOMAIN_NO_DOMAIN MagicDomain = iota
    MAGIC_DOMAIN_POISON
    MAGIC_DOMAIN_ACID
    MAGIC_DOMAIN_CURSE
    MAGIC_DOMAIN_BLOOD
    MAGIC_DOMAIN_MAX
)

type Damage_t uint16
type SpriteType_t uint16
type Luck_t uint16
type Attr_t uint16

type Sex int8

const (
    FEMALE = iota
    MALE
)

type HairStyle int8

const (
    HAIR_STYLE1 = iota
    HAIR_STYLE2
    HAIR_STYLE3
)

type Color_t uint16
type Alignment_t int
type HP_t uint16
type MP_t uint16
type ToHit_t uint16
type Defense_t uint16
type Protection_t uint16
type Speed_t uint8
type Exp_t uint32
type Level_t uint8
type Bonus_t uint16
type SkillBonus_t uint16
type Gold_t uint32
type Fame_t uint32
type SkillType_t uint16
type Silver_t uint16
type Steal_t uint8
type Regen_t uint8
type Elemental_t uint16
type MonsterType_t uint16

type Vision_t uint8
type SkillPoint_t uint8

const (
    SKILL_DOMAIN_BLADE = iota
    SKILL_DOMAIN_SWORD
    SKILL_DOMAIN_GUN
    SKILL_DOMAIN_HEAL
    SKILL_DOMAIN_ENCHANT
    SKILL_DOMAIN_ETC
    SKILL_DOMAIN_VAMPIRE
    SKILL_DOMAIN_OUSTERS
    SKILL_DOMAIN_MAX
)

type Turn_t uint32
type Moral_t uint8
type SkillLevel_t uint8
type ItemID_t uint32
type ExpLevel_t uint16
type ZoneID_t uint16
type ZoneLevel_t uint8
type DarkLevel_t uint8
type LightLevel_t uint8
type NPCType_t uint16
type Coord_t uint8
type Rank_t uint8
type RankExp_t uint32
type GuildMemberRank_t uint8

type ElementalType int8

const (
    ELEMENTAL_ANY  ElementalType = -1
    ELEMENTAL_FIRE               = iota
    ELEMENTAL_WATER
    ELEMENTAL_EARTH
    ELEMENTAL_WIND
    ELEMENTAL_SUM
    ELEMENTAL_MAX
)
