package main

import (
    "bufio"
    "encoding/binary"
    "io"
    "io/ioutil"
    "os"
)

// Read SMP file to get Zone info
func ReadSMP(fileName string, zone *Zone) error {
    fd, err := os.Open(fileName)
    if err != nil {
        return err
    }
    defer fd.Close()

    // read zone version
    var buf [200]byte
    var strLen uint32
    binary.Read(fd, binary.LittleEndian, &strLen)
    skip := io.LimitReader(fd, int64(strLen))
    ioutil.ReadAll(skip)

    // read zone id
    binary.Read(fd, binary.LittleEndian, &zone.ZoneID)

    // read zone group id
    var ignore uint16
    binary.Read(fd, binary.LittleEndian, &ignore)

    // read zone name
    binary.Read(fd, binary.LittleEndian, &strLen)
    if strLen > 0 {
        fd.Read(buf[:strLen])
    }

    // read zone type & level
    binary.Read(fd, binary.LittleEndian, &zone.ZoneType)
    binary.Read(fd, binary.LittleEndian, &zone.ZoneLevel)

    // read zone description
    binary.Read(fd, binary.LittleEndian, &strLen)
    if strLen > 0 {
        fd.Read(buf[:strLen])
    }

    // read zone width & height
    binary.Read(fd, binary.LittleEndian, &zone.Width)
    binary.Read(fd, binary.LittleEndian, &zone.Height)

    return readSMPTiles(fd, zone)
}

func readSMPTiles(fd *os.File, zone *Zone) error {
    zone.Tiles = make([]Tile, zone.Width*zone.Height)

    sectorWidth := zone.Width / SECTOR_SIZE
    sectorHeight := zone.Height / SECTOR_SIZE
    zone.Sectors = make([]Sector, sectorWidth*sectorHeight)

    for i := 0; i < int(zone.Width); i++ {
        for j := 0; j < int(zone.Height); j++ {
            sx := i / SECTOR_SIZE
            sy := i / SECTOR_SIZE
            tile := zone.Tile(i, j)
            tile.Sector = zone.Sector(sx, sy)
        }
    }

    for i := 0; i < int(sectorWidth); i++ {
        for j := 0; j < int(sectorHeight); j++ {
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
            flag, err := bufio.NewReader(fd).ReadByte()
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
                if readSMPPortal(fd, zone) != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func readSMPPortal(fd io.Reader, zone *Zone) error {
    c, err := bufio.NewReader(fd).ReadByte()
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
        size, err := bufio.NewReader(fd).ReadByte()
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

func ReadSMI(fileName string, zone *Zone) error {
    fd, err := os.Open(fileName)
    if err != nil {
        return err
    }
    defer fd.Close()

    var size int32
    binary.Read(fd, binary.LittleEndian, &size)

    var buf [5]byte
    for i := 0; i < int(size); i++ {
        fd.Read(buf[:])
        level := ZoneLevel_t(buf[0])
        left := buf[1]
        top := buf[2]
        right := buf[3]
        bottom := buf[4]

        for bx := left; bx <= right; bx++ {
            for by := top; by <= bottom; by++ {
                *(zone.Level(int(bx), int(by))) = level
            }
        }
    }
    return nil
}
