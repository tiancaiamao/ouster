package data

import (
    "encoding/binary"
    "io"
)

type EffectInfo struct {
    EList []uint16
}

func (info *EffectInfo) Read(reader io.Reader) error {
    var num uint8
    binary.Read(reader, binary.LittleEndian, &num)
    info.EList = make([]uint16, num)
    for i := 0; i < int(num); i++ {
        binary.Read(reader, binary.LittleEndian, &info.EList[i])
    }
    return nil
}

func (info *EffectInfo) Write(writer io.Writer) error {
    num := len(info.EList)
    binary.Write(writer, binary.LittleEndian, uint8(num))
    for i := 0; i < int(num); i++ {
        binary.Write(writer, binary.LittleEndian, &info.EList[i])
    }
    return nil
}

func (info *EffectInfo) Size() uint32 {
    return uint32(1 + len(info.EList)*2)
}
