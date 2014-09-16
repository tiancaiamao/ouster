package util

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

type Sex_t int8

const (
    FEMALE Sex_t = iota
    MALE
)

type HairStyle int8

const (
    HAIR_STYLE1 = iota
    HAIR_STYLE2
    HAIR_STYLE3
)

type Color_t uint16
type Alignment_t int32
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
type SkillExp_t uint32

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

type ZoneType uint8

const (
    ZONE_NORMAL_FIELD = iota
    ZONE_NORMAL_DUNGEON
    ZONE_SLAYER_GUILD
    ZONE_RESERVED_SLAYER_GUILD
    ZONE_PC_VAMPIRE_LAIR
    ZONE_NPC_VAMPIRE_LAIR
    ZONE_NPC_HOME
    ZONE_NPC_SHOP
    ZONE_RANDOM_MAP
    ZONE_CASTLE
)

type ZoneAccessMode uint8

const (
    ZONE_ACCESS_PUBLIE = iota
    ZONE_ACCESS_PRIVATE
)

const (
    LEFT = iota
    LEFTDOWN
    DOWN
    RIGHTDOWN
    RIGHT
    RIGHTUP
    UP
    LEFTUP
    DIR_MAX
    DIR_NONE = DIR_MAX
)

const (
    NO_SAFE_ZONE       byte = 0x00
    SLAYER_SAFE_ZONE   byte = 0x01
    VAMPIRE_SAFE_ZONE  byte = 0x02
    COMPLETE_SAFE_ZONE byte = 0x04
    NO_PK_ZONE         byte = 0x08
    SAFE_ZONE          byte = 0x17
    OUSTERS_SAFE_ZONE  byte = 0x10
)

type Weather uint8

const (
    WEATHER_CLEAR Weather = iota
    WEATHER_RAINY
    WEATHER_SNOWY
    WEATHER_MAX
)

type MoveMode uint8

const (
    MOVE_MODE_WALKING MoveMode = iota
    MOVE_MODE_FLYING
    MOVE_MODE_BURROWING
    MOVE_MODE_MAX
)

type WeatherLevel_t uint8

type TPOINT struct {
    X   int
    Y   int
}
