package main

import ()

type POINT struct {
    X   int
    Y   int
}

var (
    dirMoveMask = [9]POINT{
        POINT{-1, 0},  // 0 == LEFT
        POINT{-1, 1},  // 1 == LEFTDOWN
        POINT{0, 1},   // 2 == DOWN
        POINT{1, 1},   // 3 == RIGHTDOWN
        POINT{1, 0},   // 4 == RIGHT
        POINT{1, -1},  // 5 == RIGHTUP
        POINT{0, -1},  // 6 == UP
        POINT{-1, -1}, // 7 == LEFTUP
        POINT{0, 0},   // 8 == DIR_MAX, NONE
    }
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

type BPOINT struct {
    X   byte
    Y   byte
}

type Zone struct {
    ZoneID         ZoneID_t
    ZoneType       ZoneType
    ZoneLevel      ZoneLevel_t
    ZoneAccessMode ZoneAccessMode
    DarkLevel      DarkLevel_t
    LightLevel     LightLevel_t
    Width          ZoneCoord_t
    Height         ZoneCoord_t
    Tiles          []Tile
    Levels         []ZoneLevel_t
    Sectors        []Sector
    SectorWidth    int
    SectorHeight   int

    // 玩家管理
    // NPC管理
    // 怪物管理
    // Effect管理
    // 天气管理

    NPCTypes     []NPCType_t
    MonsterTypes []MonsterType_t

    IsPKZone       bool
    IsNoPortalZone bool
    IsMasterLair   bool
    IsCastle       bool
    IsHolyLand     bool
    IsCastleZone   bool

    MonsterRegenPosition []BPOINT
    EmptyTilePosition    []BPOINT
}

func NewZone(zoneID ZoneID_t) *Zone {
    return &Zone{
        ZoneID: zoneID,
    }
}

func (zone *Zone) load() {
    // 读取ZoneInfo
    zoneInfo := gZoneInfoManager.GetZoneInfo(zone.ZoneID)
    zone.IsPKZone = zoneInfo.IsPKZone
    zone.IsNoPortalZone = zoneInfo.IsNoPortalZone
    zone.IsMasterLair = zoneInfo.IsMasterLair
    zone.IsHolyLand = zoneInfo.IsHolyLand

    err := ReadSMP(zoneInfo.SMPFileName, zone)
    if err != nil {
        panic(err)
    }

    if len(zone.MonsterRegenPosition) == 0 {
        outerMinX := zone.Width / 7
        outerMinY := zone.Height / 7
        outerMaxX := zone.Width - outerMinX
        outerMaxY := zone.Height - outerMinY

        for y := outerMinY; y < outerMaxY; y++ {
            for x := outerMinX; x < outerMaxX; x++ {
                tile := zone.Tile(int(x), int(y))

                if !tile.HasPortal() && !tile.IsGroundBlocked() &&
                    !tile.IsAirBlocked() && !tile.IsUndergroundBlocked() {
                    zone.MonsterRegenPosition = append(zone.MonsterRegenPosition, BPOINT{byte(x), byte(y)})
                }
            }
        }
    }

    err = ReadSMI(zoneInfo.SMIFileName, zone)
    if err != nil {
        panic(err)
    }

    zone.loadMonster()
    zone.loadItem()
    zone.loadNPC()
    zone.loadEffect()
}

func (zone *Zone) loadEffect() {

}

func (zone *Zone) loadNPC() {

}

func (zne *Zone) loadItem() {

}

func (zone *Zone) loadMonster() {
}

func (zone *Zone) Tile(x, y int) *Tile {
    return &zone.Tiles[y+x*int(zone.Width)]
}

func (zone *Zone) Level(x, y int) *ZoneLevel_t {
    return &zone.Levels[y+x*int(zone.Height)]
}

func (zone *Zone) Sector(x, y int) *Sector {
    return &zone.Sectors[y+x*int(zone.SectorHeight)]
}
