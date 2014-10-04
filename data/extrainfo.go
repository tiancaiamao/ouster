package data

import (
    "encoding/binary"
    "io"
)

type ExtraInfo struct {
    ExtraSlotInfoList []PCItemInfo
}

func (info *ExtraInfo) Read(reader io.Reader) error {
    var num uint8
    binary.Read(reader, binary.LittleEndian, &num)
    info.ExtraSlotInfoList = make([]PCItemInfo, num)
    for i := 0; i < int(num); i++ {
        info.ExtraSlotInfoList[i].Read(reader)
    }
    return nil
}

func (info *ExtraInfo) Write(writer io.Writer) error {
    num := uint8(len(info.ExtraSlotInfoList))
    binary.Write(writer, binary.LittleEndian, num)
    for i := 0; i < int(num); i++ {
        info.ExtraSlotInfoList[i].Write(writer)
    }
    return nil
}
