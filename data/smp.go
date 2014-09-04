package data

import (
    "encoding/binary"
    . "github.com/tiancaiamao/ouster/util"
    "io"
    "io/ioutil"
    "os"
)

type SMP struct {
    Version     string
    ZoneGroupID uint8
    Name        string
    Description string
    ZoneID      ZoneID_t
    ZoneType    ZoneType
    ZoneLevel   ZoneLevel_t

    ZoneAccessMode ZoneAccessMode
    DarkLevel      DarkLevel_t
    LightLevel     LightLevel_t
    Width          ZoneCoord_t
    Height         ZoneCoord_t

    Data []byte
}

func ReadSMP(fileName string) (*SMP, error) {
    fd, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer fd.Close()

    ret := new(SMP)
    // read zone version
    var buf [200]byte
    var strLen uint32
    binary.Read(fd, binary.LittleEndian, &strLen)
    skip := io.LimitReader(fd, int64(strLen))
    ioutil.ReadAll(skip)

    // read zone id
    binary.Read(fd, binary.LittleEndian, &ret.ZoneID)

    // read zone group id
    var ignore uint16
    binary.Read(fd, binary.LittleEndian, &ignore)

    // read zone name
    binary.Read(fd, binary.LittleEndian, &strLen)
    if strLen > 0 {
        fd.Read(buf[:strLen])
    }

    // read zone type & level
    binary.Read(fd, binary.LittleEndian, &ret.ZoneType)
    binary.Read(fd, binary.LittleEndian, &ret.ZoneLevel)

    // read zone description
    binary.Read(fd, binary.LittleEndian, &strLen)
    if strLen > 0 {
        fd.Read(buf[:strLen])
    }

    // read zone width & height
    binary.Read(fd, binary.LittleEndian, &ret.Width)
    binary.Read(fd, binary.LittleEndian, &ret.Height)

    data, err := ioutil.ReadAll(fd)
    if err != nil {
        return nil, err
    }

    ret.Data = data
    return ret, nil
}

type SSIRecord struct {
    Level  uint8
    Left   uint8
    Top    uint8
    Right  uint8
    Bottom uint8
}

type SSI []SSIRecord

func ReadSSI(fileName string) (SSI, error) {
    fd, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer fd.Close()

    var size uint32
    err = binary.Read(fd, binary.LittleEndian, &size)
    if err != nil {
        return nil, err
    }

    ret := make(SSI, size)
    var buf [5]byte
    for i := 0; i < int(size); i++ {
        _, err = fd.Read(buf[:])
        if err != nil {
            return nil, err
        }
        ret[i].Level = buf[0]
        ret[i].Left = buf[1]
        ret[i].Top = buf[2]
        ret[i].Right = buf[3]
        ret[i].Bottom = buf[4]
    }
    return ret, nil
}
