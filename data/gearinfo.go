package data

import (
    "encoding/binary"
    "io"
    // . "github.com/tiancaiamao/ouster/util"
)

type GearInfo struct {
    GearSlotInfoList []GearSlotInfo
}

func (info *GearInfo) Read(reader io.Reader) error {
    var num uint8
    err := binary.Read(reader, binary.LittleEndian, &num)
    if err != nil {
        return err
    }
    info.GearSlotInfoList = make([]GearSlotInfo, num)
    for i := 0; i < int(num); i++ {
        info.GearSlotInfoList[i].Read(reader)
    }
    return nil
}

func (info *GearInfo) Write(writer io.Writer) error {
    num := uint8(len(info.GearSlotInfoList))
    binary.Write(writer, binary.LittleEndian, num)
    for i := 0; i < int(num); i++ {
        info.GearSlotInfoList[i].Write(writer)
    }
    return nil
}
