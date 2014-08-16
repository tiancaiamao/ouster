package data

import (
    "encoding/binary"
    "errors"
    "io"
    // "log"
)

const (
    TILE_GROUND_BLOCKED = iota
    TILE_AIR_BLOCKED
    TILE_UNDERGROUND_BLOCKED
    TILE_WALKING_CREATURE
    TILE_FLYING_CREATURE
    TILE_BURROWING_CREATURE
    TILE_ITEM
    TILE_OBSTACLE
    TILE_EFFECT
    TILE_BUILDING
    TILE_PORTAL
    TILE_TERRAIN
)

const (
    NO_SAFE_ZONE       = 0x00
    SLAYER_SAFE_ZONE   = 0x01
    VAMPIRE_SAFE_ZONE  = 0x02
    COMPLETE_SAFE_ZONE = 0x04
    NO_PK_ZONE         = 0x08
    SAFE_ZONE          = 0x17
    OUSTERS_SAFE_ZONE  = 0x10
)

type Monster struct {
    MonsterType uint16
    Count       uint8
}
type Map struct {
    Name        string
    Desc        string
    ZoneType    uint8
    ZoneID      uint16
    ZoneLevel   uint8
    Width       uint16
    Height      uint16
    Data        []byte
    MonsterInfo []Monster
}

func Load(smp io.Reader) (*Map, error) {
    var tmp [32]byte

    _, err := smp.Read(tmp[:4])
    if err != nil {
        return nil, errors.New("read version length error")
    }
    versionLen := binary.LittleEndian.Uint32(tmp[:4])
    // log.Println("versionLen=", versionLen)

    version := make([]byte, versionLen)
    _, err = smp.Read(version)
    if err != nil {
        return nil, errors.New("read version error")
    }
    // log.Println("version=", string(version))

    _, err = smp.Read(tmp[:2])
    if err != nil {
        return nil, errors.New("read zone id error")
    }
    zoneID := binary.LittleEndian.Uint16(tmp[:2])
    // log.Println("zoneID=", zoneID)

    _, err = smp.Read(tmp[:2])
    if err != nil {
        return nil, errors.New("read zone group id error")
    }
    // zoneGroupID := binary.LittleEndian.Uint16(tmp[:2])
    // log.Println("zoneGroupID=", zoneGroupID)

    _, err = smp.Read(tmp[:4])
    if err != nil {
        return nil, errors.New("read zone name length error")
    }
    // log.Println("zonenameLen byte = ", tmp[:4])
    zonenameLen := binary.LittleEndian.Uint32(tmp[:4])
    // log.Println("zonenameLen=", zonenameLen)

    var zoneName []byte
    if zonenameLen > 0 {
        zoneName := make([]byte, zonenameLen)
        _, err = smp.Read(zoneName)
        if err != nil {
            return nil, errors.New("read zone name error")
        }

        // log.Println("zoneName=", zoneName)
    }

    _, err = smp.Read(tmp[:1])
    if err != nil {
        return nil, errors.New("read zone type error")
    }
    zoneType := tmp[0]
    // log.Println("zoneType=", zoneType)

    _, err = smp.Read(tmp[:1])
    if err != nil {
        return nil, errors.New("read zone level error")
    }
    zoneLevel := tmp[0]
    // log.Println("zoneLevel=", zoneLevel)

    _, err = smp.Read(tmp[:4])
    descLen := binary.LittleEndian.Uint32(tmp[:4])
    // log.Println("descLen=", descLen)

    var desc []byte
    if descLen > 0 {
        desc := make([]byte, descLen)
        _, err = smp.Read(desc)
        if err != nil {
            return nil, errors.New("read zone description error")
        }
        // log.Println("desc=", desc)
    }

    _, err = smp.Read(tmp[:2])
    if err != nil {
        return nil, errors.New("read width error")
    }
    width := binary.LittleEndian.Uint16(tmp[:2])
    // log.Println("width=", width)

    _, err = smp.Read(tmp[:2])
    if err != nil {
        return nil, errors.New("read height error")
    }
    height := binary.LittleEndian.Uint16(tmp[:2])
    // log.Println("height=", height)

    // NOTE: use width*height directly overflow uint16 and get 0 !!!
    flags := make([]byte, int(width)*int(height))
    _, err = smp.Read(flags)
    if err != nil {
        return nil, errors.New("read flag error")
    }

    // for y := 0; y < height; y++ {
    // 	for x := 0; x < width; x++ {
    //
    // 	}
    // }
    return &Map{
        Name:      string(zoneName),
        Desc:      string(desc),
        ZoneType:  zoneType,
        ZoneID:    zoneID,
        ZoneLevel: zoneLevel,
        Width:     width,
        Height:    height,
        Data:      flags,
    }, nil
}
