package main

import (
    "encoding/binary"
    "errors"
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/log"
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
    "math/rand"
    // "io"
    "bytes"
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
    // SectorWidth    int
    // SectorHeight	 int

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

func (zone *Zone) isMasterLair() bool {
    // TODO
    return true
}

// TODO
// func (zone *Zone) getMasterLairManager() *MasterLairManager {
//		 return nil
// }

func (zone *Zone) getZoneLevel(x, y ZoneCoord_t) ZoneLevel_t {
    // TODO
    return 0
}
func (zone *Zone) getWidth() ZoneCoord_t {
    return zone.Width
}

func (zone *Zone) getHeight() ZoneCoord_t {
    return zone.Height
}

func NewZone(zoneID ZoneID_t) *Zone {
    return &Zone{
        ZoneID: zoneID,
    }
}

func (zone *Zone) load(smp *data.SMP, ssi data.SSI) {
    // 读取ZoneInfo
    zone.ZoneID = smp.ZoneID

    zoneInfo := gZoneInfoManager.GetZoneInfo(zone.ZoneID)
    if zoneInfo != nil {
        zone.IsPKZone = zoneInfo.IsPKZone
        zone.IsNoPortalZone = zoneInfo.IsNoPortalZone
        zone.IsMasterLair = zoneInfo.IsMasterLair
        zone.IsHolyLand = zoneInfo.IsHolyLand
    }

    zone.ZoneType = smp.ZoneType
    zone.ZoneLevel = smp.ZoneLevel
    zone.ZoneAccessMode = smp.ZoneAccessMode

    zone.Width = smp.Width
    zone.Height = smp.Height
    zone.Tiles = make([]Tile, int(zone.Width)*int(zone.Height))

    sectorWidth := int(zone.Width)/SECTOR_SIZE + 1
    sectorHeight := int(zone.Height)/SECTOR_SIZE + 1
    zone.Sectors = make([]Sector, sectorWidth*sectorHeight)

    for i := 0; i < int(zone.Width); i++ {
        for j := 0; j < int(zone.Height); j++ {
            sx := i / SECTOR_SIZE
            sy := i / SECTOR_SIZE
            tile := zone.Tile(i, j)
            // log.Debugln(sx, sy)
            tile.Sector = zone.Sector(sx, sy)
        }
    }

    for i := 1; i < sectorWidth-1; i++ {
        for j := 1; j < sectorHeight-1; j++ {
            for d := 0; d < 9; d++ {
                sectorx := i + dirMoveMask[d].X
                sectory := j + dirMoveMask[d].Y

                zone.Sector(i, j).NearbySectors[d] = zone.Sector(sectorx, sectory)
            }
        }
    }

    for y := 0; y < int(zone.Height); y++ {
        for x := 0; x < int(zone.Width); x++ {
            tile := zone.Tile(x, y)
            flag, err := smp.Data.ReadByte()
            if (flag & 0x01) != 0x00 {
                tile.SetBlocked(MOVE_MODE_BURROWING)
            }
            if (flag & 0x02) != 0x00 {
                tile.SetBlocked(MOVE_MODE_WALKING)
            }
            if (flag & 0x04) != 0x00 {
                tile.SetBlocked(MOVE_MODE_FLYING)
            }
            if flag == 0 {
                zone.MonsterRegenPosition = append(zone.MonsterRegenPosition, BPOINT{byte(x), byte(y)})
            }
            if (flag&0x07) != 0x07 && (zone.IsMasterLair || zone.ZoneID == 3002) {
                zone.EmptyTilePosition = append(zone.EmptyTilePosition, BPOINT{byte(x), byte(y)})
            }
            if (flag & 0x80) != 0x00 {
                if readSMPPortal(smp.Data, zone) != nil {
                    panic(err)
                }
            }
        }
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

    zone.Levels = make([]ZoneLevel_t, len(zone.Tiles))
    for i := 0; i < len(ssi); i++ {
        record := ssi[i]
        for bx := record.Left; bx < record.Right; bx++ {
            for by := record.Bottom; by < record.Top; by++ {
                if int(bx)*int(by) >= len(zone.Tiles) || ZoneCoord_t(bx) >= zone.Width || ZoneCoord_t(by) >= zone.Height {
                    // log.Debugln(bx, by, len(zone.Tiles))
                } else {
                    *(zone.Level(int(bx), int(by))) = record.Level
                }
            }
        }
    }

    zone.loadItem()
    zone.loadNPC()
    zone.loadEffect()
}

func readSMPPortal(fd *bytes.Buffer, zone *Zone) error {
    c, err := fd.ReadByte()
    if err != nil {
        return err
    }

    typ := PortalType(c)
    // portalType := PORTAL_NORMAL
    // bAddPortal := true
    var targetZoneID ZoneID_t
    var targetX byte
    var targetY byte

    if typ == PORTAL_NORMAL || typ == PORTAL_SLAYER ||
        typ == PORTAL_VAMPIRE || typ == PORTAL_OUSTER {
        binary.Read(fd, binary.LittleEndian, &targetZoneID)
        binary.Read(fd, binary.LittleEndian, &targetX)
        binary.Read(fd, binary.LittleEndian, &targetY)
    } else if typ == PORTAL_MULTI_TARGET {
        size, err := fd.ReadByte()
        if err != nil {
            return err
        }

        for i := 0; i < int(size); i++ {
            binary.Read(fd, binary.LittleEndian, &targetZoneID)
            binary.Read(fd, binary.LittleEndian, &targetX)
            binary.Read(fd, binary.LittleEndian, &targetY)
        }
        //...
    } else if typ == PORTAL_GUILD {
        binary.Read(fd, binary.LittleEndian, &targetZoneID)
        binary.Read(fd, binary.LittleEndian, &targetX)
        binary.Read(fd, binary.LittleEndian, &targetY)
        //...
    } else if typ == PORTAL_BATTLE {
        binary.Read(fd, binary.LittleEndian, &targetZoneID)
        binary.Read(fd, binary.LittleEndian, &targetX)
        binary.Read(fd, binary.LittleEndian, &targetY)
    }
    return nil
}

func (zone *Zone) loadEffect() {

}

func (zone *Zone) loadNPC() {

}

func (zne *Zone) loadItem() {

}

func (zone *Zone) Tile(x, y int) *Tile {
    return &zone.Tiles[y+x*int(zone.Height)]
}

func (zone *Zone) getTile(x, y ZoneCoord_t) *Tile {
    return &zone.Tiles[y+x*zone.Height]
}

func (zone *Zone) Level(x, y int) *ZoneLevel_t {
    return &zone.Levels[y+x*int(zone.Height)]
}

func (zone *Zone) Sector(x, y int) *Sector {
    sectorHeight := zone.Height / SECTOR_SIZE
    return &zone.Sectors[y+x*int(sectorHeight)]
}

// 只允许scene访问，不允许其它goroutine访问
func (zone *Zone) movePC(agent *Agent, cx ZoneCoord_t, cy ZoneCoord_t, dir Dir_t) {
    pc := agent.PlayerCreatureInstance()
    if !pc.IsAbleToMove() {
        agent.sendPacket(packet.GCMoveErrorPacket{
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
            agent.sendPacket(packet.GCMoveErrorPacket{
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
        agent.sendPacket(packet.GCMoveErrorPacket{
            X:  uint8(pc.X),
            Y:  uint8(pc.Y),
        })
        return
    }

    newTile := zone.Tile(int(nx), int(ny))
    oldTile := zone.Tile(int(cx), int(cy))

    // Tile上有东西了则不能移动
    if newTile.IsBlocked(pc.MoveMode) || newTile.HasCreature(pc.MoveMode) {
        agent.sendPacket(packet.GCMoveErrorPacket{
            X:  uint8(pc.X),
            Y:  uint8(pc.Y),
        })
        return
    }

    oldTile.DeleteCreature(pc.ObjectID)
    newTile.AddCreature(agent)

    // 检查地雷/陷阱

    pc.X = nx
    pc.Y = ny
    pc.Dir = dir
    agent.sendPacket(packet.GCMoveOKPacket{
        X:   uint8(nx),
        Y:   uint8(ny),
        Dir: uint8(dir),
    })

    zone.movePCBroadcast(agent, cx, cy, nx, ny)
}

func (zone *Zone) movePCBroadcast(agent *Agent, x1 ZoneCoord_t, y1 ZoneCoord_t, x2 ZoneCoord_t, y2 ZoneCoord_t) {
    // gcMove := packet.GCMovePacket{
    //     ObjectID: uint32(agent.PlayerCreatureInstance().ObjectID),
    //     X:        uint8(agent.X),
    //     Y:        uint8(agent.Y),
    //     Dir:      uint8(agent.Dir),
    // }

    beginX := x2 - ZoneCoord_t(MaxViewportWidth) - 1
    if beginX < 0 {
        beginX = 0
    }
    endX := x2 + ZoneCoord_t(MaxViewportWidth) + 1
    if endX > zone.Width {
        endX = zone.Width
    }
    beginY := y2 - ZoneCoord_t(MaxViewportUpperHeight) - 1
    if beginY < 0 {
        beginY = 0
    }
    endY := y2 + ZoneCoord_t(MaxViewportUpperHeight) + 1
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
                        monster := v.(*Monster)
                        // 怪物进入玩家视线
                        agent.sendPacket(&packet.GCAddMonster{
                            ObjectID:    monster.ObjectID,
                            MonsterType: monster.MonsterType,
                            MonsterName: monster.Name,
                            // MainColor:   monster.MainColor,
                            // SubColor:		monster.SubColor,
                            X:   Coord_t(monster.X),
                            Y:   Coord_t(monster.Y),
                            Dir: monster.Dir,
                            // EffectInfo  []EffectInfo
                            CurrentHP: monster.HP[ATTR_CURRENT],
                            MaxHP:     monster.HP[ATTR_MAX],
                            // FromFlag    byte
                        })
                        log.Debugln("发现了一个怪物...")
                        // 把玩家放到怪物的敌人列表中
                        monster.addEnemy(agent)

                    case *Slayer:
                        // pc.sendPacket(packet.GCAddSlayer{})
                        slayer := v.(*Slayer)
                        if canSee(slayer, agent) {
                            // slayer.sendPacket(packet.GCAddSlayer{})
                        } else {
                            // slayer.sendPacket(packet.GCDeleteObject{})
                        }
                    case *Vampire:
                        vampire := v.(*Vampire)
                        if canSee(agent, vampire) {
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
                        if canSee(vampire, agent) {
                            // TODO:添加或者移除
                        }
                    case *Ouster:
                        ouster := v.(*Ouster)
                        if canSee(agent, ouster) {
                            // pc.sendPacket(packet.GCAddOuster{})
                        }
                        if canSee(ouster, agent) {
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
                        case *SlayerCorpse:
                        // pc.sendPacket(packet.GCAddSlayerCorpse{})
                        case *VampireCorpse:
                        // pc.sendPacket(packet.GCAddVampireCorpse{})
                        case *OusterCorpse:
                        // pc.sendPacket(packet.GCAddOusterCorpse{})
                        case *MonsterCorpse:
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

func (zone *Zone) getRandomMonsterRegenPosition() BPOINT {
    r := rand.Intn(len(zone.MonsterRegenPosition))
    return zone.MonsterRegenPosition[r]
}

func findSuitablePosition(zone *Zone, cx ZoneCoord_t, cy ZoneCoord_t, mode MoveMode) (pt TPOINT, err error) {
    x := cx
    y := cy
    sx := ZoneCoord_t(1)
    sy := ZoneCoord_t(0)
    maxCount := 1
    count := 1
    checkCount := 300

    for checkCount >= 0 {
        if x > 0 && y > 0 && x < zone.Width && y < zone.Height {
            rTile := zone.Tile(int(x), int(y))
            if !rTile.isBlocked(mode) && !rTile.hasPortal() {
                pt.X = int(x)
                pt.Y = int(y)
                return
            }
        }
        x += sx
        y += sy
        count--
        if count == 0 {
            if sx == 0 {
                maxCount++
            }
            temp := sx
            sx = -sy
            sy = temp

            count = maxCount
        }
        checkCount--
    }

    err = errors.New("找不到可用的点了")
    return
}

func (zone *Zone) moveCreature(CreatureInterface, ZoneCoord_t, ZoneCoord_t, Dir_t) {
    // TODO
}

func (zone *Zone) moveFastMonster(*Monster, ZoneCoord_t, ZoneCoord_t, ZoneCoord_t, ZoneCoord_t, SkillType_t) bool {
    // TODO
    return true
}

func (zone *Zone) broadcastPacket(cx ZoneCoord_t, cy ZoneCoord_t, packet packet.Packet, owner *Agent) {
    var (
        ix   ZoneCoord_t
        iy   ZoneCoord_t
        endx ZoneCoord_t
        endy ZoneCoord_t
    )

    endx = cx + ZoneCoord_t(MaxViewportWidth) + 1 //+ Range
    if endx > zone.Width {
        endx = zone.Width
    }
    endy = cy + ZoneCoord_t(MaxViewportLowerHeight) + 1 //+ Range
    if endy > zone.Height {
        endy = zone.Height
    }

    ix = cx - ZoneCoord_t(MaxViewportWidth) - 1 //- Range
    if ix < 0 {
        ix = 0
    }
    iy = cy - ZoneCoord_t(MaxViewportUpperHeight) - 1 //- Range
    if iy < 0 {
        iy = 0
    }

    for ; ix <= endx; ix++ {
        for ; iy <= endy; iy++ {
            tile := zone.Tile(int(ix), int(iy))

            if tile.HasCreature(MOVE_MODE_WALKING) {
                for _, v := range tile.Objects {
                    agent, ok := v.(*Agent)
                    if ok && agent != owner {
                        if owner != nil {
                            if canSee(agent, owner) {
                                agent.sendPacket(packet)
                            }
                        } else {
                            agent.sendPacket(packet)
                        }
                    }
                }
            }
        }
    }
}
