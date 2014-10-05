package data

import (
    "encoding/binary"
    . "github.com/tiancaiamao/ouster/util"
    "io"
)

type SubItemInfo struct {
    ObjectID ObjectID_t
    IClass   byte
    ItemType ItemType_t
    ItemNum  ItemNum_t
    SlotID   SlotID_t
}

func (info *SubItemInfo) Size() uint32 {
    return 9
}

func (info *SubItemInfo) Read(reader io.Reader) error {
    binary.Read(reader, binary.LittleEndian, &info.ObjectID)
    binary.Read(reader, binary.LittleEndian, &info.IClass)
    binary.Read(reader, binary.LittleEndian, &info.ItemType)
    binary.Read(reader, binary.LittleEndian, &info.ItemNum)
    binary.Read(reader, binary.LittleEndian, &info.SlotID)
    return nil
}

func (info *SubItemInfo) Write(writer io.Writer) error {
    binary.Write(writer, binary.LittleEndian, &info.ObjectID)
    binary.Write(writer, binary.LittleEndian, &info.IClass)
    binary.Write(writer, binary.LittleEndian, &info.ItemType)
    binary.Write(writer, binary.LittleEndian, &info.ItemNum)
    binary.Write(writer, binary.LittleEndian, &info.SlotID)
    return nil
}
