package data

import (
    "encoding/binary"
    "io"
)

type InventoryInfo struct {
    InventorySlotInfoList []InventorySlotInfo
}

func (info *InventoryInfo) Size() uint32 {
    sz := uint32(1)
    for i := 0; i < len(info.InventorySlotInfoList); i++ {
        sz += info.InventorySlotInfoList[i].Size()
    }
    return sz
}

func (info *InventoryInfo) Write(writer io.Writer) error {
    num := uint8(len(info.InventorySlotInfoList))
    binary.Write(writer, binary.LittleEndian, num)
    for i := 0; i < int(num); i++ {
        info.InventorySlotInfoList[i].Write(writer)
    }
    return nil
}

func (info *InventoryInfo) Read(reader io.Reader) error {
    var num uint8
    binary.Read(reader, binary.LittleEndian, &num)
    info.InventorySlotInfoList = make([]InventorySlotInfo, num)
    for i := 0; i < int(num); i++ {
        info.InventorySlotInfoList[i].Read(reader)
    }
    return nil
}
