package main

import (
    "github.com/tiancaiamao/ouster/packet"
)

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

// 只允许scene访问，不允许其它goroutine访问
func (zone *Zone) movePC(pcItf PlayerCreatureInterface, cx ZoneCoord_t, cy ZoneCoord_t, dir Dir_t) {
    pc := pcItf.PlayerCreatureInstance()
    if !pc.IsAbleToMove() {
        pc.sendPacket(packet.GCMoveErrorPacket{
            X:  uint8(pc.X),
            Y:  uint8(pc.Y),
        })
        return
    }

    // 检查做弊
    if cx != pc.X || cy != pc.Y {
        difX := cx - pc.X
        difY := cy - pc.Y
        if difX < 0 {
            difX = -difX
        }
        if difY < 0 {
            difY = -difY
        }
        if difX > 6 || difY > 6 {
            pc.sendPacket(packet.GCMoveErrorPacket{
                X:  uint8(pc.X),
                Y:  uint8(pc.Y),
            })
            return
        }
    }

    // 超出地图边界
    nx := cx
    ny := cy
    nx = nx + ZoneCoord_t(dirMoveMask[dir].X)
    ny = ny + ZoneCoord_t(dirMoveMask[dir].Y)
    if nx < 0 || nx >= zone.Width || ny < 0 || ny >= zone.Height {
        pc.sendPacket(packet.GCMoveErrorPacket{
            X:  uint8(pc.X),
            Y:  uint8(pc.Y),
        })
        return
    }

    newTile := zone.Tile(int(nx), int(ny))
    oldTile := zone.Tile(int(cx), int(cy))

    // Tile上有东西了则不能移动
    if newTile.IsBlocked(pc.MoveMode) || newTile.HasCreature(pc.MoveMode) {
        pc.sendPacket(packet.GCMoveErrorPacket{
            X:  uint8(pc.X),
            Y:  uint8(pc.Y),
        })
        return
    }

    oldTile.DeleteCreature(pc.ObjectID)
    newTile.AddCreature(pcItf)

    // 检查地雷/陷阱

    pc.X = nx
    pc.Y = ny
    pc.Dir = dir
    pc.sendPacket(packet.GCMoveOKPacket{
        X:   uint8(nx),
        Y:   uint8(ny),
        Dir: uint8(dir),
    })

    zone.movePCBroadcast(pcItf, cx, cy, nx, ny)
}

func (zone *Zone) movePCBroadcast(pcItf PlayerCreatureInterface, x1 ZoneCoord_t, y1 ZoneCoord_t, x2 ZoneCoord_t, y2 ZoneCoord_t) {
    // pc := pcItf.PlayerCreatureInstance()
    // gcMove := packet.GCMovePacket{
    //     ObjectID: uint32(pc.ObjectID),
    //     X:        uint8(pc.X),
    //     Y:        uint8(pc.Y),
    //     Dir:      uint8(pc.Dir),
    // }

    beginX := x2 - ZoneCoord_t(maxViewportWidth) - 1
    if beginX < 0 {
        beginX = 0
    }
    endX := x2 + ZoneCoord_t(maxViewportWidth) + 1
    if endX > zone.Width {
        endX = zone.Width
    }
    beginY := y2 - ZoneCoord_t(maxViewportUpperHeight) - 1
    if beginY < 0 {
        beginY = 0
    }
    endY := y2 + ZoneCoord_t(maxViewportUpperHeight) + 1
    if endY > zone.Height {
        endY = zone.Height
    }

    for i := beginX; i < endX; i++ {
        for j := beginY; j < endY; j++ {
            tile := zone.Tile(int(i), int(j))
            for _, v := range tile.Objects {
                objectClass := v.ObjectClass()
                if objectClass == OBJECT_CLASS_CREATURE {
                    switch v.(type) {
                    case *Monster:
                        // 怪物进入玩家视线
                        // pc.sendPacket(packet.MonsterAddPackt{})
                        // 把玩家放到怪物的敌人列表中
                    case Slayer:
                        // pc.sendPacket(packet.GCAddSlayer{})
                        slayer := v.(Slayer)
                        if canSee(slayer, pcItf) {
                            // slayer.sendPacket(packet.GCAddSlayer{})
                        } else {
                            // slayer.sendPacket(packet.GCDeleteObject{})
                        }
                    case Vampire:
                        vampire := v.(Vampire)
                        if canSee(pcItf, vampire) {
                            if vampire.IsFlag(EFFECT_CLASS_HIDE) {
                                // pc.sendPacket(packet.GCAddBurrowingCreaturePacket{
                                //     ObjectID: vampire.ObjectID,
                                //     Name:     vampire.Name,
                                //     X:        vampire.X,
                                //     Y:        vampire.Y,
                                // })
                            } else {
                                // pc.sendPacket(packet.GCAddVampire{})
                            }
                        }
                        if canSee(vampire, pcItf) {
                            // TODO:添加或者移除
                        }
                    case Ouster:
                        ouster := v.(Ouster)
                        if canSee(pcItf, ouster) {
                            // pc.sendPacket(packet.GCAddOuster{})
                        }
                        if canSee(ouster, pcItf) {
                            // TODO:添加或者移除
                        }
                    case *NPC:
                        // pc.sendPacket(packet.GCAddNPC{})
                    }
                } else if objectClass == OBJECT_CLASS_ITEM {
                    item := v.(ItemInterface)
                    itemClass := item.ItemClass()
                    if itemClass == ITEM_CLASS_CORPSE {
                        switch v.(type) {
                        case SlayerCorpse:
                            // pc.sendPacket(packet.GCAddSlayerCorpse{})
                        case VampireCorpse:
                            // pc.sendPacket(packet.GCAddVampireCorpse{})
                        case OusterCorpse:
                            // pc.sendPacket(packet.GCAddOusterCorpse{})
                        case MonsterCorpse:
                            // pc.sendPacket(packet.GCAddMonsterCorpse{})
                        }
                    }
                } else if objectClass == OBJECT_CLASS_EFFECT {
                    // TODO:...
                    // pc.sendPacket(packet.GCAddEffectToTile{})
                }
            }
        }
    }
}

func (zone *Zone) heartbeat() {
	zone.processMonsters()
	zone.processNPCs()
	
	zone.processEffects()
}

func (zone *Zone) processMonsters() {
}

func (zone *Zone) processNPCs() {
}

func (zone *Zone) processEffects() {
}