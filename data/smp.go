package data

import (
    "encoding/binary"
    "io"
    "io/ioutil"
    "os"
)

type SMP struct {
    Version     string
    ZoneID      uint16
    ZoneGroupID uint16
    Name        string
    ZoneType    uint8
    ZoneLevel   uint8
    Desc        string
    Width       uint16
    Height      uint16
    Data        []byte
    Levels      []uint8
}

func ReadSMP(fileName string, zone *SMP) error {
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

    data, err := ioutil.ReadAll(fd)
    if err != nil {
        return err
    }

    zone.Data = data
    return nil
}

func ReadSMI(fileName string, zone *SMP) error {
    fd, err := os.Open(fileName)
    if err != nil {
        return err
    }
    defer fd.Close()

    var size uint
    binary.Read(fd, binary.LittleEndian, &size)

    var buf [5]byte
    for i := 0; i < int(size); i++ {
        fd.Read(buf[:])
        level := buf[0]
        left := buf[1]
        top := buf[2]
        right := buf[3]
        bottom := buf[4]

        for bx := left; bx <= right; bx++ {
            for by := top; by <= bottom; by++ {
                zone.Levels[by+bx*uint8(zone.Height)] = level
            }
        }
    }
    return nil
}
